package notificationtemplate

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

func (h *Handler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNotificationTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	template, err := h.service.CreateTemplate(r.Context(), req)
	if err != nil {
		if err == ErrInvalidInput {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, template)
}

func (h *Handler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.ListTemplates(r.Context())
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch notification templates")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	templateIDStr := r.PathValue("template_id")
	templateID, err := uuid.Parse(templateIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid template ID format")
		return
	}

	var req models.UpdateNotificationTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateTemplate(r.Context(), templateID, req); err != nil {
		if err == ErrTemplateNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		if err == ErrInvalidInput {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
