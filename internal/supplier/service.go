package supplier

import (
	"context"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrSupplierNotFound = errors.New("supplier not found")
	ErrInvalidInput     = errors.New("invalid input data")
)

type Service interface {
	CreateSupplier(ctx context.Context, req models.CreateSupplierRequest) (*models.SupplierModel, error)
	ListSuppliers(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListSuppliersResponse, error)
	GetSupplier(ctx context.Context, id uuid.UUID) (*models.SupplierModel, error)
	UpdateSupplier(ctx context.Context, id uuid.UUID, req models.UpdateSupplierRequest) error
	ActivateSupplier(ctx context.Context, id uuid.UUID) error
	DeactivateSupplier(ctx context.Context, id uuid.UUID) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) CreateSupplier(ctx context.Context, req models.CreateSupplierRequest) (*models.SupplierModel, error) {
	if req.Name == "" {
		return nil, ErrInvalidInput
	}

	supplier := &models.SupplierModel{
		ID:        uuid.New(),
		Name:      req.Name,
		Phone:     req.Phone,
		Email:     req.Email,
		GSTIN:     req.GSTIN,
		Address:   req.Address,
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateSupplier(ctx, supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

func (s *serviceImpl) ListSuppliers(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListSuppliersResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	suppliers, total, err := s.repo.ListSuppliers(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListSuppliersResponse{
		Suppliers:  suppliers,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) GetSupplier(ctx context.Context, id uuid.UUID) (*models.SupplierModel, error) {
	supplier, err := s.repo.GetSupplierByID(ctx, id)
	if err != nil {
		return nil, ErrSupplierNotFound
	}
	return supplier, nil
}

func (s *serviceImpl) UpdateSupplier(ctx context.Context, id uuid.UUID, req models.UpdateSupplierRequest) error {
	supplier, err := s.repo.GetSupplierByID(ctx, id)
	if err != nil {
		return ErrSupplierNotFound
	}

	if req.Name != nil {
		supplier.Name = *req.Name
	}
	if req.Phone != nil {
		supplier.Phone = req.Phone
	}
	if req.Email != nil {
		supplier.Email = req.Email
	}
	if req.GSTIN != nil {
		supplier.GSTIN = req.GSTIN
	}
	if req.Address != nil {
		supplier.Address = req.Address
	}

	supplier.UpdatedAt = time.Now()

	return s.repo.UpdateSupplier(ctx, id, supplier)
}

func (s *serviceImpl) ActivateSupplier(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateSupplierStatus(ctx, id, "ACTIVE")
}

func (s *serviceImpl) DeactivateSupplier(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateSupplierStatus(ctx, id, "INACTIVE")
}
