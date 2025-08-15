package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/sirupsen/logrus"
)

// RateLimiterConfig holds configuration for rate limiting
type RateLimiterConfig struct {
	Max        int           // Maximum number of requests
	Expiration time.Duration // Time window
	KeySuffix  string        // Suffix for the key to differentiate between different limiters
	Message    string        // Custom message when limit is reached
}

// DefaultRateLimiter returns the default rate limiter for general routes
func DefaultRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,             // 100 requests
		Expiration: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit per IP address
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Log rate limit violation
			logger.LogError(
				"DefaultRateLimiter",
				"Rate limit exceeded for general routes",
				nil,
				logrus.Fields{
					"ip":         c.IP(),
					"method":     c.Method(),
					"path":       c.Path(),
					"user_agent": c.Get("User-Agent"),
					"limit":      100,
					"window":     "1 minute",
				},
			)
			return c.Status(429).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
		},
	})
}

// APIRateLimiter returns a more restrictive rate limiter for API routes
func APIRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        50,              // 50 requests
		Expiration: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + "-api" // Different key for API endpoints
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Log API rate limit violation
			logger.LogError(
				"APIRateLimiter",
				"API rate limit exceeded",
				nil,
				logrus.Fields{
					"ip":         c.IP(),
					"method":     c.Method(),
					"path":       c.Path(),
					"user_agent": c.Get("User-Agent"),
					"limit":      50,
					"window":     "1 minute",
				},
			)
			return c.Status(429).JSON(fiber.Map{
				"error":       "API rate limit exceeded",
				"message":     "Too many API requests. Please try again later.",
				"retry_after": "60 seconds",
			})
		},
	})
}

// AuthRateLimiter returns a very restrictive rate limiter for authentication routes
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,              // 10 requests
		Expiration: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + "-auth" // Different key for auth endpoints
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Log authentication rate limit violation with high severity
			logger.LogError(
				"AuthRateLimiter",
				"Authentication rate limit exceeded - potential brute force attack",
				nil,
				logrus.Fields{
					"ip":         c.IP(),
					"method":     c.Method(),
					"path":       c.Path(),
					"user_agent": c.Get("User-Agent"),
					"limit":      10,
					"window":     "1 minute",
					"severity":   "HIGH",
					"alert_type": "SECURITY",
				},
			)
			return c.Status(429).JSON(fiber.Map{
				"error":       "Authentication rate limit exceeded",
				"message":     "Too many authentication attempts. Please try again later.",
				"retry_after": "60 seconds",
			})
		},
	})
}

// CustomRateLimiter creates a rate limiter with custom configuration
func CustomRateLimiter(config RateLimiterConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        config.Max,
		Expiration: config.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			key := c.IP()
			if config.KeySuffix != "" {
				key += "-" + config.KeySuffix
			}
			return key
		},
		LimitReached: func(c *fiber.Ctx) error {
			message := config.Message
			if message == "" {
				message = "Rate limit exceeded. Please try again later."
			}
			// Log custom rate limit violation
			logger.LogError(
				"CustomRateLimiter",
				"Custom rate limit exceeded",
				nil,
				logrus.Fields{
					"ip":         c.IP(),
					"method":     c.Method(),
					"path":       c.Path(),
					"user_agent": c.Get("User-Agent"),
					"limit":      config.Max,
					"window":     config.Expiration.String(),
					"key_suffix": config.KeySuffix,
				},
			)
			return c.Status(429).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"message": message,
			})
		},
	})
}
