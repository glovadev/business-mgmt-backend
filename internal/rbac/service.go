package rbac

import (
	"context"
	"errors"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrRoleNotFound = errors.New("role not found")
)

type Service interface {
	ListRoles(ctx context.Context) ([]models.Role, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (*models.Role, error)
	ListPermissions(ctx context.Context) ([]models.Permission, error)
	UpdateRolePermissions(ctx context.Context, roleID uuid.UUID, req models.UpdateRolePermissionsRequest) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) ListRoles(ctx context.Context) ([]models.Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *serviceImpl) GetRole(ctx context.Context, roleID uuid.UUID) (*models.Role, error) {
	role, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *serviceImpl) ListPermissions(ctx context.Context) ([]models.Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func (s *serviceImpl) UpdateRolePermissions(ctx context.Context, roleID uuid.UUID, req models.UpdateRolePermissionsRequest) error {
	_, err := s.GetRole(ctx, roleID)
	if err != nil {
		return err
	}

	return s.repo.UpdateRolePermissions(ctx, roleID, req.PermissionCodes)
}
