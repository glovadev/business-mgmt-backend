package auditlog

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	ListAuditLogs(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.AuditLogModel, int, error)
	GetAuditLogByID(ctx context.Context, id uuid.UUID) (*models.AuditLogModel, error)
	GetEntityAuditHistory(ctx context.Context, entityType string, entityID uuid.UUID) ([]models.AuditLogModel, error)
}
