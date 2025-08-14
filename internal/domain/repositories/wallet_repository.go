package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *entities.Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Wallet, error)
	GetOne(ctx context.Context, filter map[string]interface{}) (*entities.Wallet, error)
	GetAll(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.Wallet, error)
	Update(ctx context.Context, wallet *entities.Wallet) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
	CountWithFilters(ctx context.Context, queryParams *dto.QueryParams) (int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(ctx context.Context, wallet *entities.Wallet) error {
	if err := r.db.WithContext(ctx).Create(wallet).Error; err != nil {
		return err
	}
	return nil
}

func (r *walletRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Wallet, error) {
	var wallet entities.Wallet
	if err := r.db.Preload("User").WithContext(ctx).First(&wallet, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) GetOne(ctx context.Context, filter map[string]interface{}) (*entities.Wallet, error) {
	var wallet entities.Wallet
	query := r.db.WithContext(ctx)

	for key, value := range filter {
		query = query.Where("LOWER("+key+") = LOWER(?)", value)
	}

	if err := query.Preload("User").First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) GetAll(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.Wallet, error) {
	var wallets []*entities.Wallet
	query := r.db.WithContext(ctx)

	if queryParams.LoggedUserID != uuid.Nil {
		query = query.Where("user_id = ?", queryParams.LoggedUserID)
	}

	// Apply search if provided
	if queryParams.HasSearch() {
		searchTerm := "%" + queryParams.Search + "%"
		query = query.Where("name ILIKE ? OR type ILIKE ? OR category ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Apply custom filters
	if queryParams.HasFilters() {
		for key, value := range queryParams.Filters {
			// Only allow safe column names to prevent SQL injection
			switch key {
			case "category", "name", "type":
				query = query.Where("LOWER("+key+") = LOWER(?)", value)
			case "created_after":
				query = query.Where("created_at >= ?", value)
			case "created_before":
				query = query.Where("created_at <= ?", value)
			}
		}
	}

	// Apply sorting
	if queryParams.HasSort() {
		// Only allow safe column names for sorting
		allowedSortColumns := map[string]bool{
			"name":       true,
			"type":       true,
			"category":   true,
			"created_at": true,
			"updated_at": true,
		}

		if allowedSortColumns[queryParams.SortBy] {
			orderClause := queryParams.SortBy + " " + queryParams.SortType
			query = query.Order(orderClause)
		}
	} else {
		// Default sorting
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if queryParams.Limit > 0 {
		query = query.Limit(queryParams.Limit)
	}
	if queryParams.GetOffset() > 0 {
		query = query.Offset(queryParams.GetOffset())
	}

	if err := query.Preload("User").Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *walletRepository) Update(ctx context.Context, wallet *entities.Wallet) error {
	if err := r.db.WithContext(ctx).Save(wallet).Error; err != nil {
		return err
	}
	return nil
}

func (r *walletRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&entities.Wallet{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *walletRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.Wallet{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *walletRepository) CountWithFilters(ctx context.Context, queryParams *dto.QueryParams) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Wallet{})

	if queryParams.LoggedUserID != uuid.Nil {
		query = query.Where("user_id = ?", queryParams.LoggedUserID)
	}

	// Apply search if provided
	if queryParams.HasSearch() {
		searchTerm := "%" + queryParams.Search + "%"
		query = query.Where("name ILIKE ? OR type ILIKE ? OR category ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Apply custom filters
	if queryParams.HasFilters() {
		for key, value := range queryParams.Filters {
			// Only allow safe column names to prevent SQL injection
			switch key {
			case "type", "name", "category":
				query = query.Where(key+" = ?", value)
			case "created_after":
				query = query.Where("created_at >= ?", value)
			case "created_before":
				query = query.Where("created_at <= ?", value)
			}
		}
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// SoftDelete soft deletes a wallet by ID
func (r *walletRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	// Alternative approach 1: Update both fields manually
	result := r.db.WithContext(ctx).Model(&entities.Wallet{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("wallet not found")
	}

	return nil

	// Alternative approach 2: Fetch, call entity method, then save
	/*
		var wallet entities.Wallet
		if err := r.db.WithContext(ctx).First(&wallet, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("wallet not found")
			}
			return err
		}

		wallet.SoftDelete()
		return r.db.WithContext(ctx).Save(&wallet).Error
	*/
}

// HardDelete permanently deletes a wallet from the database
func (r *walletRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Unscoped().Delete(&entities.Wallet{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// Restore restores a soft deleted wallet by ID
func (r *walletRepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Model(&entities.Wallet{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": false,
			"deleted_at": nil,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("wallet not found")
	}

	return nil
}
