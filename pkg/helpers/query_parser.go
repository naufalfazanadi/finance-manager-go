package helpers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
)

// ParsePaginationQuery parses pagination parameters from Fiber context
func ParsePaginationQuery(c *fiber.Ctx) *dto.PaginationQuery {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	return &dto.PaginationQuery{
		Page:  page,
		Limit: limit,
	}
}

// ParseFilterQuery parses filter parameters from Fiber context
func ParseFilterQuery(c *fiber.Ctx) *dto.FilterQuery {
	search := c.Query("search")
	sortBy := c.Query("sort_by")
	sortType := strings.ToLower(c.Query("sort_type"))

	// Validate sort direction
	if sortType != "asc" && sortType != "desc" {
		sortType = "asc"
	}

	// Parse custom filters (e.g., ?name=john&role=user)
	filters := make(map[string]string)
	queryMap := c.Queries()
	excludedParams := map[string]bool{
		"page":      true,
		"limit":     true,
		"search":    true,
		"sort_by":   true,
		"sort_type": true,
	}

	for key, value := range queryMap {
		if !excludedParams[key] && value != "" {
			filters[key] = value
		}
	}

	return &dto.FilterQuery{
		Search:   search,
		SortBy:   sortBy,
		SortType: sortType,
		Filters:  filters,
	}
}

// ParseQueryParams parses both pagination and filter parameters
func ParseQueryParams(c *fiber.Ctx) *dto.QueryParams {
	return &dto.QueryParams{
		PaginationQuery: ParsePaginationQuery(c),
		FilterQuery:     ParseFilterQuery(c),
	}
}

// GetOffset calculates the offset for database queries
func GetOffset(p *dto.PaginationQuery) int {
	return (p.Page - 1) * p.Limit
}

// NewPaginationMeta creates pagination metadata
func NewPaginationMeta(page, limit int, total int64) *dto.PaginationMeta {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}

	return &dto.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// HasSearch checks if search query is provided
func HasSearch(f *dto.FilterQuery) bool {
	return f.Search != ""
}

// HasFilters checks if any custom filters are provided
func HasFilters(f *dto.FilterQuery) bool {
	return len(f.Filters) > 0
}

// HasSort checks if sorting is provided
func HasSort(f *dto.FilterQuery) bool {
	return f.SortBy != ""
}

// GetFilterValue gets a specific filter value
func GetFilterValue(f *dto.FilterQuery, key string) (string, bool) {
	value, exists := f.Filters[key]
	return value, exists
}
