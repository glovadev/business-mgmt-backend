package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationModel struct {
	ID                uuid.UUID  `json:"id"`
	CustomerID        uuid.UUID  `json:"customer_id"`
	InvoiceID         *uuid.UUID `json:"invoice_id,omitempty"`
	ReminderRuleID    *uuid.UUID `json:"reminder_rule_id,omitempty"`
	PhoneNumber       string     `json:"phone_number"`
	TriggerType       string     `json:"trigger_type"`      // AUTOMATIC, MANUAL
	NotificationType  string     `json:"notification_type"` // WHATSAPP, SMS, etc.
	TemplateID        uuid.UUID  `json:"template_id"`
	Payload           string     `json:"payload"` // JSONB stored as string
	Provider          string     `json:"provider"`
	ProviderMessageID *string    `json:"provider_message_id,omitempty"`
	Status            string     `json:"status"` // PENDING, PROCESSING, SENT, DELIVERED, FAILED
	RetryCount        int        `json:"retry_count"`
	NextRetryAt       *time.Time `json:"next_retry_at,omitempty"`
	LastError         *string    `json:"last_error,omitempty"`
	TriggeredBy       *uuid.UUID `json:"triggered_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	SentAt            *time.Time `json:"sent_at,omitempty"`
	DeliveredAt       *time.Time `json:"delivered_at,omitempty"`
}

type SendManualReminderRequest struct {
	TemplateID string  `json:"template_id"`
	Notes      *string `json:"notes,omitempty"`
}

type SendManualReminderResponse struct {
	Data struct {
		NotificationID string `json:"notification_id"`
		Status         string `json:"status"`
	} `json:"data"`
}

type ListNotificationsResponse struct {
	Notifications []NotificationModel `json:"notifications"`
	TotalCount    int                 `json:"total_count"`
	Page          int                 `json:"page"`
	PageSize      int                 `json:"page_size"`
}

type WhatsAppWebhookPayload struct {
	// Abstract structure for generic WhatsApp webhooks
	MessageID string `json:"message_id"`
	Status    string `json:"status"` // sent, delivered, read, failed
	Error     string `json:"error,omitempty"`
}
