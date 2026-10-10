package configuration

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

func (h *Handler) GetCompanySettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.service.GetCompanySettings(r.Context())
	if err != nil {
		if err == ErrSettingsNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch company settings")
		return
	}

	resp := models.GetCompanySettingsResponse{Data: *settings}
	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) UpdateCompanySettings(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateCompanySettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var currentUserID uuid.UUID
	if val := r.Context().Value("user_id"); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			currentUserID = id
		}
	}

	if err := h.service.UpdateCompanySettings(r.Context(), req, currentUserID); err != nil {
		if err == ErrInvalidPolicy {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err == ErrSettingsNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
