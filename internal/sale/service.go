package sale

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
	ErrInvoiceNotFound = errors.New("invoice not found")
	ErrInvalidInput    = errors.New("invalid input data")
	ErrInvalidState    = errors.New("invalid state for operation")
)

// In a real app, this would be injected via constructor to fetch product/inventory stats
type ProductInventoryFetcher interface {
	GetProductGSTRate(ctx context.Context, productID uuid.UUID) (float64, error)
	GetAverageUnitCost(ctx context.Context, storeID, productID uuid.UUID) (float64, error)
}

type Service interface {
	CreateSale(ctx context.Context, req models.CreateSaleRequest, idempotencyKey string, currentUserID *uuid.UUID) (*models.InvoiceModel, error)
	ListSales(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListSalesResponse, error)
	GetSale(ctx context.Context, id uuid.UUID) (*models.InvoiceModel, error)
	GetInvoiceByNumber(ctx context.Context, number string) (*models.InvoiceModel, error)
	CancelSale(ctx context.Context, id uuid.UUID, reason, idempotencyKey string) error
	GetSalePayments(ctx context.Context, id uuid.UUID) ([]models.PaymentAllocationModel, error)
}

type serviceImpl struct {
	repo    Repository
	fetcher ProductInventoryFetcher
}

func NewService(repo Repository, fetcher ProductInventoryFetcher) Service {
	return &serviceImpl{repo: repo, fetcher: fetcher}
}

func parseStrAmt(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func formatFloat4(f float64) string {
	return fmt.Sprintf("%.4f", f)
}

func (s *serviceImpl) CreateSale(ctx context.Context, req models.CreateSaleRequest, idempotencyKey string, currentUserID *uuid.UUID) (*models.InvoiceModel, error) {
	if req.StoreID == "" || req.CustomerID == "" || len(req.Items) == 0 {
		return nil, ErrInvalidInput
	}

	storeID, _ := uuid.Parse(req.StoreID)
	customerID, _ := uuid.Parse(req.CustomerID)

	invoice := &models.InvoiceModel{
		ID:            uuid.New(),
		InvoiceNumber: fmt.Sprintf("INV-%d", time.Now().Unix()),
		StoreID:       storeID,
		CustomerID:    customerID,
		InvoiceDate:   req.InvoiceDate,
		DueDate:       req.DueDate,
		Status:        "CONFIRMED",
		CreatedBy:     currentUserID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	var subtotal, totalDiscount, totalGST, grandTotal float64

	invoice.Items = make([]models.InvoiceItemModel, 0, len(req.Items))
	for _, itemReq := range req.Items {
		productID, _ := uuid.Parse(itemReq.ProductID)
		qty := parseStrAmt(itemReq.Quantity)
		unitPrice := parseStrAmt(itemReq.UnitPrice)
		discount := parseStrAmt(itemReq.DiscountAmount)
		
		var gstRate, unitCost float64
		if s.fetcher != nil {
			gstRate, _ = s.fetcher.GetProductGSTRate(ctx, productID)
			unitCost, _ = s.fetcher.GetAverageUnitCost(ctx, storeID, productID)
		} else {
			// Mock fallbacks if fetcher isn't implemented
			gstRate = 18.0
			unitCost = 0.0
		}

		taxableAmount := (qty * unitPrice) - discount
		gstAmount := taxableAmount * (gstRate / 100.0)
		lineTotal := taxableAmount + gstAmount

		subtotal += (qty * unitPrice)
		totalDiscount += discount
		totalGST += gstAmount
		grandTotal += lineTotal

		invoice.Items = append(invoice.Items, models.InvoiceItemModel{
			ID:             uuid.New(),
			InvoiceID:      invoice.ID,
			ProductID:      productID,
			Quantity:       fmt.Sprintf("%.3f", qty),
			UnitPrice:      formatFloat(unitPrice),
			DiscountAmount: formatFloat(discount),
			TaxableAmount:  formatFloat(taxableAmount),
			GSTRate:        formatFloat(gstRate),
			GSTAmount:      formatFloat(gstAmount),
			LineTotal:      formatFloat(lineTotal),
			UnitCost:       formatFloat4(unitCost),
			CreatedAt:      time.Now(),
		})
	}

	invoice.Subtotal = formatFloat(subtotal)
	invoice.DiscountAmount = formatFloat(totalDiscount)
	invoice.GSTAmount = formatFloat(totalGST)
	invoice.TotalAmount = formatFloat(grandTotal)
	
	paidAmount := 0.0
	if req.InitialPayment != nil {
		paidAmount = parseStrAmt(req.InitialPayment.Amount)
	}
	
	outstandingAmount := grandTotal - paidAmount
	if outstandingAmount < 0 {
		outstandingAmount = 0
	}
	
	invoice.PaidAmount = formatFloat(paidAmount)
	invoice.OutstandingAmount = formatFloat(outstandingAmount)

	if outstandingAmount == 0 {
		invoice.PaymentStatus = "PAID"
	} else if paidAmount > 0 {
		invoice.PaymentStatus = "PARTIALLY_PAID"
	} else {
		invoice.PaymentStatus = "UNPAID"
	}

	if err := s.repo.ExecuteSaleTransaction(ctx, invoice, req.InitialPayment); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *serviceImpl) ListSales(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListSalesResponse, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	invoices, total, err := s.repo.ListInvoices(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListSalesResponse{
		Sales:      invoices,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) GetSale(ctx context.Context, id uuid.UUID) (*models.InvoiceModel, error) {
	return s.repo.GetInvoiceByID(ctx, id)
}

func (s *serviceImpl) GetInvoiceByNumber(ctx context.Context, number string) (*models.InvoiceModel, error) {
	return s.repo.GetInvoiceByNumber(ctx, number)
}

func (s *serviceImpl) CancelSale(ctx context.Context, id uuid.UUID, reason, idempotencyKey string) error {
	invoice, err := s.repo.GetInvoiceByID(ctx, id)
	if err != nil {
		return ErrInvoiceNotFound
	}

	if invoice.Status == "CANCELLED" {
		return nil // idempotent
	}
	
	paid := parseStrAmt(invoice.PaidAmount)
	if paid > 0 {
		// "A sale cannot be cancelled automatically when there are unreversed payments"
		return errors.New("cannot cancel sale with existing payments; process refund first")
	}

	return s.repo.ExecuteCancelSaleTransaction(ctx, id, reason)
}

func (s *serviceImpl) GetSalePayments(ctx context.Context, id uuid.UUID) ([]models.PaymentAllocationModel, error) {
	return s.repo.GetInvoicePayments(ctx, id)
}
