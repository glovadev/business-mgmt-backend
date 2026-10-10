package notification

import (
	"context"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvoiceFullyPaid     = errors.New("invoice is fully paid, cannot send reminder")
	ErrValidation           = errors.New("validation failed")
	ErrCannotRetry          = errors.New("notification cannot be retried (must be FAILED)")
)

type Service interface {
	SendManualReminder(ctx context.Context, invoiceID uuid.UUID, req models.SendManualReminderRequest, currentUserID uuid.UUID) (*models.SendManualReminderResponse, error)
	ListNotifications(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListNotificationsResponse, error)
	GetNotification(ctx context.Context, id uuid.UUID) (*models.NotificationModel, error)
	RetryFailedNotification(ctx context.Context, id uuid.UUID) error
	HandleWhatsAppWebhook(ctx context.Context, payload models.WhatsAppWebhookPayload, signature string) error
}

type serviceImpl struct {
	repo      Repository
	validator ValidationLocker
}

func NewService(repo Repository, validator ValidationLocker) Service {
	return &serviceImpl{repo: repo, validator: validator}
}

func (s *serviceImpl) SendManualReminder(ctx context.Context, invoiceID uuid.UUID, req models.SendManualReminderRequest, currentUserID uuid.UUID) (*models.SendManualReminderResponse, error) {
	if req.TemplateID == "" {
		return nil, ErrValidation
	}

	templateID, err := uuid.Parse(req.TemplateID)
	if err != nil {
		return nil, ErrValidation
	}

	outstanding, customerID, err := s.validator.GetInvoiceOutstandingAmount(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	if outstanding <= 0 {
		return nil, ErrInvoiceFullyPaid
	}

	phone, err := s.validator.GetCustomerPhoneNumber(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if phone == "" {
		return nil, errors.New("customer has no registered phone number")
	}

	// Create Notification Record
	notification := &models.NotificationModel{
		ID:               uuid.New(),
		CustomerID:       customerID,
		InvoiceID:        &invoiceID,
		PhoneNumber:      phone,
		TriggerType:      "MANUAL",
		NotificationType: "WHATSAPP", // Hardcoded for this V1 use-case
		TemplateID:       templateID,
		Payload:          "{}", // Would map real payload params here
		Provider:         "TWILIO",
		Status:           "PENDING",
		TriggeredBy:      &currentUserID,
		CreatedAt:        time.Now(),
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		return nil, err
	}

	// In a real implementation: Enqueue worker job here via Outbox pattern or RabbitMQ/SQS

	resp := &models.SendManualReminderResponse{}
	resp.Data.NotificationID = notification.ID.String()
	resp.Data.Status = "PENDING"

	return resp, nil
}

func (s *serviceImpl) ListNotifications(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListNotificationsResponse, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	notifications, total, err := s.repo.ListNotifications(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListNotificationsResponse{
		Notifications: notifications,
		TotalCount:    total,
		Page:          page,
		PageSize:      pageSize,
	}, nil
}

func (s *serviceImpl) GetNotification(ctx context.Context, id uuid.UUID) (*models.NotificationModel, error) {
	return s.repo.GetNotificationByID(ctx, id)
}

func (s *serviceImpl) RetryFailedNotification(ctx context.Context, id uuid.UUID) error {
	notification, err := s.repo.GetNotificationByID(ctx, id)
	if err != nil {
		return ErrNotificationNotFound
	}

	if notification.Status != "FAILED" {
		return ErrCannotRetry
	}

	notification.Status = "PENDING"
	notification.RetryCount += 1
	// next_retry_at logic goes here

	return s.repo.UpdateNotification(ctx, notification)
}

func (s *serviceImpl) HandleWhatsAppWebhook(ctx context.Context, payload models.WhatsAppWebhookPayload, signature string) error {
	// Validate Signature here...

	notification, err := s.repo.GetNotificationByProviderMessageID(ctx, payload.MessageID)
	if err != nil {
		// Log and safely return. We don't want to error out the webhook provider
		return nil 
	}

	// Idempotently update based on provider's state machine
	var targetStatus string
	switch payload.Status {
	case "sent":
		targetStatus = "SENT"
	case "delivered", "read":
		targetStatus = "DELIVERED"
	case "failed":
		targetStatus = "FAILED"
	}

	if targetStatus != "" {
		notification.Status = targetStatus
		now := time.Now()
		
		if targetStatus == "SENT" {
			notification.SentAt = &now
		} else if targetStatus == "DELIVERED" {
			notification.DeliveredAt = &now
		} else if targetStatus == "FAILED" {
			notification.LastError = &payload.Error
		}

		_ = s.repo.UpdateNotification(ctx, notification)
	}

	return nil
}
