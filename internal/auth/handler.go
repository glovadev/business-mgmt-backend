package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	ipAddress := r.RemoteAddr
	deviceInfo := r.UserAgent()

	resp, err := h.service.Login(r.Context(), req, ipAddress, deviceInfo)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := h.service.Refresh(r.Context(), req)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req models.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.Logout(r.Context(), req); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	utils.RespondNoContent(w)
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Assuming a middleware extracts the user ID and sets it in the request context
	userIDValue := r.Context().Value("user_id")
	if userIDValue == nil {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid user ID in context")
		return
	}

	user, err := h.service.GetMe(r.Context(), userID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, user)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value("user_id")
	if userIDValue == nil {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid user ID in context")
		return
	}

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.ChangePassword(r.Context(), userID, req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondNoContent(w)
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	_ = h.service.ForgotPassword(r.Context(), req)
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "If the email is registered, a password reset link has been sent."})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.ResetPassword(r.Context(), req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Password reset successfully."})
}

// AuthMiddleware extracts and validates the JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		tokenString := parts[1]
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Add UserID to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "roles", claims.Roles)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
