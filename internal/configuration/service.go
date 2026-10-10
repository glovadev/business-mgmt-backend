package configuration

import (
	"context"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrSettingsNotFound = errors.New("company settings not found")
	ErrInvalidPolicy    = errors.New("invalid policy: red_after_overdue_days must be greater than yellow_after_overdue_days, and yellow must be >= 0")
)

type Service interface {
	GetCompanySettings(ctx context.Context) (*models.CompanySettingsModel, error)
	UpdateCompanySettings(ctx context.Context, req models.UpdateCompanySettingsRequest, currentUserID uuid.UUID) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) GetCompanySettings(ctx context.Context) (*models.CompanySettingsModel, error) {
	return s.repo.GetCompanySettings(ctx)
}

func (s *serviceImpl) UpdateCompanySettings(ctx context.Context, req models.UpdateCompanySettingsRequest, currentUserID uuid.UUID) error {
	settings, err := s.repo.GetCompanySettings(ctx)
	if err != nil {
		return ErrSettingsNotFound
	}

	yellow := settings.YellowAfterOverdueDays
	red := settings.RedAfterOverdueDays

	if req.YellowAfterOverdueDays != nil {
		yellow = *req.YellowAfterOverdueDays
	}
	if req.RedAfterOverdueDays != nil {
		red = *req.RedAfterOverdueDays
	}

	if yellow < 0 || red <= yellow {
		return ErrInvalidPolicy
	}

	if req.DefaultCreditDays != nil { settings.DefaultCreditDays = *req.DefaultCreditDays }
	if req.CompanyName != nil { settings.CompanyName = *req.CompanyName }
	if req.Currency != nil { settings.Currency = *req.Currency }
	if req.Timezone != nil { settings.Timezone = *req.Timezone }
	if req.AdditionalCostAllocationMethod != nil { settings.AdditionalCostAllocationMethod = *req.AdditionalCostAllocationMethod }
	if req.IncludeGSTInLandedCost != nil { settings.IncludeGSTInLandedCost = *req.IncludeGSTInLandedCost }

	settings.YellowAfterOverdueDays = yellow
	settings.RedAfterOverdueDays = red
	settings.UpdatedAt = time.Now()
	settings.UpdatedBy = &currentUserID

	// In a real implementation, we would also trigger an Audit Log creation here
	return s.repo.UpdateCompanySettings(ctx, settings)
}
