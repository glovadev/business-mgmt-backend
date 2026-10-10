package expense

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

// Categories
func (h *Handler) CreateExpenseCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CreateExpenseCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	category, err := h.service.CreateExpenseCategory(r.Context(), req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, category)
}

func (h *Handler) ListExpenseCategories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filters := make(map[string]string)
	if scope := query.Get("scope"); scope != "" { filters["scope"] = scope }
	if status := query.Get("status"); status != "" { filters["status"] = status }

	resp, err := h.service.ListExpenseCategories(r.Context(), filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) UpdateExpenseCategory(w http.ResponseWriter, r *http.Request) {
	catIDStr := r.PathValue("category_id")
	catID, err := uuid.Parse(catIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid category ID format")
		return
	}

	var req models.UpdateExpenseCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateExpenseCategory(r.Context(), catID, req); err != nil {
		if err == ErrCategoryNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeactivateExpenseCategory(w http.ResponseWriter, r *http.Request) {
	catIDStr := r.PathValue("category_id")
	catID, err := uuid.Parse(catIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid category ID format")
		return
	}

	if err := h.service.DeactivateExpenseCategory(r.Context(), catID); err != nil {
		if err == ErrCategoryNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Expenses
func (h *Handler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var req models.CreateExpenseRequest
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

	expense, err := h.service.CreateExpense(r.Context(), req, currentUserID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, expense)
}

func (h *Handler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	filters := make(map[string]string)
	if storeID := query.Get("store_id"); storeID != "" { filters["store_id"] = storeID }
	if categoryID := query.Get("category_id"); categoryID != "" { filters["category_id"] = categoryID }
	if scope := query.Get("scope"); scope != "" { filters["scope"] = scope }
	if status := query.Get("status"); status != "" { filters["status"] = status }
	if createdBy := query.Get("created_by"); createdBy != "" { filters["created_by"] = createdBy }
	if fromDate := query.Get("from_date"); fromDate != "" { filters["from_date"] = fromDate }
	if toDate := query.Get("to_date"); toDate != "" { filters["to_date"] = toDate }

	resp, err := h.service.ListExpenses(r.Context(), page, pageSize, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch expenses")
		return
	}

	utils.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetExpense(w http.ResponseWriter, r *http.Request) {
	expenseIDStr := r.PathValue("expense_id")
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid expense ID format")
		return
	}

	expense, err := h.service.GetExpense(r.Context(), expenseID)
	if err != nil {
		if err == ErrExpenseNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch expense")
		return
	}

	utils.RespondJSON(w, http.StatusOK, expense)
}

func (h *Handler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	expenseIDStr := r.PathValue("expense_id")
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid expense ID format")
		return
	}

	var req models.UpdateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateExpense(r.Context(), expenseID, req); err != nil {
		if err == ErrExpenseNotFound {
			utils.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		if err == ErrCannotEdit {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ReverseExpense(w http.ResponseWriter, r *http.Request) {
	expenseIDStr := r.PathValue("expense_id")
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid expense ID format")
		return
	}

	var req models.ReverseExpenseRequest
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

	if err := h.service.ReverseExpense(r.Context(), expenseID, req.Reason, currentUserID); err != nil {
		if err == ErrCannotReverse {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
