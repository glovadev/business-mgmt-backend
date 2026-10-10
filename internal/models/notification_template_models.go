package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationTemplateModel struct {
	ID                 uuid.UUID `json:"id"`
	Code               string    `json:"code"`
	ProviderTemplateID string    `json:"provider_template_id"`
	LanguageCode       string    `json:"language_code"`
	BodyPreview        string    `json:"body_preview"`
	Status             string    `json:"status"` // ACTIVE, INACTIVE
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CreateNotificationTemplateRequest struct {
	Code               string `json:"code"`
	ProviderTemplateID string `json:"provider_template_id"`
	LanguageCode       string `json:"language_code"`
	BodyPreview        string `json:"body_preview"`
	Status             string `json:"status"`
}

type UpdateNotificationTemplateRequest struct {
	Code               *string `json:"code,omitempty"`
	ProviderTemplateID *string `json:"provider_template_id,omitempty"`
	LanguageCode       *string `json:"language_code,omitempty"`
	BodyPreview        *string `json:"body_preview,omitempty"`
	Status             *string `json:"status,omitempty"`
}

type ListNotificationTemplatesResponse struct {
	Data []NotificationTemplateModel `json:"data"`
}
