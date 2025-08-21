package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetMonthlySumByUser(ctx context.Context, userID uuid.UUID) ([]*entities.VMonthlyTransactionSum, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetMonthlySumByUser(ctx context.Context, userID uuid.UUID) ([]*entities.VMonthlyTransactionSum, error) {
	var summaries []*entities.VMonthlyTransactionSum
	if err := r.db.WithContext(ctx).
		Table("v_monthly_transaction_sum").
		Where("user_id = ?", userID).
		Order("month DESC").
		Find(&summaries).Error; err != nil {
		return nil, fmt.Errorf("query monthly summary: %w", err)
	}
	return summaries, nil
}
