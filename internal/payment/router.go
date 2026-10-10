package payment

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Assuming these are all protected by authentication middleware
	
	mux.Handle("POST /api/v1/payments", authMiddleware(http.HandlerFunc(handler.CollectCustomerPayment)))
	mux.Handle("GET /api/v1/payments", authMiddleware(http.HandlerFunc(handler.ListPayments)))
	mux.Handle("GET /api/v1/payments/{payment_id}", authMiddleware(http.HandlerFunc(handler.GetPayment)))
	mux.Handle("POST /api/v1/payments/{payment_id}/reverse", authMiddleware(http.HandlerFunc(handler.ReversePayment)))
}
