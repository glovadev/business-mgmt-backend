package notification

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/invoices/{invoice_id}/reminders", authMiddleware(http.HandlerFunc(handler.SendManualReminder)))
	mux.Handle("GET /api/v1/notifications", authMiddleware(http.HandlerFunc(handler.ListNotifications)))
	mux.Handle("GET /api/v1/notifications/{notification_id}", authMiddleware(http.HandlerFunc(handler.GetNotification)))
	mux.Handle("POST /api/v1/notifications/{notification_id}/retry", authMiddleware(http.HandlerFunc(handler.RetryFailedNotification)))
	
	// Sugar endpoint for customer notifications
	mux.Handle("GET /api/v1/customers/{customer_id}/notifications", authMiddleware(http.HandlerFunc(handler.CustomerNotificationHistory)))

	// Webhook endpoint (unauthenticated via normal means)
	mux.HandleFunc("POST /api/v1/webhooks/whatsapp", handler.WhatsAppProviderWebhook)
}
