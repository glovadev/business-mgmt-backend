package rbac

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

func (h *Handler) ListRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.ListRoles(r.Context())
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch roles")
		return
	}

	utils.RespondJSON(w, http.StatusOK, roles)
}

func (h *Handler) GetRole(w http.ResponseWriter, r *http.Request) {
	roleIDStr := r.PathValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid role ID format")
		return
	}

	role, err := h.service.GetRole(r.Context(), roleID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, role)
}

func (h *Handler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, err := h.service.ListPermissions(r.Context())
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch permissions")
		return
	}

	utils.RespondJSON(w, http.StatusOK, permissions)
}

func (h *Handler) UpdateRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleIDStr := r.PathValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid role ID format")
		return
	}

	var req models.UpdateRolePermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateRolePermissions(r.Context(), roleID, req); err != nil {
		if err == ErrRoleNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to update role permissions")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
