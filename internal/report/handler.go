package report

import (
	"net/http"

	"business-mgmt-backend/pkg/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func extractFilters(r *http.Request) map[string]string {
	query := r.URL.Query()
	filters := make(map[string]string)
	
	if storeID := query.Get("store_id"); storeID != "" { filters["store_id"] = storeID }
	if fromDate := query.Get("from_date"); fromDate != "" { filters["from_date"] = fromDate }
	if toDate := query.Get("to_date"); toDate != "" { filters["to_date"] = toDate }
	if groupBy := query.Get("group_by"); groupBy != "" { filters["group_by"] = groupBy }
	if customerID := query.Get("customer_id"); customerID != "" { filters["customer_id"] = customerID }
	if categoryID := query.Get("category_id"); categoryID != "" { filters["category_id"] = categoryID }
	if collectionStoreID := query.Get("collection_store_id"); collectionStoreID != "" { filters["collection_store_id"] = collectionStoreID }
	if paymentMethod := query.Get("payment_method"); paymentMethod != "" { filters["payment_method"] = paymentMethod }
	if scope := query.Get("scope"); scope != "" { filters["scope"] = scope }
	if riskStatus := query.Get("risk_status"); riskStatus != "" { filters["risk_status"] = riskStatus }
	if overdueOnly := query.Get("overdue_only"); overdueOnly != "" { filters["overdue_only"] = overdueOnly }

	return filters
}

func (h *Handler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetDashboardSummary(r.Context(), extractFilters(r))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch dashboard summary")
		return
	}
	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetSalesReport(r.Context(), extractFilters(r))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch sales report")
		return
	}
	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetRevenueReport(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetRevenueReport(r.Context(), extractFilters(r))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch revenue report")
		return
	}
	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetCollectionReport(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetCollectionReport(r.Context(), extractFilters(r))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch collection report")
		return
	}
	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetExpenseReport(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetExpenseReport(r.Context(), extractFilters(r))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch expense report")
		return
	}
	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetProfitReport(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetProfitReport(r.Context(), extractFilters(r))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch profit report")
		return
	}
	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetReceivablesReport(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetReceivablesReport(r.Context(), extractFilters(r))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch receivables report")
		return
	}
	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetCreditAgingReport(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetCreditAgingReport(r.Context(), extractFilters(r))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch credit aging report")
		return
	}
	utils.RespondJSON(w, http.StatusOK, resp)
}
