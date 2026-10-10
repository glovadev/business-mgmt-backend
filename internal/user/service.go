package user

import (
	"context"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input data")
)

type Service interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) error
	ListUsers(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListUsersResponse, error)
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req models.UpdateUserRequest) error
	ActivateUser(ctx context.Context, id uuid.UUID) error
	DeactivateUser(ctx context.Context, id uuid.UUID) error
	AssignRoles(ctx context.Context, id uuid.UUID, req models.AssignRolesRequest) error
	AssignStores(ctx context.Context, id uuid.UUID, req models.AssignStoresRequest) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) CreateUser(ctx context.Context, req models.CreateUserRequest) error {
	if req.Email == "" || req.FullName == "" {
		return ErrInvalidInput
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		FullName:  req.FullName,
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	// In a real application, you would generate a temporary password,
	// hash it with bcrypt, and send an account setup link via email.
	
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	if len(req.RoleCodes) > 0 {
		_ = s.repo.AssignRoles(ctx, user.ID, req.RoleCodes)
	}

	if len(req.StoreIDs) > 0 {
		_ = s.repo.AssignStores(ctx, user.ID, req.StoreIDs)
	}

	return nil
}

func (s *serviceImpl) ListUsers(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := s.repo.ListUsers(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListUsersResponse{
		Users:      users,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *serviceImpl) UpdateUser(ctx context.Context, id uuid.UUID, req models.UpdateUserRequest) error {
	return s.repo.UpdateUser(ctx, id, req.FullName, req.Phone)
}

func (s *serviceImpl) ActivateUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateUserStatus(ctx, id, "ACTIVE")
}

func (s *serviceImpl) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateUserStatus(ctx, id, "INACTIVE")
}

func (s *serviceImpl) AssignRoles(ctx context.Context, id uuid.UUID, req models.AssignRolesRequest) error {
	return s.repo.AssignRoles(ctx, id, req.RoleCodes)
}

func (s *serviceImpl) AssignStores(ctx context.Context, id uuid.UUID, req models.AssignStoresRequest) error {
	return s.repo.AssignStores(ctx, id, req.StoreIDs)
}
