package dto

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Pagination DTOs
type PaginationQuery struct {
	Page  int `json:"page" query:"page" example:"1"`
	Limit int `json:"limit" query:"limit" example:"10"`
}

type PaginationMeta struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
}

// Filter DTOs
type FilterQuery struct {
	Search   string            `json:"search" query:"search"`
	SortBy   string            `json:"sort_by" query:"sort_by"`
	SortType string            `json:"sort_type" query:"sort_type"`
	Filters  map[string]string `json:"filters" query:"filters"`
}

// Combined Query for pagination and filtering
type QueryParams struct {
	*PaginationQuery
	*FilterQuery
}

// ParsePaginationQuery parses pagination parameters from Fiber context
func ParsePaginationQuery(c *fiber.Ctx) *PaginationQuery {
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

	return &PaginationQuery{
		Page:  page,
		Limit: limit,
	}
}

// ParseFilterQuery parses filter parameters from Fiber context
func ParseFilterQuery(c *fiber.Ctx) *FilterQuery {
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

	return &FilterQuery{
		Search:   search,
		SortBy:   sortBy,
		SortType: sortType,
		Filters:  filters,
	}
}

// ParseQueryParams parses both pagination and filter parameters
func ParseQueryParams(c *fiber.Ctx) *QueryParams {
	return &QueryParams{
		PaginationQuery: ParsePaginationQuery(c),
		FilterQuery:     ParseFilterQuery(c),
	}
}

// GetOffset calculates the offset for database queries
func (p *PaginationQuery) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// NewPaginationMeta creates pagination metadata
func NewPaginationMeta(page, limit int, total int64) *PaginationMeta {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}

	return &PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// HasSearch checks if search query is provided
func (f *FilterQuery) HasSearch() bool {
	return f.Search != ""
}

// HasFilters checks if any custom filters are provided
func (f *FilterQuery) HasFilters() bool {
	return len(f.Filters) > 0
}

// HasSort checks if sorting is provided
func (f *FilterQuery) HasSort() bool {
	return f.SortBy != ""
}

// GetFilterValue gets a specific filter value
func (f *FilterQuery) GetFilterValue(key string) (string, bool) {
	value, exists := f.Filters[key]
	return value, exists
}
