package productcategory

import (
	"context"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrCategoryNotFound      = errors.New("product category not found")
	ErrCategoryAlreadyExists = errors.New("product category name already exists")
	ErrInvalidInput          = errors.New("invalid input data")
)

type Service interface {
	CreateCategory(ctx context.Context, req models.CreateProductCategoryRequest) (*models.ProductCategory, error)
	ListCategories(ctx context.Context) ([]models.ProductCategory, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*models.ProductCategory, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req models.UpdateProductCategoryRequest) error
	DeactivateCategory(ctx context.Context, id uuid.UUID) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) CreateCategory(ctx context.Context, req models.CreateProductCategoryRequest) (*models.ProductCategory, error) {
	if req.Name == "" {
		return nil, ErrInvalidInput
	}

	existing, err := s.repo.GetCategoryByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, ErrCategoryAlreadyExists
	}

	category := &models.ProductCategory{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Status:      "ACTIVE",
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *serviceImpl) ListCategories(ctx context.Context) ([]models.ProductCategory, error) {
	return s.repo.ListCategories(ctx)
}

func (s *serviceImpl) GetCategory(ctx context.Context, id uuid.UUID) (*models.ProductCategory, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *serviceImpl) UpdateCategory(ctx context.Context, id uuid.UUID, req models.UpdateProductCategoryRequest) error {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return ErrCategoryNotFound
	}

	if req.Name != nil {
		// Verify name uniqueness
		existing, err := s.repo.GetCategoryByName(ctx, *req.Name)
		if err == nil && existing != nil && existing.ID != id {
			return ErrCategoryAlreadyExists
		}
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = req.Description
	}

	return s.repo.UpdateCategory(ctx, id, category)
}

func (s *serviceImpl) DeactivateCategory(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateCategoryStatus(ctx, id, "INACTIVE")
}
