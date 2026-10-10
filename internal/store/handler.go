package store

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

func (h *Handler) CreateStore(w http.ResponseWriter, r *http.Request) {
	var req models.CreateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var currentUserID *uuid.UUID
	if val := r.Context().Value("user_id"); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			currentUserID = &id
		}
	}

	store, err := h.service.CreateStore(r.Context(), req, currentUserID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, store)
}

func (h *Handler) ListStores(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := make(map[string]string)
	if status := query.Get("status"); status != "" {
		filters["status"] = status
	}
	if q := query.Get("q"); q != "" {
		filters["q"] = q
	}

	resp, err := h.service.ListStores(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch stores")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetStore(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.PathValue("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid store ID format")
		return
	}

	store, err := h.service.GetStore(r.Context(), storeID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, store)
}

func (h *Handler) UpdateStore(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.PathValue("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid store ID format")
		return
	}

	var req models.UpdateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateStore(r.Context(), storeID, req); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to update store")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ActivateStore(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.PathValue("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid store ID format")
		return
	}

	if err := h.service.ActivateStore(r.Context(), storeID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to activate store")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeactivateStore(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.PathValue("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid store ID format")
		return
	}

	if err := h.service.DeactivateStore(r.Context(), storeID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to deactivate store")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
