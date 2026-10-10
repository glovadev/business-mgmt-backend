package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentModel struct {
	ID                   uuid.UUID  `json:"id"`
	PaymentNumber        string     `json:"payment_number"`
	CustomerID           uuid.UUID  `json:"customer_id"`
	Amount               string     `json:"amount"`
	PaymentMethod        string     `json:"payment_method"`
	TransactionReference *string    `json:"transaction_reference,omitempty"`
	CollectionStoreID    *uuid.UUID `json:"collection_store_id,omitempty"`
	CollectedBy          uuid.UUID  `json:"collected_by"`
	CollectedAt          time.Time  `json:"collected_at"`
	Status               string     `json:"status"` // POSTED, REVERSED
	ReversalOfPaymentID  *uuid.UUID `json:"reversal_of_payment_id,omitempty"`
	Notes                *string    `json:"notes,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`

	Allocations []PaymentAllocationModel `json:"allocations,omitempty"`
}

// Ensure PaymentAllocationModel maps to schema
// Defined in sale_models.go previously, but let's redefine it here or use the existing one.
// We'll define a separate type specifically for the payments module if needed, but since it's 
// in models package, we can just use the existing one from sale_models.go and extend it.
// To avoid duplication, I will redefine it properly here and rely on Go's package scoping. 
// Wait, I already defined `PaymentAllocationModel` in sale_models.go. I'll use it.

type CreatePaymentRequest struct {
	CustomerID           string                          `json:"customer_id"`
	Amount               string                          `json:"amount"`
	PaymentMethod        string                          `json:"payment_method"`
	TransactionReference *string                         `json:"transaction_reference,omitempty"`
	CollectionStoreID    *string                         `json:"collection_store_id,omitempty"`
	AllocationStrategy   string                          `json:"allocation_strategy"` // OLDEST_FIRST, EXPLICIT
	Allocations          []CreatePaymentAllocationStruct `json:"allocations,omitempty"`
	Notes                *string                         `json:"notes,omitempty"`
}

type CreatePaymentAllocationStruct struct {
	InvoiceID string `json:"invoice_id"`
	Amount    string `json:"amount"`
}

type ReversePaymentRequest struct {
	Reason string `json:"reason"`
}

type ListPaymentsResponse struct {
	Payments   []PaymentModel `json:"payments"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}
