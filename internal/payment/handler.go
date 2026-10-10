package payment

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

func (h *Handler) CollectCustomerPayment(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")

	var currentUserID uuid.UUID
	if val := r.Context().Value("user_id"); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			currentUserID = id
		}
	}

	payment, err := h.service.CollectCustomerPayment(r.Context(), req, idempotencyKey, currentUserID)
	if err != nil {
		if err == ErrInvalidAllocation || err == ErrInvoiceNotFound || err == ErrInvalidInput {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, payment)
}

func (h *Handler) ListPayments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := make(map[string]string)
	if customerID := query.Get("customer_id"); customerID != "" { filters["customer_id"] = customerID }
	if invoiceID := query.Get("invoice_id"); invoiceID != "" { filters["invoice_id"] = invoiceID }
	if collectionStoreID := query.Get("collection_store_id"); collectionStoreID != "" { filters["collection_store_id"] = collectionStoreID }
	if collectedBy := query.Get("collected_by"); collectedBy != "" { filters["collected_by"] = collectedBy }
	if paymentMethod := query.Get("payment_method"); paymentMethod != "" { filters["payment_method"] = paymentMethod }
	if status := query.Get("status"); status != "" { filters["status"] = status }
	if fromDate := query.Get("from"); fromDate != "" { filters["from"] = fromDate }
	if toDate := query.Get("to"); toDate != "" { filters["to"] = toDate }

	resp, err := h.service.ListPayments(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch payments")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetPayment(w http.ResponseWriter, r *http.Request) {
	paymentIDStr := r.PathValue("payment_id")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid payment ID format")
		return
	}

	payment, err := h.service.GetPayment(r.Context(), paymentID)
	if err != nil {
		if err == ErrPaymentNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch payment")
		return
	}

	utils.RespondJSON(w, http.StatusOK, payment)
}

func (h *Handler) ReversePayment(w http.ResponseWriter, r *http.Request) {
	paymentIDStr := r.PathValue("payment_id")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid payment ID format")
		return
	}

	var req models.ReversePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	
	var currentUserID uuid.UUID
	if val := r.Context().Value("user_id"); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			currentUserID = id
		}
	}

	if err := h.service.ReversePayment(r.Context(), paymentID, req.Reason, idempotencyKey, currentUserID); err != nil {
		if err == ErrPaymentNotPosted {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
