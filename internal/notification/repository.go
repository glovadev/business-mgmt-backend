package notification

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type ValidationLocker interface {
	GetInvoiceOutstandingAmount(ctx context.Context, invoiceID uuid.UUID) (float64, uuid.UUID, error) // Returns outstanding amount, customerID
	GetCustomerPhoneNumber(ctx context.Context, customerID uuid.UUID) (string, error)
}

type Repository interface {
	CreateNotification(ctx context.Context, notification *models.NotificationModel) error
	GetNotificationByID(ctx context.Context, id uuid.UUID) (*models.NotificationModel, error)
	GetNotificationByProviderMessageID(ctx context.Context, providerMsgID string) (*models.NotificationModel, error)
	ListNotifications(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.NotificationModel, int, error)
	UpdateNotification(ctx context.Context, notification *models.NotificationModel) error
}
