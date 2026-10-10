package customer

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Assuming these are all protected by authentication middleware
	
	mux.Handle("POST /api/v1/customers", authMiddleware(http.HandlerFunc(handler.CreateCustomer)))
	mux.Handle("GET /api/v1/customers", authMiddleware(http.HandlerFunc(handler.ListCustomers)))
	mux.Handle("GET /api/v1/customers/{customer_id}", authMiddleware(http.HandlerFunc(handler.GetCustomer)))
	mux.Handle("PATCH /api/v1/customers/{customer_id}", authMiddleware(http.HandlerFunc(handler.UpdateCustomer)))
	mux.Handle("POST /api/v1/customers/{customer_id}/activate", authMiddleware(http.HandlerFunc(handler.ActivateCustomer)))
	mux.Handle("POST /api/v1/customers/{customer_id}/deactivate", authMiddleware(http.HandlerFunc(handler.DeactivateCustomer)))
	
	// Customer Store Relationships
	mux.Handle("GET /api/v1/customers/{customer_id}/stores", authMiddleware(http.HandlerFunc(handler.GetCustomerStores)))
	mux.Handle("POST /api/v1/customers/{customer_id}/stores/{store_id}", authMiddleware(http.HandlerFunc(handler.LinkCustomerToStore)))

	// WhatsApp Opt-in / Opt-out
	mux.Handle("POST /api/v1/customers/{customer_id}/whatsapp-opt-in", authMiddleware(http.HandlerFunc(handler.OptInWhatsApp)))
	mux.Handle("POST /api/v1/customers/{customer_id}/whatsapp-opt-out", authMiddleware(http.HandlerFunc(handler.OptOutWhatsApp)))

	// Customer Ledger, Payments, Outstanding
	mux.Handle("GET /api/v1/customers/{customer_id}/ledger", authMiddleware(http.HandlerFunc(handler.GetCustomerLedger)))
	mux.Handle("GET /api/v1/customers/{customer_id}/payments", authMiddleware(http.HandlerFunc(handler.GetCustomerPayments)))
	mux.Handle("GET /api/v1/customers/{customer_id}/outstanding", authMiddleware(http.HandlerFunc(handler.GetCustomerOutstanding)))
}
