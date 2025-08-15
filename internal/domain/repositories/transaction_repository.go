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

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entities.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Transaction, error)
	GetOne(ctx context.Context, filter map[string]interface{}) (*entities.Transaction, error)
	GetAll(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.Transaction, error)
	Update(ctx context.Context, transaction *entities.Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
	CountWithFilters(ctx context.Context, queryParams *dto.QueryParams) (int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	GetByWalletID(ctx context.Context, walletID uuid.UUID) ([]*entities.Transaction, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *entities.Transaction) error {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Transaction, error) {
	var transaction entities.Transaction
	if err := r.db.Preload("User").Preload("Wallet").WithContext(ctx).First(&transaction, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetOne(ctx context.Context, filter map[string]interface{}) (*entities.Transaction, error) {
	var transaction entities.Transaction
	query := r.db.WithContext(ctx)

	for key, value := range filter {
		query = query.Where("LOWER("+key+") = LOWER(?)", value)
	}

	if err := query.Preload("User").Preload("Wallet").First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetAll(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.Transaction, error) {
	var transactions []*entities.Transaction
	query := r.db.WithContext(ctx)

	if queryParams.LoggedUserID != uuid.Nil {
		query = query.Where("user_id = ?", queryParams.LoggedUserID)
	}

	// Apply search if provided
	if queryParams.HasSearch() {
		searchTerm := "%" + queryParams.Search + "%"
		query = query.Where("name ILIKE ? OR note ILIKE ? OR t_category ILIKE ? OR type ILIKE ?", searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// Apply custom filters
	if queryParams.HasFilters() {
		for key, value := range queryParams.Filters {
			// Only allow safe column names to prevent SQL injection
			switch key {
			case "t_category", "name", "type":
				query = query.Where("LOWER("+key+") = LOWER(?)", value)
			case "wallet_id", "user_id":
				query = query.Where(key+" = ?", value)
			case "cost_min":
				query = query.Where("cost >= ?", value)
			case "cost_max":
				query = query.Where("cost <= ?", value)
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
			"cost":       true,
			"t_category": true,
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

	if err := query.Preload("User").Preload("Wallet").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) Update(ctx context.Context, transaction *entities.Transaction) error {
	if err := r.db.WithContext(ctx).Save(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&entities.Transaction{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.Transaction{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *transactionRepository) CountWithFilters(ctx context.Context, queryParams *dto.QueryParams) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Transaction{})

	if queryParams.LoggedUserID != uuid.Nil {
		query = query.Where("user_id = ?", queryParams.LoggedUserID)
	}

	// Apply search if provided
	if queryParams.HasSearch() {
		searchTerm := "%" + queryParams.Search + "%"
		query = query.Where("name ILIKE ? OR note ILIKE ? OR t_category ILIKE ? OR type ILIKE ?", searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// Apply custom filters
	if queryParams.HasFilters() {
		for key, value := range queryParams.Filters {
			// Only allow safe column names to prevent SQL injection
			switch key {
			case "t_category", "name", "type":
				query = query.Where("LOWER("+key+") = LOWER(?)", value)
			case "wallet_id", "user_id":
				query = query.Where(key+" = ?", value)
			case "cost_min":
				query = query.Where("cost >= ?", value)
			case "cost_max":
				query = query.Where("cost <= ?", value)
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

// SoftDelete soft deletes a transaction by ID
func (r *transactionRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&entities.Transaction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("transaction not found")
	}

	return nil
}

// HardDelete permanently deletes a transaction from the database
func (r *transactionRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Unscoped().Delete(&entities.Transaction{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// Restore restores a soft deleted transaction by ID
func (r *transactionRepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Model(&entities.Transaction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": false,
			"deleted_at": nil,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("transaction not found")
	}

	return nil
}

// GetByWalletID gets all transactions for a specific wallet
func (r *transactionRepository) GetByWalletID(ctx context.Context, walletID uuid.UUID) ([]*entities.Transaction, error) {
	var transactions []*entities.Transaction
	if err := r.db.WithContext(ctx).Where("wallet_id = ? AND is_deleted = false", walletID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetByUserID gets all transactions for a specific user
func (r *transactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Transaction, error) {
	var transactions []*entities.Transaction
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_deleted = false", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
