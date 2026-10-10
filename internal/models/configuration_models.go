package models

import (
	"time"

	"github.com/google/uuid"
)

type CompanySettingsModel struct {
	ID                              uuid.UUID  `json:"-"`
	CompanyName                     string     `json:"company_name"`
	Currency                        string     `json:"currency"`
	Timezone                        string     `json:"timezone"`
	DefaultCreditDays               int        `json:"default_credit_days"`
	YellowAfterOverdueDays          int        `json:"yellow_after_overdue_days"`
	RedAfterOverdueDays             int        `json:"red_after_overdue_days"`
	AdditionalCostAllocationMethod  string     `json:"additional_cost_allocation_method"`
	IncludeGSTInLandedCost          bool       `json:"include_gst_in_landed_cost"`
	CreatedAt                       time.Time  `json:"-"`
	UpdatedAt                       time.Time  `json:"-"`
	UpdatedBy                       *uuid.UUID `json:"-"`
}

type GetCompanySettingsResponse struct {
	Data CompanySettingsModel `json:"data"`
}

type UpdateCompanySettingsRequest struct {
	DefaultCreditDays              *int    `json:"default_credit_days,omitempty"`
	YellowAfterOverdueDays         *int    `json:"yellow_after_overdue_days,omitempty"`
	RedAfterOverdueDays            *int    `json:"red_after_overdue_days,omitempty"`
	CompanyName                    *string `json:"company_name,omitempty"`
	Currency                       *string `json:"currency,omitempty"`
	Timezone                       *string `json:"timezone,omitempty"`
	AdditionalCostAllocationMethod *string `json:"additional_cost_allocation_method,omitempty"`
	IncludeGSTInLandedCost         *bool   `json:"include_gst_in_landed_cost,omitempty"`
}
