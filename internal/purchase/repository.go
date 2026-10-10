package purchase

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreatePurchase(ctx context.Context, purchase *models.PurchaseModel) error
	GetPurchaseByID(ctx context.Context, id uuid.UUID) (*models.PurchaseModel, error)
	ListPurchases(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.PurchaseModel, int, error)
	UpdatePurchase(ctx context.Context, purchase *models.PurchaseModel) error
	UpdatePurchaseStatus(ctx context.Context, id uuid.UUID, status string) error

	// For the complex approval transaction
	ExecuteApprovalTransaction(ctx context.Context, purchaseID uuid.UUID, approvedBy *uuid.UUID, fn func(context.Context, *models.PurchaseModel) error) error
}
