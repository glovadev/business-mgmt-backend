package store

import (
	"context"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrStoreNotFound = errors.New("store not found")
	ErrInvalidInput  = errors.New("invalid input data")
)

type Service interface {
	CreateStore(ctx context.Context, req models.CreateStoreRequest, currentUserID *uuid.UUID) (*models.StoreModel, error)
	ListStores(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListStoresResponse, error)
	GetStore(ctx context.Context, id uuid.UUID) (*models.StoreModel, error)
	UpdateStore(ctx context.Context, id uuid.UUID, req models.UpdateStoreRequest) error
	ActivateStore(ctx context.Context, id uuid.UUID) error
	DeactivateStore(ctx context.Context, id uuid.UUID) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) CreateStore(ctx context.Context, req models.CreateStoreRequest, currentUserID *uuid.UUID) (*models.StoreModel, error) {
	if req.Code == "" || req.Name == "" || req.Address.Line1 == "" || req.Address.City == "" || req.Address.State == "" || req.Address.PostalCode == "" {
		return nil, ErrInvalidInput
	}

	store := &models.StoreModel{
		ID:        uuid.New(),
		Code:      req.Code,
		Name:      req.Name,
		Phone:     req.Phone,
		GSTIN:     req.GSTIN,
		Address:   req.Address,
		Status:    "ACTIVE",
		CreatedBy: currentUserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateStore(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *serviceImpl) ListStores(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListStoresResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	stores, total, err := s.repo.ListStores(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListStoresResponse{
		Stores:     stores,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) GetStore(ctx context.Context, id uuid.UUID) (*models.StoreModel, error) {
	store, err := s.repo.GetStoreByID(ctx, id)
	if err != nil {
		return nil, ErrStoreNotFound
	}
	return store, nil
}

func (s *serviceImpl) UpdateStore(ctx context.Context, id uuid.UUID, req models.UpdateStoreRequest) error {
	store, err := s.repo.GetStoreByID(ctx, id)
	if err != nil {
		return ErrStoreNotFound
	}

	if req.Name != nil {
		store.Name = *req.Name
	}
	if req.Phone != nil {
		store.Phone = req.Phone
	}
	if req.GSTIN != nil {
		store.GSTIN = req.GSTIN
	}
	if req.Address != nil {
		store.Address = *req.Address
	}
	store.UpdatedAt = time.Now()

	return s.repo.UpdateStore(ctx, id, store)
}

func (s *serviceImpl) ActivateStore(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateStoreStatus(ctx, id, "ACTIVE")
}

func (s *serviceImpl) DeactivateStore(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateStoreStatus(ctx, id, "INACTIVE")
}
