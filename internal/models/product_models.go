package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductModel struct {
	ID                  uuid.UUID  `json:"id"`
	SKU                 string     `json:"sku"`
	Name                string     `json:"name"`
	CategoryID          uuid.UUID  `json:"category_id"`
	Description         *string    `json:"description,omitempty"`
	HSNCode             *string    `json:"hsn_code,omitempty"`
	UnitOfMeasure       string     `json:"unit_of_measure"`
	DefaultSellingPrice string     `json:"default_selling_price"` // String for NUMERIC handling
	GSTRate             string     `json:"gst_rate"`              // String for NUMERIC handling
	Status              string     `json:"status"`
	CreatedBy           *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type CreateProductRequest struct {
	SKU                 string  `json:"sku"`
	Name                string  `json:"name"`
	CategoryID          string  `json:"category_id"`
	Description         *string `json:"description,omitempty"`
	HSNCode             *string `json:"hsn_code,omitempty"`
	UnitOfMeasure       string  `json:"unit_of_measure"`
	DefaultSellingPrice string  `json:"default_selling_price"`
	GSTRate             string  `json:"gst_rate"`
}

type UpdateProductRequest struct {
	SKU                 *string `json:"sku,omitempty"`
	Name                *string `json:"name,omitempty"`
	CategoryID          *string `json:"category_id,omitempty"`
	Description         *string `json:"description,omitempty"`
	HSNCode             *string `json:"hsn_code,omitempty"`
	UnitOfMeasure       *string `json:"unit_of_measure,omitempty"`
	DefaultSellingPrice *string `json:"default_selling_price,omitempty"`
	GSTRate             *string `json:"gst_rate,omitempty"`
}

type ListProductsResponse struct {
	Products   []ProductModel `json:"products"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

type ProductInventoryData struct {
	StoreID           uuid.UUID `json:"store_id"`
	StoreName         string    `json:"store_name"`
	AvailableQuantity string    `json:"available_quantity"`
	AverageUnitCost   string    `json:"average_unit_cost"`
}

type ProductInventoryResponse struct {
	Data []ProductInventoryData `json:"data"`
}
