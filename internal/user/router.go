package user

import (
	"net/http"

	"business-mgmt-backend/internal/auth"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	authMiddleware func(http.Handler) http.Handler,
) {
	superAdminOnly := func(next http.Handler) http.Handler {
		return authMiddleware(auth.RequireRole("SUPER_ADMIN")(next))
	}

	mux.Handle("POST /api/v1/users",
		superAdminOnly(http.HandlerFunc(handler.CreateUser)))

	mux.Handle("GET /api/v1/users",
		superAdminOnly(http.HandlerFunc(handler.ListUsers)))

	mux.Handle("GET /api/v1/users/{user_id}",
		superAdminOnly(http.HandlerFunc(handler.GetUser)))

	mux.Handle("PATCH /api/v1/users/{user_id}",
		superAdminOnly(http.HandlerFunc(handler.UpdateUser)))

	mux.Handle("POST /api/v1/users/{user_id}/activate",
		superAdminOnly(http.HandlerFunc(handler.ActivateUser)))

	mux.Handle("POST /api/v1/users/{user_id}/deactivate",
		superAdminOnly(http.HandlerFunc(handler.DeactivateUser)))

	mux.Handle("PUT /api/v1/users/{user_id}/roles",
		superAdminOnly(http.HandlerFunc(handler.AssignRoles)))

	mux.Handle("PUT /api/v1/users/{user_id}/stores",
		superAdminOnly(http.HandlerFunc(handler.AssignStores)))
}
