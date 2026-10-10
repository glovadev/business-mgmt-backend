package auth

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	// Public routes
	mux.HandleFunc("POST /api/v1/auth/login", handler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", handler.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", handler.Logout)

	// Protected routes using Go 1.22 mux pattern
	meHandler := http.HandlerFunc(handler.GetMe)
	mux.Handle("GET /api/v1/auth/me", AuthMiddleware(meHandler))

	changePwdHandler := http.HandlerFunc(handler.ChangePassword)
	mux.Handle("POST /api/v1/auth/change-password", AuthMiddleware(changePwdHandler))

	mux.HandleFunc("POST /api/v1/auth/forgot-password", handler.ForgotPassword)
	mux.HandleFunc("POST /api/v1/auth/reset-password", handler.ResetPassword)
}
