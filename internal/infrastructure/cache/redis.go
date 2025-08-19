package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/naufalfazanadi/finance-manager-go/pkg/config"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

// GetRedisClient returns a Redis client instance (singleton pattern)
func GetRedisClient() *redis.Client {
	redisOnce.Do(func() {
		cfg := config.GetConfig()

		// Create Redis client options
		options := &redis.Options{
			Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			Password:     cfg.Redis.Password,
			DB:           cfg.Redis.DB,
			MaxRetries:   cfg.Redis.MaxRetries,
			PoolSize:     cfg.Redis.PoolSize,
			MinIdleConns: cfg.Redis.MinIdle,
			MaxIdleConns: cfg.Redis.MaxIdle,
			DialTimeout:  time.Duration(cfg.Redis.DialTimeout) * time.Second,
			ReadTimeout:  time.Duration(cfg.Redis.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.Redis.ReadTimeout) * time.Second,
		}

		redisClient = redis.NewClient(options)

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := redisClient.Ping(ctx).Result()
		if err != nil {
			logrus.WithError(err).Error("Failed to connect to Redis")
			// Don't panic, just log the error - fallback to database will handle this
		} else {
			logrus.Info("Successfully connected to Redis")
		}
	})

	return redisClient
}

// IsRedisAvailable checks if Redis is available and healthy
func IsRedisAvailable() bool {
	if redisClient == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	return err == nil
}

// CloseRedisConnection closes the Redis connection
func CloseRedisConnection() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

// RedisSetData sets a value in Redis for a given key and TTL
func RedisSetData(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	client := GetRedisClient()
	if !IsRedisAvailable() {
		return fmt.Errorf("redis not available")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set data in redis: %w", err)
	}
	return nil
}

// RedisGetData gets a value from Redis for a given key and unmarshals into dest
func RedisGetData(ctx context.Context, key string, dest interface{}) error {
	client := GetRedisClient()
	if !IsRedisAvailable() {
		return fmt.Errorf("redis not available")
	}

	data, err := client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found in redis")
		}
		return fmt.Errorf("failed to get data from redis: %w", err)
	}

	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal redis data: %w", err)
	}
	return nil
}

// RedisDeleteData deletes a key from Redis
func RedisDeleteData(ctx context.Context, key string) error {
	client := GetRedisClient()
	if !IsRedisAvailable() {
		return fmt.Errorf("redis not available")
	}

	err := client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key from redis: %w", err)
	}
	return nil
}
