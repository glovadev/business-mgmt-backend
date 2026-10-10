package reminderrule

import (
	"context"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrRuleNotFound = errors.New("reminder rule not found")
	ErrInvalidInput = errors.New("invalid input data")
)

type Service interface {
	CreateRule(ctx context.Context, req models.CreateReminderRuleRequest, currentUserID uuid.UUID) (*models.ReminderRuleModel, error)
	ListRules(ctx context.Context) (*models.ListReminderRulesResponse, error)
	UpdateRule(ctx context.Context, id uuid.UUID, req models.UpdateReminderRuleRequest) error
	ActivateRule(ctx context.Context, id uuid.UUID) error
	DeactivateRule(ctx context.Context, id uuid.UUID) error
	ReorderRules(ctx context.Context, req models.ReorderReminderRulesRequest) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) CreateRule(ctx context.Context, req models.CreateReminderRuleRequest, currentUserID uuid.UUID) (*models.ReminderRuleModel, error) {
	if req.TemplateID == "" || req.SequenceNumber <= 0 || req.DaysAfterDue < 0 {
		return nil, ErrInvalidInput
	}

	templateID, err := uuid.Parse(req.TemplateID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	rule := &models.ReminderRuleModel{
		ID:             uuid.New(),
		SequenceNumber: req.SequenceNumber,
		DaysAfterDue:   req.DaysAfterDue,
		TemplateID:     templateID,
		Status:         status,
		CreatedBy:      &currentUserID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateRule(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *serviceImpl) ListRules(ctx context.Context) (*models.ListReminderRulesResponse, error) {
	rules, err := s.repo.ListRules(ctx)
	if err != nil {
		return nil, err
	}

	return &models.ListReminderRulesResponse{
		Data: rules,
	}, nil
}

func (s *serviceImpl) UpdateRule(ctx context.Context, id uuid.UUID, req models.UpdateReminderRuleRequest) error {
	rule, err := s.repo.GetRuleByID(ctx, id)
	if err != nil {
		return ErrRuleNotFound
	}

	if req.SequenceNumber != nil {
		if *req.SequenceNumber <= 0 { return ErrInvalidInput }
		rule.SequenceNumber = *req.SequenceNumber
	}
	if req.DaysAfterDue != nil {
		if *req.DaysAfterDue < 0 { return ErrInvalidInput }
		rule.DaysAfterDue = *req.DaysAfterDue
	}
	if req.TemplateID != nil {
		tID, err := uuid.Parse(*req.TemplateID)
		if err != nil { return ErrInvalidInput }
		rule.TemplateID = tID
	}

	rule.UpdatedAt = time.Now()
	return s.repo.UpdateRule(ctx, rule)
}

func (s *serviceImpl) ActivateRule(ctx context.Context, id uuid.UUID) error {
	rule, err := s.repo.GetRuleByID(ctx, id)
	if err != nil {
		return ErrRuleNotFound
	}
	rule.Status = "ACTIVE"
	rule.UpdatedAt = time.Now()
	return s.repo.UpdateRule(ctx, rule)
}

func (s *serviceImpl) DeactivateRule(ctx context.Context, id uuid.UUID) error {
	rule, err := s.repo.GetRuleByID(ctx, id)
	if err != nil {
		return ErrRuleNotFound
	}
	rule.Status = "INACTIVE"
	rule.UpdatedAt = time.Now()
	return s.repo.UpdateRule(ctx, rule)
}

func (s *serviceImpl) ReorderRules(ctx context.Context, req models.ReorderReminderRulesRequest) error {
	if len(req.Rules) == 0 {
		return ErrInvalidInput
	}

	updates := make(map[uuid.UUID]int)
	for _, r := range req.Rules {
		id, err := uuid.Parse(r.ID)
		if err != nil { return ErrInvalidInput }
		if r.SequenceNumber <= 0 { return ErrInvalidInput }
		updates[id] = r.SequenceNumber
	}

	return s.repo.ExecuteReorderTransaction(ctx, updates)
}
