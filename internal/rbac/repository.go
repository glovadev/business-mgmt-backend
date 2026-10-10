package rbac

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	ListRoles(ctx context.Context) ([]models.Role, error)
	GetRoleByID(ctx context.Context, roleID uuid.UUID) (*models.Role, error)
	ListPermissions(ctx context.Context) ([]models.Permission, error)
	UpdateRolePermissions(ctx context.Context, roleID uuid.UUID, permissionCodes []string) error
}
