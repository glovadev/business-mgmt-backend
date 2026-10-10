package configuration

import (
	"context"

	"business-mgmt-backend/internal/models"
)

type Repository interface {
	GetCompanySettings(ctx context.Context) (*models.CompanySettingsModel, error)
	UpdateCompanySettings(ctx context.Context, settings *models.CompanySettingsModel) error
}
