package user

import (
	"context"
	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListUsers(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.User, int, error)
	UpdateUser(ctx context.Context, id uuid.UUID, fullName, phone *string) error
	UpdateUserStatus(ctx context.Context, id uuid.UUID, status string) error
	AssignRoles(ctx context.Context, id uuid.UUID, roleCodes []string) error
	AssignStores(ctx context.Context, id uuid.UUID, storeIDs []uuid.UUID) error
}
