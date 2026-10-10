package supplier

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Assuming these are all protected by authentication middleware
	
	mux.Handle("POST /api/v1/suppliers", authMiddleware(http.HandlerFunc(handler.CreateSupplier)))
	mux.Handle("GET /api/v1/suppliers", authMiddleware(http.HandlerFunc(handler.ListSuppliers)))
	mux.Handle("GET /api/v1/suppliers/{supplier_id}", authMiddleware(http.HandlerFunc(handler.GetSupplier)))
	mux.Handle("PATCH /api/v1/suppliers/{supplier_id}", authMiddleware(http.HandlerFunc(handler.UpdateSupplier)))
	mux.Handle("POST /api/v1/suppliers/{supplier_id}/activate", authMiddleware(http.HandlerFunc(handler.ActivateSupplier)))
	mux.Handle("POST /api/v1/suppliers/{supplier_id}/deactivate", authMiddleware(http.HandlerFunc(handler.DeactivateSupplier)))
}
