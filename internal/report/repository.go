package report

import (
	"context"

	"business-mgmt-backend/internal/models"
)

type Repository interface {
	GetDashboardSummary(ctx context.Context, filters map[string]string) (*models.DashboardReportResponse, error)
	GetSalesReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error)
	GetRevenueReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error)
	GetCollectionReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error)
	GetExpenseReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error)
	GetProfitReport(ctx context.Context, filters map[string]string) (*models.ProfitReportResponse, error)
	GetReceivablesReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error)
	GetCreditAgingReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error)
}
