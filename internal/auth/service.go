package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"business-mgmt-backend/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserInactive        = errors.New("user is inactive")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrUserNotFound        = errors.New("user not found")
)

type Service interface {
	Login(ctx context.Context, req models.LoginRequest, ipAddress, deviceInfo string) (*models.LoginResponse, error)
	Refresh(ctx context.Context, req models.RefreshRequest) (*models.RefreshResponse, error)
	Logout(ctx context.Context, req models.LogoutRequest) error
	GetMe(ctx context.Context, userID uuid.UUID) (*models.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req models.ChangePasswordRequest) error
	ForgotPassword(ctx context.Context, req models.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req models.ResetPasswordRequest) error
}

type serviceImpl struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewService(userRepo UserRepository, sessionRepo SessionRepository) Service {
	return &serviceImpl{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *serviceImpl) Login(ctx context.Context, req models.LoginRequest, ipAddress, deviceInfo string) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Avoid returning specific user not found error for security
		return nil, ErrInvalidCredentials
	}

	if user.Status != "ACTIVE" {
		return nil, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate access token (e.g. 15 minutes)
	expiresIn := 15 * time.Minute
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Roles, expiresIn)
	if err != nil {
		return nil, err
	}

	// Generate refresh token (e.g. 7 days)
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}
	refreshTokenHash := hashRefreshToken(refreshToken)

	// Create user session
	session := &models.UserSession{
		ID:               uuid.New(),
		UserID:           user.ID,
		RefreshTokenHash: refreshTokenHash,
		DeviceInfo:       deviceInfo,
		IPAddress:        ipAddress,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:        time.Now(),
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(expiresIn.Seconds()),
		User:         *user,
	}, nil
}

func (s *serviceImpl) Refresh(ctx context.Context, req models.RefreshRequest) (*models.RefreshResponse, error) {
	tokenHash := hashRefreshToken(req.RefreshToken)

	session, err := s.sessionRepo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil || session.RevokedAt != nil || session.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.GetUserByID(ctx, session.UserID)
	if err != nil || user.Status != "ACTIVE" {
		return nil, ErrUserInactive
	}

	// Rotate refresh token
	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}
	newRefreshTokenHash := hashRefreshToken(newRefreshToken)

	// Revoke old session
	if err := s.sessionRepo.RevokeSession(ctx, session.ID); err != nil {
		return nil, err
	}

	// Create new session
	newSession := &models.UserSession{
		ID:               uuid.New(),
		UserID:           user.ID,
		RefreshTokenHash: newRefreshTokenHash,
		DeviceInfo:       session.DeviceInfo, // Carry over info
		IPAddress:        session.IPAddress,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:        time.Now(),
	}

	if err := s.sessionRepo.CreateSession(ctx, newSession); err != nil {
		return nil, err
	}

	expiresIn := 15 * time.Minute
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Roles, expiresIn)
	if err != nil {
		return nil, err
	}

	return &models.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(expiresIn.Seconds()),
	}, nil
}

func (s *serviceImpl) Logout(ctx context.Context, req models.LogoutRequest) error {
	tokenHash := hashRefreshToken(req.RefreshToken)

	session, err := s.sessionRepo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		// Already logged out or invalid, consider it successful
		return nil
	}

	return s.sessionRepo.RevokeSession(ctx, session.ID)
}

func (s *serviceImpl) GetMe(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *serviceImpl) ChangePassword(ctx context.Context, userID uuid.UUID, req models.ChangePasswordRequest) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.New("incorrect current password")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, user.ID, string(newHash))
}

func (s *serviceImpl) ForgotPassword(ctx context.Context, req models.ForgotPasswordRequest) error {
	// Generic response to avoid account enumeration
	return nil
}

func (s *serviceImpl) ResetPassword(ctx context.Context, req models.ResetPasswordRequest) error {
	return nil
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
