package models

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseModel struct {
	ID                    uuid.UUID  `json:"id"`
	PurchaseNumber        string     `json:"purchase_number"`
	StoreID               uuid.UUID  `json:"store_id"`
	SupplierID            *uuid.UUID `json:"supplier_id,omitempty"`
	SupplierInvoiceNumber *string    `json:"supplier_invoice_number,omitempty"`
	PurchaseDate          string     `json:"purchase_date"` // YYYY-MM-DD
	Status                string     `json:"status"`
	Subtotal              string     `json:"subtotal"`
	GSTAmount             string     `json:"gst_amount"`
	AdditionalCostTotal   string     `json:"additional_cost_total"`
	GrandTotal            string     `json:"grand_total"`
	Notes                 *string    `json:"notes,omitempty"`
	CreatedBy             *uuid.UUID `json:"created_by,omitempty"`
	SubmittedAt           *time.Time `json:"submitted_at,omitempty"`
	ApprovedBy            *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt            *time.Time `json:"approved_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`

	Items           []PurchaseItemModel           `json:"items,omitempty"`
	AdditionalCosts []PurchaseAdditionalCostModel `json:"additional_costs,omitempty"`
}

type PurchaseItemModel struct {
	ID                 uuid.UUID `json:"id"`
	PurchaseID         uuid.UUID `json:"purchase_id"`
	ProductID          uuid.UUID `json:"product_id"`
	Quantity           string    `json:"quantity"`
	BaseUnitCost       string    `json:"base_unit_cost"`
	BaseAmount         string    `json:"base_amount"`
	GSTRate            string    `json:"gst_rate"`
	GSTAmount          string    `json:"gst_amount"`
	AllocatedExtraCost string    `json:"allocated_extra_cost"`
	LandedUnitCost     string    `json:"landed_unit_cost"`
	TotalAmount        string    `json:"total_amount"`
}

type PurchaseAdditionalCostModel struct {
	ID          uuid.UUID `json:"id"`
	PurchaseID  uuid.UUID `json:"purchase_id"`
	CostType    string    `json:"cost_type"`
	Description *string   `json:"description,omitempty"`
	Amount      string    `json:"amount"`
}

type CreatePurchaseRequest struct {
	StoreID               string                                `json:"store_id"`
	SupplierID            *string                               `json:"supplier_id,omitempty"`
	SupplierInvoiceNumber *string                               `json:"supplier_invoice_number,omitempty"`
	PurchaseDate          string                                `json:"purchase_date"`
	Items                 []CreatePurchaseItemRequest           `json:"items"`
	AdditionalCosts       []CreatePurchaseAdditionalCostRequest `json:"additional_costs,omitempty"`
	Notes                 *string                               `json:"notes,omitempty"`
}

type CreatePurchaseItemRequest struct {
	ProductID    string `json:"product_id"`
	Quantity     string `json:"quantity"`
	BaseUnitCost string `json:"base_unit_cost"`
	GSTRate      string `json:"gst_rate"`
}

type CreatePurchaseAdditionalCostRequest struct {
	CostType    string  `json:"cost_type"`
	Description *string `json:"description,omitempty"`
	Amount      string  `json:"amount"`
}

type RejectPurchaseRequest struct {
	Reason string `json:"reason"`
}

type ListPurchasesResponse struct {
	Purchases  []PurchaseModel `json:"purchases"`
	TotalCount int             `json:"total_count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
}
