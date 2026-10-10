package report

import (
	"context"

	"business-mgmt-backend/internal/models"
)

type Service interface {
	GetDashboardSummary(ctx context.Context, filters map[string]string) (*models.DashboardReportResponse, error)
	GetSalesReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error)
	GetRevenueReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error)
	GetCollectionReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error)
	GetExpenseReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error)
	GetProfitReport(ctx context.Context, filters map[string]string) (*models.ProfitReportResponse, error)
	GetReceivablesReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error)
	GetCreditAgingReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error)
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) GetDashboardSummary(ctx context.Context, filters map[string]string) (*models.DashboardReportResponse, error) {
	return s.repo.GetDashboardSummary(ctx, filters)
}

func (s *serviceImpl) GetSalesReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error) {
	data, err := s.repo.GetSalesReport(ctx, filters)
	if err != nil {
		return nil, err
	}
	return &models.GenericReportResponse{Data: data}, nil
}

func (s *serviceImpl) GetRevenueReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error) {
	data, err := s.repo.GetRevenueReport(ctx, filters)
	if err != nil {
		return nil, err
	}
	return &models.GenericReportResponse{Data: data}, nil
}

func (s *serviceImpl) GetCollectionReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error) {
	data, err := s.repo.GetCollectionReport(ctx, filters)
	if err != nil {
		return nil, err
	}
	return &models.GenericReportResponse{Data: data}, nil
}

func (s *serviceImpl) GetExpenseReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error) {
	data, err := s.repo.GetExpenseReport(ctx, filters)
	if err != nil {
		return nil, err
	}
	return &models.GenericReportResponse{Data: data}, nil
}

func (s *serviceImpl) GetProfitReport(ctx context.Context, filters map[string]string) (*models.ProfitReportResponse, error) {
	return s.repo.GetProfitReport(ctx, filters)
}

func (s *serviceImpl) GetReceivablesReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error) {
	data, err := s.repo.GetReceivablesReport(ctx, filters)
	if err != nil {
		return nil, err
	}
	return &models.GenericReportResponse{Data: data}, nil
}

func (s *serviceImpl) GetCreditAgingReport(ctx context.Context, filters map[string]string) (*models.GenericReportResponse, error) {
	data, err := s.repo.GetCreditAgingReport(ctx, filters)
	if err != nil {
		return nil, err
	}
	return &models.GenericReportResponse{Data: data}, nil
}
