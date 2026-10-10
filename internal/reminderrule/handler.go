package reminderrule

import (
	"encoding/json"
	"net/http"

	"business-mgmt-backend/internal/models"
	"business-mgmt-backend/pkg/utils"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateRule(w http.ResponseWriter, r *http.Request) {
	var req models.CreateReminderRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var currentUserID uuid.UUID
	if val := r.Context().Value("user_id"); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			currentUserID = id
		}
	}

	rule, err := h.service.CreateRule(r.Context(), req, currentUserID)
	if err != nil {
		if err == ErrInvalidInput {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, rule)
}

func (h *Handler) ListRules(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.ListRules(r.Context())
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch reminder rules")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	ruleIDStr := r.PathValue("rule_id")
	ruleID, err := uuid.Parse(ruleIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid rule ID format")
		return
	}

	var req models.UpdateReminderRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateRule(r.Context(), ruleID, req); err != nil {
		if err == ErrRuleNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		if err == ErrInvalidInput {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ActivateRule(w http.ResponseWriter, r *http.Request) {
	ruleIDStr := r.PathValue("rule_id")
	ruleID, err := uuid.Parse(ruleIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid rule ID format")
		return
	}

	if err := h.service.ActivateRule(r.Context(), ruleID); err != nil {
		if err == ErrRuleNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeactivateRule(w http.ResponseWriter, r *http.Request) {
	ruleIDStr := r.PathValue("rule_id")
	ruleID, err := uuid.Parse(ruleIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid rule ID format")
		return
	}

	if err := h.service.DeactivateRule(r.Context(), ruleID); err != nil {
		if err == ErrRuleNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ReorderRules(w http.ResponseWriter, r *http.Request) {
	var req models.ReorderReminderRulesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.ReorderRules(r.Context(), req); err != nil {
		if err == ErrInvalidInput {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		// In a real application, catching a Postgres UNIQUE constraint violation on 'sequence_number'
		// would map to a 409 Conflict here. We abstract it away for now.
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
