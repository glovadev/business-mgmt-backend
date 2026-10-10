package productcategory

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateCategory(ctx context.Context, category *models.ProductCategory) error
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.ProductCategory, error)
	GetCategoryByName(ctx context.Context, name string) (*models.ProductCategory, error)
	ListCategories(ctx context.Context) ([]models.ProductCategory, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, category *models.ProductCategory) error
	UpdateCategoryStatus(ctx context.Context, id uuid.UUID, status string) error
}
