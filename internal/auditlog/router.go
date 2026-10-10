package auditlog

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/audit-logs", authMiddleware(http.HandlerFunc(handler.ListAuditLogs)))
	mux.Handle("GET /api/v1/audit-logs/{audit_log_id}", authMiddleware(http.HandlerFunc(handler.GetAuditLog)))
	
	// Entity Audit History
	mux.Handle("GET /api/v1/audit-logs/entity/{entity_type}/{entity_id}", authMiddleware(http.HandlerFunc(handler.GetEntityAuditHistory)))
}
