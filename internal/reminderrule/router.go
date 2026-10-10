package reminderrule

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/reminder-rules", authMiddleware(http.HandlerFunc(handler.CreateRule)))
	mux.Handle("GET /api/v1/reminder-rules", authMiddleware(http.HandlerFunc(handler.ListRules)))
	mux.Handle("PATCH /api/v1/reminder-rules/{rule_id}", authMiddleware(http.HandlerFunc(handler.UpdateRule)))
	mux.Handle("POST /api/v1/reminder-rules/{rule_id}/activate", authMiddleware(http.HandlerFunc(handler.ActivateRule)))
	mux.Handle("POST /api/v1/reminder-rules/{rule_id}/deactivate", authMiddleware(http.HandlerFunc(handler.DeactivateRule)))
	mux.Handle("PUT /api/v1/reminder-rules/order", authMiddleware(http.HandlerFunc(handler.ReorderRules)))
}
