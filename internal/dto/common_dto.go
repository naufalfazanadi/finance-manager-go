package dto

import (
	"github.com/google/uuid"
)

// Pagination DTOs
type PaginationQuery struct {
	Page  int `json:"page" query:"page" validate:"omitempty,min=1" example:"1"`
	Limit int `json:"limit" query:"limit" validate:"omitempty,min=1,max=100" example:"10"`
}

type PaginationMeta struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
}

// Filter DTOs
type FilterQuery struct {
	Search   string            `json:"search" query:"search" validate:"omitempty,max=255"`
	SortBy   string            `json:"sort_by" query:"sort_by" validate:"omitempty"`
	SortType string            `json:"sort_type" query:"sort_type" validate:"omitempty,oneof=asc desc"`
	Filters  map[string]string `json:"filters" query:"filters"`
}

// Combined Query for pagination and filtering
type QueryParams struct {
	*PaginationQuery
	*FilterQuery
	LoggedUserID uuid.UUID `json:"logged_user_id" param:"logged_user_id" validate:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// ID Parameter validation
type IDParam struct {
	ID uuid.UUID `json:"id" param:"id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// Error Response
type ErrorResponse struct {
	Error   string                 `json:"error" example:"Bad Request"`
	Message string                 `json:"message" example:"Invalid input data"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type PaginationData[T any] struct {
	Data []T             `json:"data"`
	Meta *PaginationMeta `json:"meta"`
}

// GetOffset calculates the offset for database queries
func (p *PaginationQuery) GetOffset() int {
	return (p.Page - 1) * p.Limit
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
