package customer

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

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCustomerRequest
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

	customer, err := h.service.CreateCustomer(r.Context(), req, currentUserID)
	if err != nil {
		if err == ErrCustomerAlreadyExists {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, customer)
}

func (h *Handler) ListCustomers(w http.ResponseWriter, r *http.Request) {
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
	if phone := query.Get("phone"); phone != "" {
		filters["phone"] = phone
	}
	if storeID := query.Get("store_id"); storeID != "" {
		filters["store_id"] = storeID
	}
	if riskStatus := query.Get("risk_status"); riskStatus != "" {
		filters["risk_status"] = riskStatus
	}
	if hasOutstanding := query.Get("has_outstanding"); hasOutstanding != "" {
		filters["has_outstanding"] = hasOutstanding
	}

	resp, err := h.service.ListCustomers(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch customers")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.PathValue("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid customer ID format")
		return
	}

	customer, err := h.service.GetCustomer(r.Context(), customerID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	// Wrapper for single resource response data according to API spec suggested response
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{"data": customer})
}

func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.PathValue("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid customer ID format")
		return
	}

	var req models.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateCustomer(r.Context(), customerID, req); err != nil {
		if err == ErrCustomerAlreadyExists {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to update customer")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ActivateCustomer(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.PathValue("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid customer ID format")
		return
	}

	if err := h.service.ActivateCustomer(r.Context(), customerID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to activate customer")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeactivateCustomer(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.PathValue("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid customer ID format")
		return
	}

	if err := h.service.DeactivateCustomer(r.Context(), customerID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to deactivate customer")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCustomerStores(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.PathValue("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid customer ID format")
		return
	}

	stores, err := h.service.GetCustomerStores(r.Context(), customerID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch customer stores")
		return
	}

	utils.RespondJSON(w, http.StatusOK, stores)
}

func (h *Handler) LinkCustomerToStore(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.PathValue("customer_id")
	storeIDStr := r.PathValue("store_id")
	
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid customer ID format")
		return
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid store ID format")
		return
	}

	var currentUserID *uuid.UUID
	if val := r.Context().Value("user_id"); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			currentUserID = &id
		}
	}

	if err := h.service.LinkCustomerToStore(r.Context(), customerID, storeID, currentUserID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to link customer to store")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) OptInWhatsApp(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.PathValue("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid customer ID format")
		return
	}
	optIn := true
	_ = h.service.UpdateCustomer(r.Context(), customerID, models.UpdateCustomerRequest{WhatsappOptIn: &optIn})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) OptOutWhatsApp(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.PathValue("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid customer ID format")
		return
	}
	optIn := false
	_ = h.service.UpdateCustomer(r.Context(), customerID, models.UpdateCustomerRequest{WhatsappOptIn: &optIn})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCustomerLedger(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"data": []interface{}{},
		"total_count": 0,
		"page": 1,
		"page_size": 20,
	})
}

func (h *Handler) GetCustomerPayments(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"data": []interface{}{},
		"total_count": 0,
		"page": 1,
		"page_size": 20,
	})
}

func (h *Handler) GetCustomerOutstanding(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"total_outstanding": "0.00",
		"oldest_due_date": nil,
		"max_overdue_days": 0,
		"risk_status": "NORMAL",
	})
}
