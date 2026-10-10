package expense

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware func(http.Handler) http.Handler) {
	// Expense Categories
	mux.Handle("POST /api/v1/expense-categories", authMiddleware(http.HandlerFunc(handler.CreateExpenseCategory)))
	mux.Handle("GET /api/v1/expense-categories", authMiddleware(http.HandlerFunc(handler.ListExpenseCategories)))
	mux.Handle("PATCH /api/v1/expense-categories/{category_id}", authMiddleware(http.HandlerFunc(handler.UpdateExpenseCategory)))
	mux.Handle("POST /api/v1/expense-categories/{category_id}/deactivate", authMiddleware(http.HandlerFunc(handler.DeactivateExpenseCategory)))

	// Expenses
	mux.Handle("POST /api/v1/expenses", authMiddleware(http.HandlerFunc(handler.CreateExpense)))
	mux.Handle("GET /api/v1/expenses", authMiddleware(http.HandlerFunc(handler.ListExpenses)))
	mux.Handle("GET /api/v1/expenses/{expense_id}", authMiddleware(http.HandlerFunc(handler.GetExpense)))
	mux.Handle("PATCH /api/v1/expenses/{expense_id}", authMiddleware(http.HandlerFunc(handler.UpdateExpense)))
	mux.Handle("POST /api/v1/expenses/{expense_id}/reverse", authMiddleware(http.HandlerFunc(handler.ReverseExpense)))
}
