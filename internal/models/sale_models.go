package models

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceModel struct {
	ID                uuid.UUID  `json:"id"`
	InvoiceNumber     string     `json:"invoice_number"`
	StoreID           uuid.UUID  `json:"store_id"`
	CustomerID        uuid.UUID  `json:"customer_id"`
	InvoiceDate       string     `json:"invoice_date"` // YYYY-MM-DD
	Subtotal          string     `json:"subtotal"`
	DiscountAmount    string     `json:"discount_amount"`
	GSTAmount         string     `json:"gst_amount"`
	TotalAmount       string     `json:"total_amount"`
	PaidAmount        string     `json:"paid_amount"`
	OutstandingAmount string     `json:"outstanding_amount"`
	PaymentStatus     string     `json:"payment_status"`
	DueDate           *string    `json:"due_date,omitempty"`
	Status            string     `json:"status"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	Items []InvoiceItemModel `json:"items,omitempty"`
}

type InvoiceItemModel struct {
	ID             uuid.UUID `json:"id"`
	InvoiceID      uuid.UUID `json:"invoice_id"`
	ProductID      uuid.UUID `json:"product_id"`
	Quantity       string    `json:"quantity"`
	UnitPrice      string    `json:"unit_price"`
	DiscountAmount string    `json:"discount_amount"`
	TaxableAmount  string    `json:"taxable_amount"`
	GSTRate        string    `json:"gst_rate"`
	GSTAmount      string    `json:"gst_amount"`
	LineTotal      string    `json:"line_total"`
	UnitCost       string    `json:"unit_cost"`
	CreatedAt      time.Time `json:"created_at"`
}

type InitialPaymentModel struct {
	Amount               string  `json:"amount"`
	PaymentMethod        string  `json:"payment_method"`
	TransactionReference *string `json:"transaction_reference,omitempty"`
}

type CreateSaleRequest struct {
	StoreID        string                     `json:"store_id"`
	CustomerID     string                     `json:"customer_id"` // Mandatory for V1
	InvoiceDate    string                     `json:"invoice_date"`
	Items          []CreateSaleItemRequest    `json:"items"`
	InitialPayment *InitialPaymentModel       `json:"initial_payment,omitempty"`
	DueDate        *string                    `json:"due_date,omitempty"`
	Notes          *string                    `json:"notes,omitempty"`
}

type CreateSaleItemRequest struct {
	ProductID      string `json:"product_id"`
	Quantity       string `json:"quantity"`
	UnitPrice      string `json:"unit_price"`
	DiscountAmount string `json:"discount_amount"`
}

type CancelSaleRequest struct {
	Reason string `json:"reason"`
}

type ListSalesResponse struct {
	Sales      []InvoiceModel `json:"sales"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

// Payment Allocation stub
type PaymentAllocationModel struct {
	ID        uuid.UUID `json:"id"`
	PaymentID uuid.UUID `json:"payment_id"`
	InvoiceID uuid.UUID `json:"invoice_id"`
	Amount    string    `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
