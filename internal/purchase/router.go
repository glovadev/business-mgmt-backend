package purchase

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Assuming these are all protected by authentication middleware
	
	mux.Handle("POST /api/v1/purchases", authMiddleware(http.HandlerFunc(handler.CreatePurchaseDraft)))
	mux.Handle("GET /api/v1/purchases", authMiddleware(http.HandlerFunc(handler.ListPurchases)))
	mux.Handle("GET /api/v1/purchases/{purchase_id}", authMiddleware(http.HandlerFunc(handler.GetPurchase)))
	mux.Handle("PUT /api/v1/purchases/{purchase_id}", authMiddleware(http.HandlerFunc(handler.UpdatePurchaseDraft)))
	
	// Workflow Actions
	mux.Handle("POST /api/v1/purchases/{purchase_id}/submit", authMiddleware(http.HandlerFunc(handler.SubmitPurchase)))
	mux.Handle("POST /api/v1/purchases/{purchase_id}/approve", authMiddleware(http.HandlerFunc(handler.ApprovePurchase)))
	mux.Handle("POST /api/v1/purchases/{purchase_id}/reject", authMiddleware(http.HandlerFunc(handler.RejectPurchase)))
	mux.Handle("POST /api/v1/purchases/{purchase_id}/cancel", authMiddleware(http.HandlerFunc(handler.CancelPurchase)))
}
