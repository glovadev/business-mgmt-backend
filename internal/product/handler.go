package product

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

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest
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

	product, err := h.service.CreateProduct(r.Context(), req, currentUserID)
	if err != nil {
		if err == ErrProductAlreadyExists {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, product)
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
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
	if sku := query.Get("sku"); sku != "" {
		filters["sku"] = sku
	}
	if categoryID := query.Get("category_id"); categoryID != "" {
		filters["category_id"] = categoryID
	}

	resp, err := h.service.ListProducts(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	product, err := h.service.GetProduct(r.Context(), productID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, product)
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateProduct(r.Context(), productID, req); err != nil {
		if err == ErrProductAlreadyExists {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ActivateProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	if err := h.service.ActivateProduct(r.Context(), productID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to activate product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeactivateProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	if err := h.service.DeactivateProduct(r.Context(), productID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to deactivate product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetProductInventory(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	resp, err := h.service.GetProductInventory(r.Context(), productID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch product inventory")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetPriceHistory(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"data": []interface{}{},
		"total_count": 0,
		"page": 1,
		"page_size": 20,
	})
}
