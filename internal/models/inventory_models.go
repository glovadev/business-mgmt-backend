package models

import (
	"time"

	"github.com/google/uuid"
)

type InventoryModel struct {
	StoreID           uuid.UUID `json:"store_id"`
	ProductID         uuid.UUID `json:"product_id"`
	AvailableQuantity string    `json:"available_quantity"`
	AverageUnitCost   string    `json:"average_unit_cost"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Derived fields for API responses
	SKU            string `json:"sku,omitempty"`
	ProductName    string `json:"product_name,omitempty"`
	InventoryValue string `json:"inventory_value,omitempty"`
}

type InventoryMovementModel struct {
	ID             uuid.UUID  `json:"id"`
	StoreID        uuid.UUID  `json:"store_id"`
	ProductID      uuid.UUID  `json:"product_id"`
	MovementType   string     `json:"movement_type"`
	QuantityChange string     `json:"quantity_change"`
	UnitCost       *string    `json:"unit_cost,omitempty"`
	TotalCost      *string    `json:"total_cost,omitempty"`
	SourceType     string     `json:"source_type"`
	SourceID       *uuid.UUID `json:"source_id,omitempty"`
	SourceLineID   *uuid.UUID `json:"source_line_id,omitempty"`
	Reason         *string    `json:"reason,omitempty"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
}

type ListInventoryResponse struct {
	Data       []InventoryModel `json:"data"`
	TotalCount int              `json:"total_count"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
}

type ListInventoryMovementsResponse struct {
	Data       []InventoryMovementModel `json:"data"`
	TotalCount int                      `json:"total_count"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
}

type CreateInventoryAdjustmentRequest struct {
	StoreID        *string `json:"store_id,omitempty"`
	ProductID      string  `json:"product_id"`
	AdjustmentType string  `json:"adjustment_type"` // ADJUSTMENT_IN or ADJUSTMENT_OUT
	Quantity       string  `json:"quantity"`
	Reason         *string `json:"reason,omitempty"`
}
