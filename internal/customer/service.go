package customer

import (
	"context"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrInvalidInput          = errors.New("invalid input data")
)

type Service interface {
	CreateCustomer(ctx context.Context, req models.CreateCustomerRequest, currentUserID *uuid.UUID) (*models.CustomerModel, error)
	ListCustomers(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListCustomersResponse, error)
	GetCustomer(ctx context.Context, id uuid.UUID) (*models.CustomerModel, error)
	UpdateCustomer(ctx context.Context, id uuid.UUID, req models.UpdateCustomerRequest) error
	DeactivateCustomer(ctx context.Context, id uuid.UUID) error
	ActivateCustomer(ctx context.Context, id uuid.UUID) error
	GetCustomerStores(ctx context.Context, id uuid.UUID) ([]models.StoreModel, error)
	LinkCustomerToStore(ctx context.Context, customerID, storeID uuid.UUID, currentUserID *uuid.UUID) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) CreateCustomer(ctx context.Context, req models.CreateCustomerRequest, currentUserID *uuid.UUID) (*models.CustomerModel, error) {
	if req.PhoneNumber == "" || req.Name == "" || req.StoreID == "" {
		return nil, ErrInvalidInput
	}

	storeID, err := uuid.Parse(req.StoreID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	// Normalize phone in a real app, e.g. converting 9876543210 to +919876543210
	
	existing, err := s.repo.GetCustomerByPhone(ctx, req.PhoneNumber)
	if err == nil && existing != nil {
		return nil, ErrCustomerAlreadyExists
	}

	customer := &models.CustomerModel{
		ID:            uuid.New(),
		Name:          req.Name,
		PhoneNumber:   req.PhoneNumber,
		AlternatePhone: req.AlternatePhone,
		Email:         req.Email,
		Address:       req.Address,
		GSTIN:         req.GSTIN,
		WhatsappOptIn: req.WhatsappOptIn,
		Status:        "ACTIVE",
		CreatedBy:     currentUserID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.CreateCustomer(ctx, customer); err != nil {
		return nil, err
	}

	// Link customer to the store implicitly
	_ = s.repo.LinkCustomerToStore(ctx, customer.ID, storeID, currentUserID)

	return customer, nil
}

func (s *serviceImpl) ListCustomers(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListCustomersResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	customers, total, err := s.repo.ListCustomers(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListCustomersResponse{
		Customers:  customers,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) GetCustomer(ctx context.Context, id uuid.UUID) (*models.CustomerModel, error) {
	customer, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, ErrCustomerNotFound
	}
	return customer, nil
}

func (s *serviceImpl) UpdateCustomer(ctx context.Context, id uuid.UUID, req models.UpdateCustomerRequest) error {
	customer, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return ErrCustomerNotFound
	}

	if req.Name != nil {
		customer.Name = *req.Name
	}
	if req.PhoneNumber != nil {
		// Needs careful duplicate validation and audit logging before updating phone
		existing, err := s.repo.GetCustomerByPhone(ctx, *req.PhoneNumber)
		if err == nil && existing != nil && existing.ID != id {
			return ErrCustomerAlreadyExists
		}
		customer.PhoneNumber = *req.PhoneNumber
	}
	if req.AlternatePhone != nil {
		customer.AlternatePhone = req.AlternatePhone
	}
	if req.Email != nil {
		customer.Email = req.Email
	}
	if req.Address != nil {
		customer.Address = req.Address
	}
	if req.GSTIN != nil {
		customer.GSTIN = req.GSTIN
	}
	if req.WhatsappOptIn != nil {
		customer.WhatsappOptIn = *req.WhatsappOptIn
		if !*req.WhatsappOptIn {
			now := time.Now()
			customer.WhatsappOptedOutAt = &now
		} else {
			customer.WhatsappOptedOutAt = nil
		}
	}
	if req.Notes != nil {
		customer.Notes = req.Notes
	}

	customer.UpdatedAt = time.Now()

	return s.repo.UpdateCustomer(ctx, id, customer)
}

func (s *serviceImpl) DeactivateCustomer(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateCustomerStatus(ctx, id, "INACTIVE")
}

func (s *serviceImpl) ActivateCustomer(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateCustomerStatus(ctx, id, "ACTIVE")
}

func (s *serviceImpl) GetCustomerStores(ctx context.Context, id uuid.UUID) ([]models.StoreModel, error) {
	return s.repo.GetCustomerStores(ctx, id)
}

func (s *serviceImpl) LinkCustomerToStore(ctx context.Context, customerID, storeID uuid.UUID, currentUserID *uuid.UUID) error {
	return s.repo.LinkCustomerToStore(ctx, customerID, storeID, currentUserID)
}
