package models

import (
	"time"

	"github.com/google/uuid"
)

type CustomerModel struct {
	ID                  uuid.UUID  `json:"id"`
	Name                string     `json:"name"`
	PhoneNumber         string     `json:"phone_number"`
	AlternatePhone      *string    `json:"alternate_phone,omitempty"`
	Email               *string    `json:"email,omitempty"`
	Address             *string    `json:"address,omitempty"`
	GSTIN               *string    `json:"gstin,omitempty"`
	WhatsappOptIn       bool       `json:"whatsapp_opt_in"`
	WhatsappOptedOutAt  *time.Time `json:"whatsapp_opted_out_at,omitempty"`
	Status              string     `json:"status"`
	Notes               *string    `json:"notes,omitempty"`
	CreatedBy           *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`

	// Derived fields as per API spec
	RiskStatus       *string  `json:"risk_status,omitempty"`
	TotalOutstanding *string  `json:"total_outstanding,omitempty"`
	OldestDueDate    *string  `json:"oldest_due_date,omitempty"`
	MaxOverdueDays   *int     `json:"max_overdue_days,omitempty"`
}

type CreateCustomerRequest struct {
	Name           string  `json:"name"`
	PhoneNumber    string  `json:"phone_number"`
	AlternatePhone *string `json:"alternate_phone,omitempty"`
	Email          *string `json:"email,omitempty"`
	Address        *string `json:"address,omitempty"`
	GSTIN          *string `json:"gstin,omitempty"`
	WhatsappOptIn  bool    `json:"whatsapp_opt_in"`
	StoreID        string  `json:"store_id"` // required as per docs
	Notes          *string `json:"notes,omitempty"`
}

type UpdateCustomerRequest struct {
	Name           *string `json:"name,omitempty"`
	PhoneNumber    *string `json:"phone_number,omitempty"`
	AlternatePhone *string `json:"alternate_phone,omitempty"`
	Email          *string `json:"email,omitempty"`
	Address        *string `json:"address,omitempty"`
	GSTIN          *string `json:"gstin,omitempty"`
	WhatsappOptIn  *bool   `json:"whatsapp_opt_in,omitempty"`
	Notes          *string `json:"notes,omitempty"`
}

type ListCustomersResponse struct {
	Customers  []CustomerModel `json:"customers"`
	TotalCount int             `json:"total_count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
}

type CustomerStoreRelationship struct {
	CustomerID   uuid.UUID  `json:"customer_id"`
	StoreID      uuid.UUID  `json:"store_id"`
	FirstAddedBy *uuid.UUID `json:"first_added_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
