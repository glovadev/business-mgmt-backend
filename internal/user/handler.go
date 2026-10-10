package user

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.CreateUser(r.Context(), req); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := make(map[string]string)
	if status := query.Get("status"); status != "" {
		filters["status"] = status
	}
	if role := query.Get("role"); role != "" {
		filters["role"] = role
	}
	if storeID := query.Get("store_id"); storeID != "" {
		filters["store_id"] = storeID
	}
	if q := query.Get("q"); q != "" {
		filters["q"] = q
	}

	resp, err := h.service.ListUsers(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateUser(r.Context(), userID, req); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	if err := h.service.ActivateUser(r.Context(), userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to activate user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	if err := h.service.DeactivateUser(r.Context(), userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to deactivate user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AssignRoles(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var req models.AssignRolesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.AssignRoles(r.Context(), userID, req); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to assign roles")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AssignStores(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var req models.AssignStoresRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.AssignStores(r.Context(), userID, req); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to assign stores")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
