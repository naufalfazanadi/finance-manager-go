package repositories

import (
	"context"
	"errors"

	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByEmailHash(ctx context.Context, emailHash string) (*entities.User, error)
	GetOne(ctx context.Context, filter map[string]interface{}) (*entities.User, error)
	GetAll(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
	CountWithFilters(ctx context.Context, queryParams *dto.QueryParams) (int64, error)
	// Soft delete methods
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	GetWithDeleted(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.User, error)
	GetOnlyDeleted(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.User, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetByEmailHash(ctx context.Context, emailHash string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).First(&user, "email_hash = ?", emailHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetOne(ctx context.Context, filter map[string]interface{}) (*entities.User, error) {
	var user entities.User
	query := r.db.WithContext(ctx)

	for key, value := range filter {
		query = query.Where("LOWER("+key+") = LOWER(?)", value)
	}

	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetAll(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.User, error) {
	var users []*entities.User
	query := r.db.WithContext(ctx)

	// Apply search if provided
	if queryParams.HasSearch() {
		searchTerm := "%" + queryParams.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	// Apply custom filters
	if queryParams.HasFilters() {
		for key, value := range queryParams.Filters {
			// Only allow safe column names to prevent SQL injection
			switch key {
			case "role", "name", "email":
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
			"email":      true,
			"role":       true,
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

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepositoryImpl) CountWithFilters(ctx context.Context, queryParams *dto.QueryParams) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.User{})

	// Apply search if provided
	if queryParams.HasSearch() {
		searchTerm := "%" + queryParams.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	// Apply custom filters
	if queryParams.HasFilters() {
		for key, value := range queryParams.Filters {
			// Only allow safe column names to prevent SQL injection
			switch key {
			case "role", "name", "email":
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

// SoftDelete soft deletes a user by ID
func (r *userRepositoryImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// Restore restores a soft deleted user by ID
func (r *userRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Unscoped().Model(&entities.User{}).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		return err
	}
	return nil
}

// GetWithDeleted gets all users including soft deleted ones
func (r *userRepositoryImpl) GetWithDeleted(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.User, error) {
	var users []*entities.User
	query := r.db.WithContext(ctx).Unscoped()

	// Apply search if provided
	if queryParams.HasSearch() {
		searchTerm := "%" + queryParams.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	// Apply custom filters
	if queryParams.HasFilters() {
		for key, value := range queryParams.Filters {
			// Only allow safe column names to prevent SQL injection
			switch key {
			case "role", "name", "email":
				query = query.Where(key+" = ?", value)
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
			"email":      true,
			"role":       true,
			"created_at": true,
			"updated_at": true,
			"deleted_at": true,
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

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetOnlyDeleted gets only soft deleted users
func (r *userRepositoryImpl) GetOnlyDeleted(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.User, error) {
	var users []*entities.User
	query := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL")

	// Apply search if provided
	if queryParams.HasSearch() {
		searchTerm := "%" + queryParams.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	// Apply custom filters
	if queryParams.HasFilters() {
		for key, value := range queryParams.Filters {
			// Only allow safe column names to prevent SQL injection
			switch key {
			case "role", "name", "email":
				query = query.Where(key+" = ?", value)
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
			"email":      true,
			"role":       true,
			"created_at": true,
			"updated_at": true,
			"deleted_at": true,
		}

		if allowedSortColumns[queryParams.SortBy] {
			orderClause := queryParams.SortBy + " " + queryParams.SortType
			query = query.Order(orderClause)
		}
	} else {
		// Default sorting
		query = query.Order("deleted_at DESC")
	}

	// Apply pagination
	if queryParams.Limit > 0 {
		query = query.Limit(queryParams.Limit)
	}
	if queryParams.GetOffset() > 0 {
		query = query.Offset(queryParams.GetOffset())
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// HardDelete permanently deletes a user from the database
func (r *userRepositoryImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Unscoped().Delete(&entities.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
