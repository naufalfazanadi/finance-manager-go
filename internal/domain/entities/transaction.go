package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TableName sets the table name
func (Transaction) TableName() string {
	return "transactions"
}

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string          `json:"name" gorm:"not null"`
	Cost      float64         `json:"cost" gorm:"type:decimal(20,8);not null"`
	Type      TransactionType `json:"type" gorm:"type:varchar(20);not null"`
	Note      string          `json:"note" gorm:"type:text"`
	TCategory string          `json:"t_category" gorm:"column:t_category;not null"`
	UserID    uuid.UUID       `json:"user_id" gorm:"type:uuid;not null;index"`
	WalletID  uuid.UUID       `json:"wallet_id" gorm:"type:uuid;not null;index"`
	IsDeleted bool            `json:"is_deleted" gorm:"column:is_deleted;default:false;index"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	// Belongs to User
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	// Belongs to Wallet
	Wallet Wallet `json:"wallet,omitempty" gorm:"foreignKey:WalletID"`
}

// IsSoftDeleted checks if transaction is soft deleted (either by boolean flag or DeletedAt timestamp)
func (t *Transaction) IsSoftDeleted() bool {
	return t.IsDeleted || t.DeletedAt.Valid
}

// IsActive checks if transaction is not soft deleted (neither boolean flag nor DeletedAt timestamp)
func (t *Transaction) IsActive() bool {
	return !t.IsDeleted && !t.DeletedAt.Valid
}

// SoftDelete marks the transaction as deleted (sets both boolean flag and DeletedAt timestamp)
func (t *Transaction) SoftDelete() {
	t.IsDeleted = true
	t.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
}

// Restore removes the soft delete flag (clears both boolean flag and DeletedAt)
func (t *Transaction) Restore() {
	t.IsDeleted = false
	t.DeletedAt = gorm.DeletedAt{
		Valid: false,
	}
}

// GetAbsoluteCost returns the absolute value of the transaction cost
func (t *Transaction) GetAbsoluteCost() float64 {
	if t.Cost < 0 {
		return -t.Cost
	}
	return t.Cost
}

// GetWalletImpact returns the amount that should be added to wallet balance
// For income: positive cost adds to balance
// For expense: positive cost subtracts from balance (returns negative value)
func (t *Transaction) GetWalletImpact() float64 {
	if t.Type == TransactionTypeIncome {
		return t.GetAbsoluteCost() // Always add positive amount for income
	} else {
		return -t.GetAbsoluteCost() // Always subtract positive amount for expense
	}
}
