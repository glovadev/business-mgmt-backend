package product

import (
	"context"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product sku already exists")
	ErrInvalidInput         = errors.New("invalid input data")
)

type Service interface {
	CreateProduct(ctx context.Context, req models.CreateProductRequest, currentUserID *uuid.UUID) (*models.ProductModel, error)
	ListProducts(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListProductsResponse, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*models.ProductModel, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req models.UpdateProductRequest) error
	ActivateProduct(ctx context.Context, id uuid.UUID) error
	DeactivateProduct(ctx context.Context, id uuid.UUID) error
	GetProductInventory(ctx context.Context, id uuid.UUID) (*models.ProductInventoryResponse, error)
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) CreateProduct(ctx context.Context, req models.CreateProductRequest, currentUserID *uuid.UUID) (*models.ProductModel, error) {
	if req.SKU == "" || req.Name == "" || req.CategoryID == "" || req.UnitOfMeasure == "" || req.DefaultSellingPrice == "" || req.GSTRate == "" {
		return nil, ErrInvalidInput
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	existing, err := s.repo.GetProductBySKU(ctx, req.SKU)
	if err == nil && existing != nil {
		return nil, ErrProductAlreadyExists
	}

	product := &models.ProductModel{
		ID:                  uuid.New(),
		SKU:                 req.SKU,
		Name:                req.Name,
		CategoryID:          categoryID,
		Description:         req.Description,
		HSNCode:             req.HSNCode,
		UnitOfMeasure:       req.UnitOfMeasure,
		DefaultSellingPrice: req.DefaultSellingPrice,
		GSTRate:             req.GSTRate,
		Status:              "ACTIVE",
		CreatedBy:           currentUserID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.repo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *serviceImpl) ListProducts(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListProductsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	products, total, err := s.repo.ListProducts(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListProductsResponse{
		Products:   products,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) GetProduct(ctx context.Context, id uuid.UUID) (*models.ProductModel, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (s *serviceImpl) UpdateProduct(ctx context.Context, id uuid.UUID, req models.UpdateProductRequest) error {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return ErrProductNotFound
	}

	if req.SKU != nil {
		existing, err := s.repo.GetProductBySKU(ctx, *req.SKU)
		if err == nil && existing != nil && existing.ID != id {
			return ErrProductAlreadyExists
		}
		product.SKU = *req.SKU
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.CategoryID != nil {
		if cid, err := uuid.Parse(*req.CategoryID); err == nil {
			product.CategoryID = cid
		}
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.HSNCode != nil {
		product.HSNCode = req.HSNCode
	}
	if req.UnitOfMeasure != nil {
		product.UnitOfMeasure = *req.UnitOfMeasure
	}
	if req.DefaultSellingPrice != nil {
		product.DefaultSellingPrice = *req.DefaultSellingPrice
	}
	if req.GSTRate != nil {
		product.GSTRate = *req.GSTRate
	}

	product.UpdatedAt = time.Now()

	return s.repo.UpdateProduct(ctx, id, product)
}

func (s *serviceImpl) ActivateProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateProductStatus(ctx, id, "ACTIVE")
}

func (s *serviceImpl) DeactivateProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateProductStatus(ctx, id, "INACTIVE")
}

func (s *serviceImpl) GetProductInventory(ctx context.Context, id uuid.UUID) (*models.ProductInventoryResponse, error) {
	_, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, ErrProductNotFound
	}

	inventory, err := s.repo.GetProductInventoryAcrossStores(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.ProductInventoryResponse{
		Data: inventory,
	}, nil
}
