package sale

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

func (h *Handler) CreateSale(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")

	var currentUserID *uuid.UUID
	if val := r.Context().Value("user_id"); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			currentUserID = &id
		}
	}

	sale, err := h.service.CreateSale(r.Context(), req, idempotencyKey, currentUserID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, sale)
}

func (h *Handler) ListSales(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := make(map[string]string)
	if storeID := query.Get("store_id"); storeID != "" { filters["store_id"] = storeID }
	if customerID := query.Get("customer_id"); customerID != "" { filters["customer_id"] = customerID }
	if paymentStatus := query.Get("payment_status"); paymentStatus != "" { filters["payment_status"] = paymentStatus }
	if status := query.Get("status"); status != "" { filters["status"] = status }
	if overdueOnly := query.Get("overdue_only"); overdueOnly != "" { filters["overdue_only"] = overdueOnly }
	if fromDate := query.Get("from_date"); fromDate != "" { filters["from_date"] = fromDate }
	if toDate := query.Get("to_date"); toDate != "" { filters["to_date"] = toDate }

	resp, err := h.service.ListSales(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch sales")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetSale(w http.ResponseWriter, r *http.Request) {
	saleIDStr := r.PathValue("sale_id")
	saleID, err := uuid.Parse(saleIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid sale ID format")
		return
	}

	sale, err := h.service.GetSale(r.Context(), saleID)
	if err != nil {
		if err == ErrInvoiceNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch sale")
		return
	}

	utils.RespondJSON(w, http.StatusOK, sale)
}

func (h *Handler) GetInvoiceByNumber(w http.ResponseWriter, r *http.Request) {
	invoiceNumber := r.PathValue("invoice_number")
	
	sale, err := h.service.GetInvoiceByNumber(r.Context(), invoiceNumber)
	if err != nil {
		if err == ErrInvoiceNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch sale by invoice number")
		return
	}

	utils.RespondJSON(w, http.StatusOK, sale)
}

func (h *Handler) CancelSale(w http.ResponseWriter, r *http.Request) {
	saleIDStr := r.PathValue("sale_id")
	saleID, err := uuid.Parse(saleIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid sale ID format")
		return
	}

	var req models.CancelSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")

	if err := h.service.CancelSale(r.Context(), saleID, req.Reason, idempotencyKey); err != nil {
		utils.RespondError(w, http.StatusConflict, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetSalePayments(w http.ResponseWriter, r *http.Request) {
	saleIDStr := r.PathValue("sale_id")
	saleID, err := uuid.Parse(saleIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid sale ID format")
		return
	}

	payments, err := h.service.GetSalePayments(r.Context(), saleID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch sale payments")
		return
	}

	utils.RespondJSON(w, http.StatusOK, payments)
}
