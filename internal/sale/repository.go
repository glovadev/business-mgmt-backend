package sale

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	GetInvoiceByID(ctx context.Context, id uuid.UUID) (*models.InvoiceModel, error)
	GetInvoiceByNumber(ctx context.Context, number string) (*models.InvoiceModel, error)
	ListInvoices(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.InvoiceModel, int, error)
	GetInvoicePayments(ctx context.Context, invoiceID uuid.UUID) ([]models.PaymentAllocationModel, error)
	
	// Atomic Transactions
	ExecuteSaleTransaction(ctx context.Context, invoice *models.InvoiceModel, initialPayment *models.InitialPaymentModel) error
	ExecuteCancelSaleTransaction(ctx context.Context, invoiceID uuid.UUID, reason string) error
}
