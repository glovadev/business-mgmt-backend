package reminderrule

import (
	"context"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateRule(ctx context.Context, rule *models.ReminderRuleModel) error
	ListRules(ctx context.Context) ([]models.ReminderRuleModel, error)
	GetRuleByID(ctx context.Context, id uuid.UUID) (*models.ReminderRuleModel, error)
	UpdateRule(ctx context.Context, rule *models.ReminderRuleModel) error
	
	ExecuteReorderTransaction(ctx context.Context, updates map[uuid.UUID]int) error
}
