package expense

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrCategoryNotFound = errors.New("expense category not found")
	ErrExpenseNotFound  = errors.New("expense not found")
	ErrInvalidInput     = errors.New("invalid input data")
	ErrCannotEdit       = errors.New("expense cannot be edited (prefer reversal)")
	ErrCannotReverse    = errors.New("only POSTED expenses can be reversed")
)

type Service interface {
	CreateExpenseCategory(ctx context.Context, req models.CreateExpenseCategoryRequest) (*models.ExpenseCategoryModel, error)
	ListExpenseCategories(ctx context.Context, filters map[string]string) (*models.ListExpenseCategoriesResponse, error)
	UpdateExpenseCategory(ctx context.Context, id uuid.UUID, req models.UpdateExpenseCategoryRequest) error
	DeactivateExpenseCategory(ctx context.Context, id uuid.UUID) error

	CreateExpense(ctx context.Context, req models.CreateExpenseRequest, currentUserID uuid.UUID) (*models.ExpenseModel, error)
	ListExpenses(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListExpensesResponse, error)
	GetExpense(ctx context.Context, id uuid.UUID) (*models.ExpenseModel, error)
	UpdateExpense(ctx context.Context, id uuid.UUID, req models.UpdateExpenseRequest) error
	ReverseExpense(ctx context.Context, id uuid.UUID, reason string, currentUserID uuid.UUID) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func parseStrAmt(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func (s *serviceImpl) CreateExpenseCategory(ctx context.Context, req models.CreateExpenseCategoryRequest) (*models.ExpenseCategoryModel, error) {
	if req.Name == "" || req.Scope == "" {
		return nil, ErrInvalidInput
	}

	category := &models.ExpenseCategoryModel{
		ID:        uuid.New(),
		Name:      req.Name,
		Scope:     req.Scope,
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateExpenseCategory(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *serviceImpl) ListExpenseCategories(ctx context.Context, filters map[string]string) (*models.ListExpenseCategoriesResponse, error) {
	cats, total, err := s.repo.ListExpenseCategories(ctx, filters)
	if err != nil {
		return nil, err
	}
	return &models.ListExpenseCategoriesResponse{
		Categories: cats,
		TotalCount: total,
	}, nil
}

func (s *serviceImpl) UpdateExpenseCategory(ctx context.Context, id uuid.UUID, req models.UpdateExpenseCategoryRequest) error {
	cat, err := s.repo.GetExpenseCategoryByID(ctx, id)
	if err != nil {
		return ErrCategoryNotFound
	}
	if req.Name != nil { cat.Name = *req.Name }
	if req.Scope != nil { cat.Scope = *req.Scope }

	return s.repo.UpdateExpenseCategory(ctx, cat)
}

func (s *serviceImpl) DeactivateExpenseCategory(ctx context.Context, id uuid.UUID) error {
	cat, err := s.repo.GetExpenseCategoryByID(ctx, id)
	if err != nil {
		return ErrCategoryNotFound
	}
	cat.Status = "INACTIVE"
	return s.repo.UpdateExpenseCategory(ctx, cat)
}

func (s *serviceImpl) CreateExpense(ctx context.Context, req models.CreateExpenseRequest, currentUserID uuid.UUID) (*models.ExpenseModel, error) {
	if req.CategoryID == "" || req.Amount == "" || req.Description == "" {
		return nil, ErrInvalidInput
	}

	categoryID, _ := uuid.Parse(req.CategoryID)
	
	var storeID *uuid.UUID
	if req.StoreID != nil && *req.StoreID != "" {
		id, _ := uuid.Parse(*req.StoreID)
		storeID = &id
	}

	amount := parseStrAmt(req.Amount)
	gstAmount := parseStrAmt(req.GSTAmount)

	expense := &models.ExpenseModel{
		ID:               uuid.New(),
		ExpenseNumber:    fmt.Sprintf("EXP-%d", time.Now().Unix()),
		StoreID:          storeID,
		CategoryID:       categoryID,
		ExpenseDate:      req.ExpenseDate,
		Amount:           formatFloat(amount),
		GSTAmount:        formatFloat(gstAmount),
		PaymentMethod:    req.PaymentMethod,
		Vendor:           req.Vendor,
		Description:      req.Description,
		ReceiptObjectKey: req.ReceiptObjectKey,
		Status:           "POSTED",
		CreatedBy:        currentUserID,
		CreatedAt:        time.Now(),
	}

	if err := s.repo.CreateExpense(ctx, expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *serviceImpl) ListExpenses(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListExpensesResponse, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	expenses, total, err := s.repo.ListExpenses(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListExpensesResponse{
		Expenses:   expenses,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) GetExpense(ctx context.Context, id uuid.UUID) (*models.ExpenseModel, error) {
	return s.repo.GetExpenseByID(ctx, id)
}

func (s *serviceImpl) UpdateExpense(ctx context.Context, id uuid.UUID, req models.UpdateExpenseRequest) error {
	expense, err := s.repo.GetExpenseByID(ctx, id)
	if err != nil {
		return ErrExpenseNotFound
	}
	
	if expense.Status == "REVERSED" {
		return ErrCannotEdit
	}

	if req.CategoryID != nil {
		if cid, err := uuid.Parse(*req.CategoryID); err == nil {
			expense.CategoryID = cid
		}
	}
	if req.ExpenseDate != nil { expense.ExpenseDate = *req.ExpenseDate }
	if req.Amount != nil {
		amount := parseStrAmt(*req.Amount)
		expense.Amount = formatFloat(amount)
	}
	if req.GSTAmount != nil {
		gst := parseStrAmt(*req.GSTAmount)
		expense.GSTAmount = formatFloat(gst)
	}
	if req.PaymentMethod != nil { expense.PaymentMethod = req.PaymentMethod }
	if req.Vendor != nil { expense.Vendor = req.Vendor }
	if req.Description != nil { expense.Description = *req.Description }
	if req.ReceiptObjectKey != nil { expense.ReceiptObjectKey = req.ReceiptObjectKey }

	return s.repo.UpdateExpense(ctx, expense)
}

func (s *serviceImpl) ReverseExpense(ctx context.Context, id uuid.UUID, reason string, currentUserID uuid.UUID) error {
	expense, err := s.repo.GetExpenseByID(ctx, id)
	if err != nil {
		return ErrExpenseNotFound
	}

	if expense.Status != "POSTED" {
		return ErrCannotReverse
	}

	return s.repo.ExecuteExpenseReversalTransaction(ctx, id, reason, currentUserID)
}
