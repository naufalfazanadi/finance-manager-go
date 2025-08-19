package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
)

const (
	UserCacheKeyPrefix  = "user:"
	DefaultUserCacheTTL = 24 * time.Hour
)

type CachedUser struct {
	ID        uuid.UUID         `json:"id"`
	Email     string            `json:"email"`
	Name      string            `json:"name"`
	Role      entities.UserRole `json:"role"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func SetUser(ctx context.Context, user *entities.User, ttl time.Duration) error {
	cachedUser := &CachedUser{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	key := fmt.Sprintf("%s%s", UserCacheKeyPrefix, user.ID.String())
	if ttl == 0 {
		ttl = DefaultUserCacheTTL
	}
	err := RedisSetData(ctx, key, cachedUser, ttl)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"user_id": user.ID.String(),
		"ttl":     ttl.String(),
	}).Debug("User cached successfully")
	return nil
}

func GetUser(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	key := fmt.Sprintf("%s%s", UserCacheKeyPrefix, userID.String())
	var cachedUser CachedUser
	err := RedisGetData(ctx, key, &cachedUser)
	if err != nil {
		return nil, err
	}
	user := &entities.User{
		ID:        cachedUser.ID,
		Email:     cachedUser.Email,
		Name:      cachedUser.Name,
		Role:      cachedUser.Role,
		CreatedAt: cachedUser.CreatedAt,
		UpdatedAt: cachedUser.UpdatedAt,
	}
	logrus.WithField("user_id", userID.String()).Debug("User retrieved from cache")
	return user, nil
}

func DeleteUser(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("%s%s", UserCacheKeyPrefix, userID.String())
	err := RedisDeleteData(ctx, key)
	if err != nil {
		return err
	}
	logrus.WithField("user_id", userID.String()).Debug("User deleted from cache")
	return nil
}
