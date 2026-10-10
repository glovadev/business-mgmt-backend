package notificationtemplate

import (
	"context"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrTemplateNotFound = errors.New("notification template not found")
	ErrInvalidInput     = errors.New("invalid input data")
)

type Service interface {
	CreateTemplate(ctx context.Context, req models.CreateNotificationTemplateRequest) (*models.NotificationTemplateModel, error)
	ListTemplates(ctx context.Context) (*models.ListNotificationTemplatesResponse, error)
	UpdateTemplate(ctx context.Context, id uuid.UUID, req models.UpdateNotificationTemplateRequest) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) CreateTemplate(ctx context.Context, req models.CreateNotificationTemplateRequest) (*models.NotificationTemplateModel, error) {
	if req.Code == "" || req.ProviderTemplateID == "" || req.LanguageCode == "" {
		return nil, ErrInvalidInput
	}

	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	template := &models.NotificationTemplateModel{
		ID:                 uuid.New(),
		Code:               req.Code,
		ProviderTemplateID: req.ProviderTemplateID,
		LanguageCode:       req.LanguageCode,
		BodyPreview:        req.BodyPreview,
		Status:             status,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.repo.CreateTemplate(ctx, template); err != nil {
		return nil, err
	}
	return template, nil
}

func (s *serviceImpl) ListTemplates(ctx context.Context) (*models.ListNotificationTemplatesResponse, error) {
	templates, err := s.repo.ListTemplates(ctx)
	if err != nil {
		return nil, err
	}

	return &models.ListNotificationTemplatesResponse{
		Data: templates,
	}, nil
}

func (s *serviceImpl) UpdateTemplate(ctx context.Context, id uuid.UUID, req models.UpdateNotificationTemplateRequest) error {
	template, err := s.repo.GetTemplateByID(ctx, id)
	if err != nil {
		return ErrTemplateNotFound
	}

	if req.Code != nil {
		if *req.Code == "" { return ErrInvalidInput }
		template.Code = *req.Code
	}
	if req.ProviderTemplateID != nil {
		if *req.ProviderTemplateID == "" { return ErrInvalidInput }
		template.ProviderTemplateID = *req.ProviderTemplateID
	}
	if req.LanguageCode != nil {
		if *req.LanguageCode == "" { return ErrInvalidInput }
		template.LanguageCode = *req.LanguageCode
	}
	if req.BodyPreview != nil {
		template.BodyPreview = *req.BodyPreview
	}
	if req.Status != nil {
		template.Status = *req.Status
	}

	template.UpdatedAt = time.Now()
	return s.repo.UpdateTemplate(ctx, template)
}
