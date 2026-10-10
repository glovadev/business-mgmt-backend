package auditlog

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Service interface {
	ListAuditLogs(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListAuditLogsResponse, error)
	GetAuditLog(ctx context.Context, id uuid.UUID) (*models.AuditLogModel, error)
	GetEntityAuditHistory(ctx context.Context, entityType string, entityID uuid.UUID) ([]models.AuditLogModel, error)
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) ListAuditLogs(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListAuditLogsResponse, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	logs, total, err := s.repo.ListAuditLogs(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListAuditLogsResponse{
		AuditLogs:  logs,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) GetAuditLog(ctx context.Context, id uuid.UUID) (*models.AuditLogModel, error) {
	return s.repo.GetAuditLogByID(ctx, id)
}

func (s *serviceImpl) GetEntityAuditHistory(ctx context.Context, entityType string, entityID uuid.UUID) ([]models.AuditLogModel, error) {
	return s.repo.GetEntityAuditHistory(ctx, entityType, entityID)
}
