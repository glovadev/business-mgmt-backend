package inventory

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

func (h *Handler) GetStoreInventory(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.PathValue("store_id")
	if storeIDStr == "" {
		storeIDStr = r.URL.Query().Get("store_id")
	}
	if storeIDStr == "" {
		storeIDStr = "00000000-0000-0000-0000-000000000001"
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid store ID format")
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := make(map[string]string)
	if q := query.Get("q"); q != "" { filters["q"] = q }
	if productID := query.Get("product_id"); productID != "" { filters["product_id"] = productID }
	if categoryID := query.Get("category_id"); categoryID != "" { filters["category_id"] = categoryID }
	if lowStock := query.Get("low_stock"); lowStock != "" { filters["low_stock"] = lowStock }

	resp, err := h.service.GetStoreInventory(r.Context(), storeID, page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch store inventory")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetProductInventoryInStore(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.PathValue("store_id")
	productIDStr := r.PathValue("product_id")
	
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid store ID format")
		return
	}
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	inventory, err := h.service.GetProductInventoryInStore(r.Context(), storeID, productID)
	if err != nil {
		if err == ErrInventoryNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch product inventory")
		return
	}

	utils.RespondJSON(w, http.StatusOK, inventory)
}

func (h *Handler) GetInventoryMovementHistory(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.PathValue("store_id")
	productIDStr := r.PathValue("product_id")
	
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid store ID format")
		return
	}
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := make(map[string]string)
	if movementType := query.Get("movement_type"); movementType != "" { filters["movement_type"] = movementType }
	if from := query.Get("from"); from != "" { filters["from"] = from }
	if to := query.Get("to"); to != "" { filters["to"] = to }

	resp, err := h.service.GetInventoryMovementHistory(r.Context(), storeID, productID, page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch inventory movements")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) CreateInventoryAdjustment(w http.ResponseWriter, r *http.Request) {
	var req models.CreateInventoryAdjustmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	storeIDStr := r.PathValue("store_id")
	if storeIDStr == "" && req.StoreID != nil {
		storeIDStr = *req.StoreID
	}
	if storeIDStr == "" {
		storeIDStr = "00000000-0000-0000-0000-000000000001"
	}

	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid store ID format")
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")

	var currentUserID uuid.UUID
	if val := r.Context().Value("user_id"); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			currentUserID = id
		}
	}

	if err := h.service.CreateInventoryAdjustment(r.Context(), storeID, req, idempotencyKey, currentUserID); err != nil {
		if err == ErrInsufficientStock || err == ErrInvalidAdjustment {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetLowStockReport(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.PathValue("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid store ID format")
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))
	threshold := query.Get("threshold")

	resp, err := h.service.GetLowStockReport(r.Context(), storeID, threshold, page, pageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch low stock report")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}
