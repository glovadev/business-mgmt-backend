package product

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Assuming these are all protected by authentication middleware
	
	mux.Handle("POST /api/v1/products", authMiddleware(http.HandlerFunc(handler.CreateProduct)))
	mux.Handle("GET /api/v1/products", authMiddleware(http.HandlerFunc(handler.ListProducts)))
	mux.Handle("GET /api/v1/products/{product_id}", authMiddleware(http.HandlerFunc(handler.GetProduct)))
	mux.Handle("PATCH /api/v1/products/{product_id}", authMiddleware(http.HandlerFunc(handler.UpdateProduct)))
	mux.Handle("POST /api/v1/products/{product_id}/activate", authMiddleware(http.HandlerFunc(handler.ActivateProduct)))
	mux.Handle("POST /api/v1/products/{product_id}/deactivate", authMiddleware(http.HandlerFunc(handler.DeactivateProduct)))
	
	// Inventory and Price history
	mux.Handle("GET /api/v1/products/{product_id}/inventory", authMiddleware(http.HandlerFunc(handler.GetProductInventory)))
	mux.Handle("GET /api/v1/products/{product_id}/price-history", authMiddleware(http.HandlerFunc(handler.GetPriceHistory)))
}
