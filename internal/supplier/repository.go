package supplier

import (
	"context"
	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateSupplier(ctx context.Context, supplier *models.SupplierModel) error
	GetSupplierByID(ctx context.Context, id uuid.UUID) (*models.SupplierModel, error)
	ListSuppliers(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.SupplierModel, int, error)
	UpdateSupplier(ctx context.Context, id uuid.UUID, supplier *models.SupplierModel) error
	UpdateSupplierStatus(ctx context.Context, id uuid.UUID, status string) error
}
