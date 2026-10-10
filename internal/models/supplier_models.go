package models

import (
	"time"

	"github.com/google/uuid"
)

type SupplierModel struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Phone     *string    `json:"phone,omitempty"`
	Email     *string    `json:"email,omitempty"`
	GSTIN     *string    `json:"gstin,omitempty"`
	Address   *string    `json:"address,omitempty"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateSupplierRequest struct {
	Name    string  `json:"name"`
	Phone   *string `json:"phone,omitempty"`
	Email   *string `json:"email,omitempty"`
	GSTIN   *string `json:"gstin,omitempty"`
	Address *string `json:"address,omitempty"`
}

type UpdateSupplierRequest struct {
	Name    *string `json:"name,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Email   *string `json:"email,omitempty"`
	GSTIN   *string `json:"gstin,omitempty"`
	Address *string `json:"address,omitempty"`
}

type ListSuppliersResponse struct {
	Suppliers  []SupplierModel `json:"suppliers"`
	TotalCount int             `json:"total_count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
}
