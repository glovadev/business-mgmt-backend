package auth

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.UserSession) error
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (*models.UserSession, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
}
