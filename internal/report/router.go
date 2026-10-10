package report

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/reports/dashboard", authMiddleware(http.HandlerFunc(handler.GetDashboardSummary)))
	mux.Handle("GET /api/v1/reports/sales", authMiddleware(http.HandlerFunc(handler.GetSalesReport)))
	mux.Handle("GET /api/v1/reports/revenue", authMiddleware(http.HandlerFunc(handler.GetRevenueReport)))
	mux.Handle("GET /api/v1/reports/collections", authMiddleware(http.HandlerFunc(handler.GetCollectionReport)))
	mux.Handle("GET /api/v1/reports/expenses", authMiddleware(http.HandlerFunc(handler.GetExpenseReport)))
	mux.Handle("GET /api/v1/reports/profit", authMiddleware(http.HandlerFunc(handler.GetProfitReport)))
	mux.Handle("GET /api/v1/reports/receivables", authMiddleware(http.HandlerFunc(handler.GetReceivablesReport)))
	mux.Handle("GET /api/v1/reports/receivables/aging", authMiddleware(http.HandlerFunc(handler.GetCreditAgingReport)))
}
