package purchase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrPurchaseNotFound  = errors.New("purchase not found")
	ErrInvalidState      = errors.New("invalid state transition")
	ErrInvalidInput      = errors.New("invalid input data")
	ErrSelfApproval      = errors.New("cannot approve own purchase")
)

type Service interface {
	CreatePurchaseDraft(ctx context.Context, req models.CreatePurchaseRequest, currentUserID *uuid.UUID) (*models.PurchaseModel, error)
	ListPurchases(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListPurchasesResponse, error)
	GetPurchase(ctx context.Context, id uuid.UUID) (*models.PurchaseModel, error)
	UpdatePurchaseDraft(ctx context.Context, id uuid.UUID, req models.CreatePurchaseRequest) error
	SubmitPurchaseForApproval(ctx context.Context, id uuid.UUID) error
	ApprovePurchase(ctx context.Context, id uuid.UUID, currentUserID *uuid.UUID) error
	RejectPurchase(ctx context.Context, id uuid.UUID, reason string) error
	CancelPurchase(ctx context.Context, id uuid.UUID) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

// Helper to convert string to float64, defaulting to 0 on error (in production use shopspring/decimal)
func parseStrAmt(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.4f", f)
}

func formatFloat2(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func (s *serviceImpl) calculatePurchaseTotals(req *models.CreatePurchaseRequest, purchase *models.PurchaseModel) {
	var subtotal, totalGST, additionalCostTotal float64

	// Additional costs
	purchase.AdditionalCosts = make([]models.PurchaseAdditionalCostModel, 0, len(req.AdditionalCosts))
	for _, ac := range req.AdditionalCosts {
		amt := parseStrAmt(ac.Amount)
		additionalCostTotal += amt
		purchase.AdditionalCosts = append(purchase.AdditionalCosts, models.PurchaseAdditionalCostModel{
			ID:          uuid.New(),
			PurchaseID:  purchase.ID,
			CostType:    ac.CostType,
			Description: ac.Description,
			Amount:      formatFloat2(amt),
		})
	}

	// Items and base amounts
	var totalQty float64
	purchase.Items = make([]models.PurchaseItemModel, 0, len(req.Items))
	for _, item := range req.Items {
		qty := parseStrAmt(item.Quantity)
		baseCost := parseStrAmt(item.BaseUnitCost)
		gstRate := parseStrAmt(item.GSTRate)

		baseAmount := qty * baseCost
		gstAmount := baseAmount * (gstRate / 100.0)

		subtotal += baseAmount
		totalGST += gstAmount
		totalQty += qty

		productID, _ := uuid.Parse(item.ProductID)
		purchase.Items = append(purchase.Items, models.PurchaseItemModel{
			ID:                 uuid.New(),
			PurchaseID:         purchase.ID,
			ProductID:          productID,
			Quantity:           formatFloat(qty),
			BaseUnitCost:       formatFloat(baseCost),
			BaseAmount:         formatFloat2(baseAmount),
			GSTRate:            formatFloat2(gstRate),
			GSTAmount:          formatFloat2(gstAmount),
			AllocatedExtraCost: "0.00",
			LandedUnitCost:     "0.0000",
			TotalAmount:        formatFloat2(baseAmount + gstAmount),
		})
	}

	// Allocate extra costs (simplified: proportion by quantity for this mock)
	for i := range purchase.Items {
		qty := parseStrAmt(purchase.Items[i].Quantity)
		if totalQty > 0 {
			proportion := qty / totalQty
			allocated := additionalCostTotal * proportion
			purchase.Items[i].AllocatedExtraCost = formatFloat2(allocated)
			
			baseCost := parseStrAmt(purchase.Items[i].BaseUnitCost)
			landedCost := baseCost + (allocated / qty)
			purchase.Items[i].LandedUnitCost = formatFloat(landedCost)
		}
	}

	purchase.Subtotal = formatFloat2(subtotal)
	purchase.GSTAmount = formatFloat2(totalGST)
	purchase.AdditionalCostTotal = formatFloat2(additionalCostTotal)
	purchase.GrandTotal = formatFloat2(subtotal + totalGST + additionalCostTotal)
}

func (s *serviceImpl) CreatePurchaseDraft(ctx context.Context, req models.CreatePurchaseRequest, currentUserID *uuid.UUID) (*models.PurchaseModel, error) {
	if req.StoreID == "" || len(req.Items) == 0 {
		return nil, ErrInvalidInput
	}

	storeID, err := uuid.Parse(req.StoreID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	var supplierID *uuid.UUID
	if req.SupplierID != nil {
		if id, err := uuid.Parse(*req.SupplierID); err == nil {
			supplierID = &id
		}
	}

	purchase := &models.PurchaseModel{
		ID:                    uuid.New(),
		PurchaseNumber:        fmt.Sprintf("PUR-%d", time.Now().Unix()),
		StoreID:               storeID,
		SupplierID:            supplierID,
		SupplierInvoiceNumber: req.SupplierInvoiceNumber,
		PurchaseDate:          req.PurchaseDate,
		Status:                "DRAFT",
		Notes:                 req.Notes,
		CreatedBy:             currentUserID,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	s.calculatePurchaseTotals(&req, purchase)

	if err := s.repo.CreatePurchase(ctx, purchase); err != nil {
		return nil, err
	}

	return purchase, nil
}

func (s *serviceImpl) ListPurchases(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListPurchasesResponse, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	purchases, total, err := s.repo.ListPurchases(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListPurchasesResponse{
		Purchases:  purchases,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) GetPurchase(ctx context.Context, id uuid.UUID) (*models.PurchaseModel, error) {
	return s.repo.GetPurchaseByID(ctx, id)
}

func (s *serviceImpl) UpdatePurchaseDraft(ctx context.Context, id uuid.UUID, req models.CreatePurchaseRequest) error {
	purchase, err := s.repo.GetPurchaseByID(ctx, id)
	if err != nil {
		return ErrPurchaseNotFound
	}
	if purchase.Status != "DRAFT" {
		return ErrInvalidState
	}

	// For simplicity, recompute everything
	purchase.SupplierInvoiceNumber = req.SupplierInvoiceNumber
	purchase.PurchaseDate = req.PurchaseDate
	purchase.Notes = req.Notes
	purchase.UpdatedAt = time.Now()

	s.calculatePurchaseTotals(&req, purchase)

	return s.repo.UpdatePurchase(ctx, purchase)
}

func (s *serviceImpl) SubmitPurchaseForApproval(ctx context.Context, id uuid.UUID) error {
	purchase, err := s.repo.GetPurchaseByID(ctx, id)
	if err != nil {
		return ErrPurchaseNotFound
	}
	if purchase.Status != "DRAFT" {
		return ErrInvalidState
	}
	
	// Create approval history event 'SUBMITTED' (omitted in mock)

	return s.repo.UpdatePurchaseStatus(ctx, id, "PENDING_APPROVAL")
}

func (s *serviceImpl) ApprovePurchase(ctx context.Context, id uuid.UUID, currentUserID *uuid.UUID) error {
	purchase, err := s.repo.GetPurchaseByID(ctx, id)
	if err != nil {
		return ErrPurchaseNotFound
	}
	if purchase.Status != "PENDING_APPROVAL" {
		return ErrInvalidState
	}
	
	if purchase.CreatedBy != nil && currentUserID != nil && *purchase.CreatedBy == *currentUserID {
		return ErrSelfApproval
	}

	return s.repo.ExecuteApprovalTransaction(ctx, id, currentUserID, func(ctx context.Context, lockedPurchase *models.PurchaseModel) error {
		// Business transaction logic:
		// 1. Calculate landed costs (already done during creation)
		// 2. Lock inventory rows
		// 3. Update weighted average cost
		// 4. Increase inventory
		// 5. Create inventory movements
		// 6. Mark purchase APPROVED
		// 7. Create approval history
		// 8. Create audit log
		
		lockedPurchase.Status = "APPROVED"
		now := time.Now()
		lockedPurchase.ApprovedAt = &now
		lockedPurchase.ApprovedBy = currentUserID
		lockedPurchase.UpdatedAt = now

		return nil
	})
}

func (s *serviceImpl) RejectPurchase(ctx context.Context, id uuid.UUID, reason string) error {
	purchase, err := s.repo.GetPurchaseByID(ctx, id)
	if err != nil {
		return ErrPurchaseNotFound
	}
	if purchase.Status != "PENDING_APPROVAL" {
		return ErrInvalidState
	}

	return s.repo.UpdatePurchaseStatus(ctx, id, "REJECTED")
}

func (s *serviceImpl) CancelPurchase(ctx context.Context, id uuid.UUID) error {
	purchase, err := s.repo.GetPurchaseByID(ctx, id)
	if err != nil {
		return ErrPurchaseNotFound
	}
	if purchase.Status != "DRAFT" {
		return ErrInvalidState
	}

	return s.repo.UpdatePurchaseStatus(ctx, id, "CANCELLED")
}
