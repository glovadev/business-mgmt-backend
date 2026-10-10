package store

import (
	"context"
	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateStore(ctx context.Context, store *models.StoreModel) error
	GetStoreByID(ctx context.Context, id uuid.UUID) (*models.StoreModel, error)
	ListStores(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.StoreModel, int, error)
	UpdateStore(ctx context.Context, id uuid.UUID, store *models.StoreModel) error
	UpdateStoreStatus(ctx context.Context, id uuid.UUID, status string) error
}
