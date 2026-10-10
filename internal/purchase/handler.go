package purchase

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

func (h *Handler) CreatePurchaseDraft(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePurchaseRequest
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

	purchase, err := h.service.CreatePurchaseDraft(r.Context(), req, currentUserID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, purchase)
}

func (h *Handler) ListPurchases(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := make(map[string]string)
	if storeID := query.Get("store_id"); storeID != "" { filters["store_id"] = storeID }
	if supplierID := query.Get("supplier_id"); supplierID != "" { filters["supplier_id"] = supplierID }
	if status := query.Get("status"); status != "" { filters["status"] = status }
	if fromDate := query.Get("from_date"); fromDate != "" { filters["from_date"] = fromDate }
	if toDate := query.Get("to_date"); toDate != "" { filters["to_date"] = toDate }
	if createdBy := query.Get("created_by"); createdBy != "" { filters["created_by"] = createdBy }

	resp, err := h.service.ListPurchases(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch purchases")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetPurchase(w http.ResponseWriter, r *http.Request) {
	purchaseIDStr := r.PathValue("purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid purchase ID format")
		return
	}

	purchase, err := h.service.GetPurchase(r.Context(), purchaseID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, purchase)
}

func (h *Handler) UpdatePurchaseDraft(w http.ResponseWriter, r *http.Request) {
	purchaseIDStr := r.PathValue("purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid purchase ID format")
		return
	}

	var req models.CreatePurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdatePurchaseDraft(r.Context(), purchaseID, req); err != nil {
		if err == ErrInvalidState {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SubmitPurchase(w http.ResponseWriter, r *http.Request) {
	purchaseIDStr := r.PathValue("purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid purchase ID format")
		return
	}

	if err := h.service.SubmitPurchaseForApproval(r.Context(), purchaseID); err != nil {
		if err == ErrInvalidState {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ApprovePurchase(w http.ResponseWriter, r *http.Request) {
	purchaseIDStr := r.PathValue("purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid purchase ID format")
		return
	}
	
	// Idempotency-Key header is expected
	// idempotencyKey := r.Header.Get("Idempotency-Key")

	var currentUserID *uuid.UUID
	if val := r.Context().Value("user_id"); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			currentUserID = &id
		}
	}

	if err := h.service.ApprovePurchase(r.Context(), purchaseID, currentUserID); err != nil {
		if err == ErrSelfApproval || err == ErrInvalidState {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RejectPurchase(w http.ResponseWriter, r *http.Request) {
	purchaseIDStr := r.PathValue("purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid purchase ID format")
		return
	}

	var req models.RejectPurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.RejectPurchase(r.Context(), purchaseID, req.Reason); err != nil {
		if err == ErrInvalidState {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CancelPurchase(w http.ResponseWriter, r *http.Request) {
	purchaseIDStr := r.PathValue("purchase_id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid purchase ID format")
		return
	}

	if err := h.service.CancelPurchase(r.Context(), purchaseID); err != nil {
		if err == ErrInvalidState {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
