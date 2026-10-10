package inventory

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrInventoryNotFound = errors.New("inventory record not found")
	ErrInsufficientStock = errors.New("insufficient stock for adjustment")
	ErrInvalidAdjustment = errors.New("invalid adjustment type")
	ErrInvalidInput      = errors.New("invalid input data")
)

type Service interface {
	GetStoreInventory(ctx context.Context, storeID uuid.UUID, page, pageSize int, filters map[string]string) (*models.ListInventoryResponse, error)
	GetProductInventoryInStore(ctx context.Context, storeID, productID uuid.UUID) (*models.InventoryModel, error)
	GetInventoryMovementHistory(ctx context.Context, storeID, productID uuid.UUID, page, pageSize int, filters map[string]string) (*models.ListInventoryMovementsResponse, error)
	CreateInventoryAdjustment(ctx context.Context, storeID uuid.UUID, req models.CreateInventoryAdjustmentRequest, idempotencyKey string, createdBy uuid.UUID) error
	GetLowStockReport(ctx context.Context, storeID uuid.UUID, threshold string, page, pageSize int) (*models.ListInventoryResponse, error)
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func parseStrAmt(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.3f", f)
}

func (s *serviceImpl) GetStoreInventory(ctx context.Context, storeID uuid.UUID, page, pageSize int, filters map[string]string) (*models.ListInventoryResponse, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	items, total, err := s.repo.GetStoreInventory(ctx, storeID, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListInventoryResponse{
		Data:       items,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) GetProductInventoryInStore(ctx context.Context, storeID, productID uuid.UUID) (*models.InventoryModel, error) {
	item, err := s.repo.GetProductInventoryInStore(ctx, storeID, productID)
	if err != nil {
		return nil, ErrInventoryNotFound
	}
	return item, nil
}

func (s *serviceImpl) GetInventoryMovementHistory(ctx context.Context, storeID, productID uuid.UUID, page, pageSize int, filters map[string]string) (*models.ListInventoryMovementsResponse, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	movements, total, err := s.repo.GetInventoryMovements(ctx, storeID, productID, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListInventoryMovementsResponse{
		Data:       movements,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) CreateInventoryAdjustment(ctx context.Context, storeID uuid.UUID, req models.CreateInventoryAdjustmentRequest, idempotencyKey string, createdBy uuid.UUID) error {
	if req.ProductID == "" || req.Quantity == "" {
		return ErrInvalidInput
	}
	
	qtyChange := parseStrAmt(req.Quantity)
	if qtyChange <= 0 {
		return ErrInvalidInput
	}

	if req.AdjustmentType != "ADJUSTMENT_IN" && req.AdjustmentType != "ADJUSTMENT_OUT" {
		return ErrInvalidAdjustment
	}

	return s.repo.ExecuteInventoryAdjustment(ctx, storeID, req, createdBy, func(ctx context.Context, lockedInv *models.InventoryModel) error {
		currentQty := parseStrAmt(lockedInv.AvailableQuantity)
		
		if req.AdjustmentType == "ADJUSTMENT_OUT" {
			if currentQty < qtyChange {
				return ErrInsufficientStock
			}
			lockedInv.AvailableQuantity = formatFloat(currentQty - qtyChange)
		} else {
			lockedInv.AvailableQuantity = formatFloat(currentQty + qtyChange)
		}

		// The repo function handles creating the inventory_movements record and audit log.
		return nil
	})
}

func (s *serviceImpl) GetLowStockReport(ctx context.Context, storeID uuid.UUID, threshold string, page, pageSize int) (*models.ListInventoryResponse, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	
	if threshold == "" {
		threshold = "10.000" // Default threshold if not provided
	}

	items, total, err := s.repo.GetLowStockReport(ctx, storeID, threshold, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &models.ListInventoryResponse{
		Data:       items,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
