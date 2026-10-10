package payment

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

// Define an interface for what the transaction expects from the database layer
type InvoiceLocker interface {
	GetOutstandingInvoicesByCustomer(ctx context.Context, customerID uuid.UUID) ([]*models.InvoiceModel, error)
	GetInvoiceByID(ctx context.Context, invoiceID uuid.UUID) (*models.InvoiceModel, error)
}

type Repository interface {
	GetPaymentByID(ctx context.Context, id uuid.UUID) (*models.PaymentModel, error)
	ListPayments(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.PaymentModel, int, error)
	
	// Atomic Transactions
	ExecutePaymentTransaction(ctx context.Context, payment *models.PaymentModel, allocationStrategy string, explicitAllocations []models.CreatePaymentAllocationStruct, fn func(context.Context, InvoiceLocker) ([]models.PaymentAllocationModel, error)) error
	ExecutePaymentReversalTransaction(ctx context.Context, paymentID uuid.UUID, reason string, reversedBy uuid.UUID) error
}
