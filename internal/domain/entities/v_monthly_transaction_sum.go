package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VMonthlyTransactionSum represents the view, not a table (ignored by GORM migration)
type VMonthlyTransactionSum struct {
	UserID           uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id;index"`
	WalletID         uuid.UUID `json:"wallet_id" gorm:"type:uuid;column:wallet_id"`
	Month            string    `json:"month" gorm:"column:month"`
	TransactionCount int64     `json:"transaction_count" gorm:"column:transaction_count"`
	TotalCost        float64   `json:"total_cost" gorm:"column:total_cost"`
}

// MigrateVMonthlyTransactionSumView creates the view in the database
func MigrateVMonthlyTransactionSumView(db *gorm.DB) error {
	return db.Exec(`
CREATE OR REPLACE VIEW v_monthly_transaction_sum AS
SELECT
  user_id,
  wallet_id,
  TO_CHAR(DATE_TRUNC('month', created_at), 'YYYY-MM') AS month,
  COUNT(*) AS transaction_count,
  SUM(cost) AS total_cost
FROM
  transactions
GROUP BY
  user_id,
  wallet_id,
  DATE_TRUNC('month', created_at);
`).Error
}
