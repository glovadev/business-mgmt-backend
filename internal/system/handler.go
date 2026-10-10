package system

import (
	"net/http"

	"business-mgmt-backend/pkg/utils"
)

type Handler struct {
	// Database connection or dependency checkers would be injected here for Readiness
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	// E.g., err := h.db.PingContext(r.Context())
	// if err != nil { w.WriteHeader(http.StatusServiceUnavailable); return }
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"version": "1.0.0", // This could be populated by build flags
		"commit":  "unknown",
	})
}
