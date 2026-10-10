package productcategory

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Assuming these are all protected by authentication middleware
	
	mux.Handle("POST /api/v1/product-categories", authMiddleware(http.HandlerFunc(handler.CreateCategory)))
	mux.Handle("GET /api/v1/product-categories", authMiddleware(http.HandlerFunc(handler.ListCategories)))
	mux.Handle("POST /api/v1/categories", authMiddleware(http.HandlerFunc(handler.CreateCategory)))
	mux.Handle("GET /api/v1/categories", authMiddleware(http.HandlerFunc(handler.ListCategories)))
	mux.Handle("GET /api/v1/product-categories/{category_id}", authMiddleware(http.HandlerFunc(handler.GetCategory)))
	mux.Handle("PATCH /api/v1/product-categories/{category_id}", authMiddleware(http.HandlerFunc(handler.UpdateCategory)))
	mux.Handle("POST /api/v1/product-categories/{category_id}/deactivate", authMiddleware(http.HandlerFunc(handler.DeactivateCategory)))
}
