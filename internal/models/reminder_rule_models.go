package models

import (
	"time"

	"github.com/google/uuid"
)

type ReminderRuleModel struct {
	ID             uuid.UUID  `json:"id"`
	SequenceNumber int        `json:"sequence_number"`
	DaysAfterDue   int        `json:"days_after_due"`
	TemplateID     uuid.UUID  `json:"template_id"`
	Status         string     `json:"status"` // ACTIVE, INACTIVE
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateReminderRuleRequest struct {
	SequenceNumber int    `json:"sequence_number"`
	DaysAfterDue   int    `json:"days_after_due"`
	TemplateID     string `json:"template_id"`
	Status         string `json:"status"` // ACTIVE by default
}

type UpdateReminderRuleRequest struct {
	SequenceNumber *int    `json:"sequence_number,omitempty"`
	DaysAfterDue   *int    `json:"days_after_due,omitempty"`
	TemplateID     *string `json:"template_id,omitempty"`
}

type ReorderRuleItem struct {
	ID             string `json:"id"`
	SequenceNumber int    `json:"sequence_number"`
}

type ReorderReminderRulesRequest struct {
	Rules []ReorderRuleItem `json:"rules"`
}

type ListReminderRulesResponse struct {
	Data []ReminderRuleModel `json:"data"`
}
