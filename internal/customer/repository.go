package customer

import (
	"context"
	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateCustomer(ctx context.Context, customer *models.CustomerModel) error
	GetCustomerByID(ctx context.Context, id uuid.UUID) (*models.CustomerModel, error)
	GetCustomerByPhone(ctx context.Context, phone string) (*models.CustomerModel, error)
	ListCustomers(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.CustomerModel, int, error)
	UpdateCustomer(ctx context.Context, id uuid.UUID, customer *models.CustomerModel) error
	UpdateCustomerStatus(ctx context.Context, id uuid.UUID, status string) error

	// Store Relationships
	LinkCustomerToStore(ctx context.Context, customerID, storeID uuid.UUID, addedBy *uuid.UUID) error
	GetCustomerStores(ctx context.Context, customerID uuid.UUID) ([]models.StoreModel, error)
}
