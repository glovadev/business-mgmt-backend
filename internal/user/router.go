package user

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	
	// Create a subrouter or simply register paths. Using Go 1.22 syntax:
	// Assuming these are all protected by authentication middleware
	
	mux.Handle("POST /api/v1/users", authMiddleware(http.HandlerFunc(handler.CreateUser)))
	mux.Handle("GET /api/v1/users", authMiddleware(http.HandlerFunc(handler.ListUsers)))
	mux.Handle("GET /api/v1/users/{user_id}", authMiddleware(http.HandlerFunc(handler.GetUser)))
	mux.Handle("PATCH /api/v1/users/{user_id}", authMiddleware(http.HandlerFunc(handler.UpdateUser)))
	mux.Handle("POST /api/v1/users/{user_id}/activate", authMiddleware(http.HandlerFunc(handler.ActivateUser)))
	mux.Handle("POST /api/v1/users/{user_id}/deactivate", authMiddleware(http.HandlerFunc(handler.DeactivateUser)))
	mux.Handle("PUT /api/v1/users/{user_id}/roles", authMiddleware(http.HandlerFunc(handler.AssignRoles)))
	mux.Handle("PUT /api/v1/users/{user_id}/stores", authMiddleware(http.HandlerFunc(handler.AssignStores)))
}
