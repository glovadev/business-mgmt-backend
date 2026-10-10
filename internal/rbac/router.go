package rbac

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Assuming these are all protected by authentication middleware
	
	mux.Handle("GET /api/v1/roles", authMiddleware(http.HandlerFunc(handler.ListRoles)))
	mux.Handle("GET /api/v1/roles/{role_id}", authMiddleware(http.HandlerFunc(handler.GetRole)))
	
	mux.Handle("GET /api/v1/permissions", authMiddleware(http.HandlerFunc(handler.ListPermissions)))
	
	mux.Handle("PUT /api/v1/roles/{role_id}/permissions", authMiddleware(http.HandlerFunc(handler.UpdateRolePermissions)))
}
