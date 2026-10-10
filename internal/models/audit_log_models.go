package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLogModel struct {
	ID          uuid.UUID  `json:"id"`
	ActorUserID uuid.UUID  `json:"actor_user_id"`
	StoreID     *uuid.UUID `json:"store_id,omitempty"`
	Action      string     `json:"action"`
	EntityType  string     `json:"entity_type"`
	EntityID    uuid.UUID  `json:"entity_id"`
	OldValue    *string    `json:"old_value,omitempty"` // JSONB as string
	NewValue    *string    `json:"new_value,omitempty"` // JSONB as string
	RequestID   string     `json:"request_id"`
	IPAddress   *string    `json:"ip_address,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type ListAuditLogsResponse struct {
	AuditLogs  []AuditLogModel `json:"audit_logs"`
	TotalCount int             `json:"total_count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
}
