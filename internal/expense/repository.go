package expense

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	// Expense Categories
	CreateExpenseCategory(ctx context.Context, category *models.ExpenseCategoryModel) error
	GetExpenseCategoryByID(ctx context.Context, id uuid.UUID) (*models.ExpenseCategoryModel, error)
	ListExpenseCategories(ctx context.Context, filters map[string]string) ([]models.ExpenseCategoryModel, int, error)
	UpdateExpenseCategory(ctx context.Context, category *models.ExpenseCategoryModel) error
	
	// Expenses
	CreateExpense(ctx context.Context, expense *models.ExpenseModel) error
	GetExpenseByID(ctx context.Context, id uuid.UUID) (*models.ExpenseModel, error)
	ListExpenses(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.ExpenseModel, int, error)
	UpdateExpense(ctx context.Context, expense *models.ExpenseModel) error
	
	// Atomic Transactions
	ExecuteExpenseReversalTransaction(ctx context.Context, expenseID uuid.UUID, reason string, currentUserID uuid.UUID) error
}
