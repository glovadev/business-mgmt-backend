package models

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	Line1       string  `json:"line1"`
	Line2       *string `json:"line2,omitempty"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	PostalCode  string  `json:"postal_code"`
	CountryCode string  `json:"country_code"`
}

type StoreModel struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone,omitempty"`
	GSTIN     *string   `json:"gstin,omitempty"`
	Address   Address   `json:"address"`
	Status    string    `json:"status"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateStoreRequest struct {
	Code    string  `json:"code"`
	Name    string  `json:"name"`
	Phone   *string `json:"phone,omitempty"`
	GSTIN   *string `json:"gstin,omitempty"`
	Address Address `json:"address"`
}

type UpdateStoreRequest struct {
	Name    *string  `json:"name,omitempty"`
	Phone   *string  `json:"phone,omitempty"`
	GSTIN   *string  `json:"gstin,omitempty"`
	Address *Address `json:"address,omitempty"`
}

type ListStoresResponse struct {
	Stores     []StoreModel `json:"stores"`
	TotalCount int          `json:"total_count"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
}
