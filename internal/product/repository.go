package product

import (
	"context"
	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateProduct(ctx context.Context, product *models.ProductModel) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.ProductModel, error)
	GetProductBySKU(ctx context.Context, sku string) (*models.ProductModel, error)
	ListProducts(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.ProductModel, int, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, product *models.ProductModel) error
	UpdateProductStatus(ctx context.Context, id uuid.UUID, status string) error
	
	GetProductInventoryAcrossStores(ctx context.Context, id uuid.UUID) ([]models.ProductInventoryData, error)
}
