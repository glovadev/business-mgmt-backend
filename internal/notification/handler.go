package notification

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *Handler) SendManualReminder(w http.ResponseWriter, r *http.Request) {
	invoiceIDStr := r.PathValue("invoice_id")
	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid invoice ID format")
		return
	}

	var req models.SendManualReminderRequest
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

	resp, err := h.service.SendManualReminder(r.Context(), invoiceID, req, currentUserID)
	if err != nil {
		if err == ErrValidation || err == ErrInvoiceFullyPaid {
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusAccepted, resp)
}

func (h *Handler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := make(map[string]string)
	if customerID := query.Get("customer_id"); customerID != "" { filters["customer_id"] = customerID }
	if invoiceID := query.Get("invoice_id"); invoiceID != "" { filters["invoice_id"] = invoiceID }
	if triggerType := query.Get("trigger_type"); triggerType != "" { filters["trigger_type"] = triggerType }
	if status := query.Get("status"); status != "" { filters["status"] = status }
	if notifType := query.Get("notification_type"); notifType != "" { filters["notification_type"] = notifType }
	if fromDate := query.Get("from"); fromDate != "" { filters["from"] = fromDate }
	if toDate := query.Get("to"); toDate != "" { filters["to"] = toDate }

	resp, err := h.service.ListNotifications(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetNotification(w http.ResponseWriter, r *http.Request) {
	notificationIDStr := r.PathValue("notification_id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid notification ID format")
		return
	}

	notification, err := h.service.GetNotification(r.Context(), notificationID)
	if err != nil {
		if err == ErrNotificationNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch notification")
		return
	}

	utils.RespondJSON(w, http.StatusOK, notification)
}

func (h *Handler) RetryFailedNotification(w http.ResponseWriter, r *http.Request) {
	notificationIDStr := r.PathValue("notification_id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid notification ID format")
		return
	}

	if err := h.service.RetryFailedNotification(r.Context(), notificationID); err != nil {
		if err == ErrNotificationNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		if err == ErrCannotRetry {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CustomerNotificationHistory(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.PathValue("customer_id")
	if _, err := uuid.Parse(customerIDStr); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid customer ID format")
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := map[string]string{
		"customer_id": customerIDStr,
	}

	resp, err := h.service.ListNotifications(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch customer notifications")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) WhatsAppProviderWebhook(w http.ResponseWriter, r *http.Request) {
	var payload models.WhatsAppWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest) // Generic provider failure response
		return
	}

	signature := r.Header.Get("X-Hub-Signature") // specific to the provider

	// Ignore errors securely, just 200 OK so provider stops retrying.
	_ = h.service.HandleWhatsAppWebhook(r.Context(), payload, signature)
	
	w.WriteHeader(http.StatusOK)
}
