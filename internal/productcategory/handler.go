package productcategory

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

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	category, err := h.service.CreateCategory(r.Context(), req)
	if err != nil {
		if err == ErrCategoryAlreadyExists {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, category)
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.ListCategories(r.Context())
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	utils.RespondJSON(w, http.StatusOK, categories)
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.PathValue("category_id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid category ID format")
		return
	}

	category, err := h.service.GetCategory(r.Context(), categoryID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, category)
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.PathValue("category_id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid category ID format")
		return
	}

	var req models.UpdateProductCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateCategory(r.Context(), categoryID, req); err != nil {
		if err == ErrCategoryAlreadyExists {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to update category")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeactivateCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.PathValue("category_id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid category ID format")
		return
	}

	if err := h.service.DeactivateCategory(r.Context(), categoryID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to deactivate category")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
