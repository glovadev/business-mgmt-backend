package models

import (
	"time"

	"github.com/google/uuid"
)

type ExpenseCategoryModel struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Scope     string    `json:"scope"` // STORE, COMPANY, BOTH
	Status    string    `json:"status"` // ACTIVE, INACTIVE
	CreatedAt time.Time `json:"created_at"`
}

type CreateExpenseCategoryRequest struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
}

type UpdateExpenseCategoryRequest struct {
	Name  *string `json:"name,omitempty"`
	Scope *string `json:"scope,omitempty"`
}

type ExpenseModel struct {
	ID                  uuid.UUID  `json:"id"`
	ExpenseNumber       string     `json:"expense_number"`
	StoreID             *uuid.UUID `json:"store_id,omitempty"`
	CategoryID          uuid.UUID  `json:"category_id"`
	ExpenseDate         string     `json:"expense_date"`
	Amount              string     `json:"amount"`
	GSTAmount           string     `json:"gst_amount"`
	PaymentMethod       *string    `json:"payment_method,omitempty"`
	Vendor              *string    `json:"vendor,omitempty"`
	Description         string     `json:"description"`
	ReceiptObjectKey    *string    `json:"receipt_object_key,omitempty"`
	Status              string     `json:"status"` // POSTED, REVERSED
	CreatedBy           uuid.UUID  `json:"created_by"`
	ReversalOfExpenseID *uuid.UUID `json:"reversal_of_expense_id,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
}

type CreateExpenseRequest struct {
	StoreID          *string `json:"store_id,omitempty"`
	CategoryID       string  `json:"category_id"`
	ExpenseDate      string  `json:"expense_date"`
	Amount           string  `json:"amount"`
	GSTAmount        string  `json:"gst_amount"`
	PaymentMethod    *string `json:"payment_method,omitempty"`
	Vendor           *string `json:"vendor,omitempty"`
	Description      string  `json:"description"`
	ReceiptObjectKey *string `json:"receipt_object_key,omitempty"`
}

type UpdateExpenseRequest struct {
	CategoryID       *string `json:"category_id,omitempty"`
	ExpenseDate      *string `json:"expense_date,omitempty"`
	Amount           *string `json:"amount,omitempty"`
	GSTAmount        *string `json:"gst_amount,omitempty"`
	PaymentMethod    *string `json:"payment_method,omitempty"`
	Vendor           *string `json:"vendor,omitempty"`
	Description      *string `json:"description,omitempty"`
	ReceiptObjectKey *string `json:"receipt_object_key,omitempty"`
}

type ReverseExpenseRequest struct {
	Reason string `json:"reason"`
}

type ListExpensesResponse struct {
	Expenses   []ExpenseModel `json:"expenses"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

type ListExpenseCategoriesResponse struct {
	Categories []ExpenseCategoryModel `json:"categories"`
	TotalCount int                    `json:"total_count"`
}
