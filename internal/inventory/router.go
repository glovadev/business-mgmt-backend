package inventory

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Assuming these are all protected by authentication middleware
	
	mux.Handle("GET /api/v1/stores/{store_id}/inventory", authMiddleware(http.HandlerFunc(handler.GetStoreInventory)))
	mux.Handle("GET /api/v1/stores/{store_id}/inventory/{product_id}", authMiddleware(http.HandlerFunc(handler.GetProductInventoryInStore)))
	mux.Handle("GET /api/v1/stores/{store_id}/inventory/{product_id}/movements", authMiddleware(http.HandlerFunc(handler.GetInventoryMovementHistory)))
	mux.Handle("POST /api/v1/stores/{store_id}/inventory/adjustments", authMiddleware(http.HandlerFunc(handler.CreateInventoryAdjustment)))
	mux.Handle("GET /api/v1/stores/{store_id}/inventory/low-stock", authMiddleware(http.HandlerFunc(handler.GetLowStockReport)))

	// Aliases for top-level inventory endpoints
	mux.Handle("GET /api/v1/inventory", authMiddleware(http.HandlerFunc(handler.GetStoreInventory)))
	mux.Handle("POST /api/v1/inventory/adjustments", authMiddleware(http.HandlerFunc(handler.CreateInventoryAdjustment)))
}
