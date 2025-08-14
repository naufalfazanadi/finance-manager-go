package entities

import (
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/pkg/config"
	"github.com/naufalfazanadi/finance-manager-go/pkg/encryption"
	"github.com/naufalfazanadi/finance-manager-go/pkg/minio"
	"gorm.io/gorm"
)

// UserRole represents user roles in the system
type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// TableName sets the table name
func (User) TableName() string {
	return "users"
}

type User struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email              string         `json:"email" gorm:"-"`                         // Decrypted email (not stored)
	EmailHash          string         `json:"-" gorm:"column:email_hash;uniqueIndex"` // Hash for indexing
	EmailEncrypted     string         `json:"-" gorm:"column:email_encrypted"`        // Encrypted email storage
	BirthDate          *time.Time     `json:"birth_date" gorm:"-"`                    // Decrypted birth date (not stored)
	BirthDateEncrypted string         `json:"-" gorm:"column:birth_date_encrypted"`   // Encrypted birth date storage
	Name               string         `json:"name" gorm:"not null"`
	Password           string         `json:"-" gorm:"not null"`
	Role               UserRole       `json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	ProfilePhoto       string         `json:"profile_photo" gorm:"column:profile_photo"`
	IsDeleted          bool           `json:"is_deleted" gorm:"column:is_deleted;default:false;index"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	// One-to-Many: User can have multiple wallets (if you need multiple wallets per user)
	Wallets []Wallet `json:"wallets,omitempty" gorm:"foreignKey:UserID"`
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// IsUser checks if user has user role
func (u *User) IsUser() bool {
	return u.Role == UserRoleUser
}

// IsSoftDeleted checks if user is soft deleted (either by boolean flag or DeletedAt timestamp)
func (u *User) IsSoftDeleted() bool {
	return u.IsDeleted || u.DeletedAt.Valid
}

// IsActive checks if user is not soft deleted (neither boolean flag nor DeletedAt timestamp)
func (u *User) IsActive() bool {
	return !u.IsDeleted && !u.DeletedAt.Valid
}

// GetAge calculates the user's age based on birth date
func (u *User) GetAge() *int {
	if u.BirthDate == nil {
		return nil
	}

	now := time.Now()
	age := now.Year() - u.BirthDate.Year()

	// Adjust if birthday hasn't occurred this year
	if now.Month() < u.BirthDate.Month() ||
		(now.Month() == u.BirthDate.Month() && now.Day() < u.BirthDate.Day()) {
		age--
	}

	return &age
}

// GetProfilePhotoURL returns the full URL for the profile photo
func (u *User) GetProfilePhotoURL() string {
	if u.ProfilePhoto == "" {
		return ""
	}

	// Get config and use public bucket for profile photos
	cfg := config.GetConfig()
	return minio.GetFullUrl(cfg.Minio.PrivateBucket, u.ProfilePhoto)
}

// SoftDelete marks the user as deleted (sets both boolean flag and DeletedAt timestamp)
func (u *User) SoftDelete() {
	u.IsDeleted = true
	u.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
}

// Restore removes the soft delete flag (clears both boolean flag and DeletedAt)
func (u *User) Restore() {
	u.IsDeleted = false
	u.DeletedAt = gorm.DeletedAt{
		Valid: false,
	}
}

// SetPII encrypts and sets the email address and birth date using the current User entity values
func (u *User) SetPII() error {
	// Handle email encryption
	if u.Email == "" {
		u.EmailHash = ""
		u.EmailEncrypted = ""
	} else {
		// Hash the email for indexing
		hashResult := encryption.HashSHA256(u.Email)
		if hashResult.Error != nil {
			return hashResult.Error
		}

		// Encrypt the email for storage
		encResult := encryption.EncryptAES128GCM(u.Email)
		if encResult.Error != nil {
			return encResult.Error
		}

		u.EmailHash = hashResult.Data.(string)

		// Convert encrypted bytes to base64 string for storage
		encryptedBytes := encResult.Data.([]byte)
		u.EmailEncrypted = base64.StdEncoding.EncodeToString(encryptedBytes)
	}

	// Handle birth date encryption
	if u.BirthDate == nil {
		u.BirthDateEncrypted = ""
	} else {
		// Format birth date as ISO 8601 string for encryption
		birthDateStr := u.BirthDate.Format(time.RFC3339)

		// Encrypt the birth date for storage (no hashing needed)
		encResult := encryption.EncryptAES128GCM(birthDateStr)
		if encResult.Error != nil {
			return encResult.Error
		}

		// Convert encrypted bytes to base64 string for storage
		encryptedBytes := encResult.Data.([]byte)
		u.BirthDateEncrypted = base64.StdEncoding.EncodeToString(encryptedBytes)
	}

	return nil
}

// LoadPII decrypts and loads the email address and birth date
func (u *User) LoadPII() error {
	// Handle email decryption
	if u.EmailEncrypted == "" {
		u.Email = ""
	} else {
		// Pass the base64 encoded string directly to DecryptAES128GCM
		decResult := encryption.DecryptAES128GCM(u.EmailEncrypted)
		if decResult.Error != nil {
			return decResult.Error
		}
		u.Email = decResult.Data.(string)
	}

	// Handle birth date decryption
	if u.BirthDateEncrypted == "" {
		u.BirthDate = nil
	} else {
		// Decrypt the birth date
		decResult := encryption.DecryptAES128GCM(u.BirthDateEncrypted)
		if decResult.Error != nil {
			return decResult.Error
		}

		// Parse the decrypted string back to time.Time
		birthDateStr := decResult.Data.(string)
		parsedTime, err := time.Parse(time.RFC3339, birthDateStr)
		if err != nil {
			return err
		}
		u.BirthDate = &parsedTime
	}

	return nil
}

// BeforeCreate hook to generate UUID, set default role, and encrypt email
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.Role == "" {
		u.Role = UserRoleUser
	}

	// Encrypt email and birth date if they're set and not already encrypted
	if (u.Email != "" && u.EmailEncrypted == "") || (u.BirthDate != nil && u.BirthDateEncrypted == "") {
		if err := u.SetPII(); err != nil {
			return err
		}
	}

	return nil
}

// BeforeUpdate hook - runs before update operations
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// If email or birth date is being updated, encrypt them
	if u.Email != "" || u.BirthDate != nil {
		// Check if we need to encrypt (email or birth date changed)
		if err := u.SetPII(); err != nil {
			return err
		}
	}
	return nil
}

// BeforeSave hook - runs before both create and update
func (u *User) BeforeSave(tx *gorm.DB) error {
	// Encrypt email and birth date if they're set and not already encrypted
	if (u.Email != "" && u.EmailEncrypted == "") || (u.BirthDate != nil && u.BirthDateEncrypted == "") {
		if err := u.SetPII(); err != nil {
			return err
		}
	}
	return nil
}

// AfterFind hook - runs after record is found/retrieved
func (u *User) AfterFind(tx *gorm.DB) error {
	// Decrypt email after loading from database
	if err := u.LoadPII(); err != nil {
		// Log error but don't fail the operation
		// You might want to add proper logging here
		return nil
	}
	return nil
}

// AfterCreate hook - runs after record is created
func (u *User) AfterCreate(tx *gorm.DB) error {
	// Add any post-create logic here
	// For example: send welcome email, create related records, etc.
	return nil
}

// AfterUpdate hook - runs after record is updated
func (u *User) AfterUpdate(tx *gorm.DB) error {
	// Add any post-update logic here
	// For example: invalidate cache, send notifications, etc.
	return nil
}

// AfterSave hook - runs after both create and update
func (u *User) AfterSave(tx *gorm.DB) error {
	// Add any post-save logic here
	// For example: audit logging, cache updates, etc.
	return nil
}
