package file

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
	ErrUnauthorized = errors.New("unauthorized to access this object")
)

type Service interface {
	GetUploadURL(ctx context.Context, req models.UploadUrlRequest) (*models.UploadUrlResponse, error)
	GetDownloadURL(ctx context.Context, objectKey string, currentUserID uuid.UUID) (*models.DownloadUrlResponse, error)
}

type serviceImpl struct {
	// StorageProvider (AWS S3, GCP Cloud Storage) would be injected here
}

func NewService() Service {
	return &serviceImpl{}
}

func (s *serviceImpl) GetUploadURL(ctx context.Context, req models.UploadUrlRequest) (*models.UploadUrlResponse, error) {
	if req.FileName == "" || req.Purpose == "" {
		return nil, ErrInvalidInput
	}

	purposePrefix := strings.ToLower(req.Purpose)
	objectKey := fmt.Sprintf("%s/%s/%s", purposePrefix, uuid.New().String(), req.FileName)

	// Mock signed URL generation
	signedURL := fmt.Sprintf("https://storage.example.com/upload?key=%s&sig=mock-signature", objectKey)

	resp := &models.UploadUrlResponse{}
	resp.Data.UploadURL = signedURL
	resp.Data.ObjectKey = objectKey
	resp.Data.ExpiresIn = 300 // 5 minutes

	return resp, nil
}

func (s *serviceImpl) GetDownloadURL(ctx context.Context, objectKey string, currentUserID uuid.UUID) (*models.DownloadUrlResponse, error) {
	if objectKey == "" {
		return nil, ErrInvalidInput
	}

	// In a real application, we would check if the current user has access to the underlying record 
	// (e.g. check if the objectKey matches a receipt they are authorized to view)
	// For V1, we simulate passing the authorization check.

	// Mock signed URL generation
	signedURL := fmt.Sprintf("https://storage.example.com/download?key=%s&sig=mock-signature", objectKey)

	resp := &models.DownloadUrlResponse{}
	resp.Data.DownloadURL = signedURL
	resp.Data.ExpiresIn = 300 // 5 minutes

	return resp, nil
}
