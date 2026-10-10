package inventory

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	GetStoreInventory(ctx context.Context, storeID uuid.UUID, page, pageSize int, filters map[string]string) ([]models.InventoryModel, int, error)
	GetProductInventoryInStore(ctx context.Context, storeID, productID uuid.UUID) (*models.InventoryModel, error)
	GetLowStockReport(ctx context.Context, storeID uuid.UUID, threshold string, page, pageSize int) ([]models.InventoryModel, int, error)
	
	GetInventoryMovements(ctx context.Context, storeID, productID uuid.UUID, page, pageSize int, filters map[string]string) ([]models.InventoryMovementModel, int, error)
	
	// Complex transaction
	ExecuteInventoryAdjustment(ctx context.Context, storeID uuid.UUID, adjustment models.CreateInventoryAdjustmentRequest, createdBy uuid.UUID, fn func(context.Context, *models.InventoryModel) error) error
}
