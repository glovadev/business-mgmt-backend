package notificationtemplate

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/notification-templates", authMiddleware(http.HandlerFunc(handler.CreateTemplate)))
	mux.Handle("GET /api/v1/notification-templates", authMiddleware(http.HandlerFunc(handler.ListTemplates)))
	mux.Handle("PATCH /api/v1/notification-templates/{template_id}", authMiddleware(http.HandlerFunc(handler.UpdateTemplate)))
}
