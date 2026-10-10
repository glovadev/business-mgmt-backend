package auditlog

import (
	"net/http"
	"strconv"

	"business-mgmt-backend/pkg/utils"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := make(map[string]string)
	if actor := query.Get("actor_user_id"); actor != "" { filters["actor_user_id"] = actor }
	if store := query.Get("store_id"); store != "" { filters["store_id"] = store }
	if action := query.Get("action"); action != "" { filters["action"] = action }
	if eType := query.Get("entity_type"); eType != "" { filters["entity_type"] = eType }
	if eID := query.Get("entity_id"); eID != "" { filters["entity_id"] = eID }
	if from := query.Get("from"); from != "" { filters["from"] = from }
	if to := query.Get("to"); to != "" { filters["to"] = to }

	resp, err := h.service.ListAuditLogs(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch audit logs")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetAuditLog(w http.ResponseWriter, r *http.Request) {
	logIDStr := r.PathValue("audit_log_id")
	logID, err := uuid.Parse(logIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid audit log ID format")
		return
	}

	log, err := h.service.GetAuditLog(r.Context(), logID)
	if err != nil {
		// Treating nil return as 404 for simplicity if repo supports it
		utils.RespondError(w, http.StatusNotFound, "Audit log not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, log)
}

func (h *Handler) GetEntityAuditHistory(w http.ResponseWriter, r *http.Request) {
	eType := r.PathValue("entity_type")
	eIDStr := r.PathValue("entity_id")
	
	eID, err := uuid.Parse(eIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid entity ID format")
		return
	}

	history, err := h.service.GetEntityAuditHistory(r.Context(), eType, eID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch entity audit history")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"data": history,
	})
}
