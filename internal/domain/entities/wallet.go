package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TableName sets the table name
func (Wallet) TableName() string {
	return "wallets"
}

type Wallet struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null"`
	Type      string         `json:"type" gorm:"not null"`
	Category  string         `json:"category" gorm:"not null"`
	Balance   float64        `json:"balance" gorm:"type:decimal(20,8);default:0"`
	Currency  string         `json:"currency" gorm:"not null;default:'IDR'"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	IsDeleted bool           `json:"is_deleted" gorm:"column:is_deleted;default:false;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	// Belongs to User
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// IsSoftDeleted checks if wallet is soft deleted (either by boolean flag or DeletedAt timestamp)
func (u *Wallet) IsSoftDeleted() bool {
	return u.IsDeleted || u.DeletedAt.Valid
}

// IsActive checks if wallet is not soft deleted (neither boolean flag nor DeletedAt timestamp)
func (u *Wallet) IsActive() bool {
	return !u.IsDeleted && !u.DeletedAt.Valid
}

// SoftDelete marks the wallet as deleted (sets both boolean flag and DeletedAt timestamp)
func (u *Wallet) SoftDelete() {
	u.IsDeleted = true
	u.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
}

// Restore removes the soft delete flag (clears both boolean flag and DeletedAt)
func (u *Wallet) Restore() {
	u.IsDeleted = false
	u.DeletedAt = gorm.DeletedAt{
		Valid: false,
	}
}
