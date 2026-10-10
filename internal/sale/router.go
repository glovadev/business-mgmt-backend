package sale

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Assuming these are all protected by authentication middleware
	
	mux.Handle("POST /api/v1/sales", authMiddleware(http.HandlerFunc(handler.CreateSale)))
	mux.Handle("GET /api/v1/sales", authMiddleware(http.HandlerFunc(handler.ListSales)))
	mux.Handle("GET /api/v1/sales/{sale_id}", authMiddleware(http.HandlerFunc(handler.GetSale)))
	mux.Handle("GET /api/v1/invoices/{invoice_number}", authMiddleware(http.HandlerFunc(handler.GetInvoiceByNumber)))
	
	mux.Handle("POST /api/v1/sales/{sale_id}/cancel", authMiddleware(http.HandlerFunc(handler.CancelSale)))
	mux.Handle("GET /api/v1/sales/{sale_id}/payments", authMiddleware(http.HandlerFunc(handler.GetSalePayments)))
}
