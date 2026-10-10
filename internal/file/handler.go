package file

import (
	"encoding/json"
	"net/http"

	"business-mgmt-backend/internal/models"
	"business-mgmt-backend/pkg/utils"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUploadURL(w http.ResponseWriter, r *http.Request) {
	var req models.UploadUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := h.service.GetUploadURL(r.Context(), req)
	if err != nil {
		if err == ErrInvalidInput {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate upload URL")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetDownloadURL(w http.ResponseWriter, r *http.Request) {
	objectKey := r.URL.Query().Get("object_key")
	if objectKey == "" {
		utils.RespondError(w, http.StatusBadRequest, "Missing object_key")
		return
	}

	var currentUserID uuid.UUID
	if val := r.Context().Value("user_id"); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			currentUserID = id
		}
	}

	resp, err := h.service.GetDownloadURL(r.Context(), objectKey, currentUserID)
	if err != nil {
		if err == ErrInvalidInput {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err == ErrUnauthorized {
			utils.RespondError(w, http.StatusForbidden, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate download URL")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}
