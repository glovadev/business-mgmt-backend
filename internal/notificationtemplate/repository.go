package notificationtemplate

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateTemplate(ctx context.Context, template *models.NotificationTemplateModel) error
	ListTemplates(ctx context.Context) ([]models.NotificationTemplateModel, error)
	GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.NotificationTemplateModel, error)
	UpdateTemplate(ctx context.Context, template *models.NotificationTemplateModel) error
}
