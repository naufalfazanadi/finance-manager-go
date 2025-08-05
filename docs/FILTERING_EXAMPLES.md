# User Filtering and Sorting Examples

This document shows how to use the enhanced GetAll method with filtering and sorting capabilities.

## API Usage Examples

### 1. Basic Pagination
```
GET /api/users?page=2&limit=10
```

### 2. Search Users
```
GET /api/users?search=john
```
This will search for users where name or email contains "john" (case-insensitive).

### 3. Filter by Role
```
GET /api/users?role=admin
```

### 4. Sort Users
```
GET /api/users?sort_by=created_at&sort_type=desc
```
Available sort columns: `name`, `email`, `role`, `created_at`, `updated_at`

### 5. Date Range Filtering
```
GET /api/users?created_after=2024-01-01&created_before=2024-12-31
```

### 6. Combined Filtering
```
GET /api/users?search=john&role=user&sort_by=name&sort_type=asc&page=1&limit=20
```

## Implementation Details

### Supported Filters
- `role`: Filter by user role (admin, user, etc.)
- `name`: Exact name match
- `email`: Exact email match
- `created_after`: Users created after this date
- `created_before`: Users created before this date

### Search Functionality
- Searches in both `name` and `email` fields
- Case-insensitive search using ILIKE
- Supports partial matches

### Security Features
- SQL injection protection through whitelisted columns
- Only safe column names are allowed for filtering and sorting
- Parameterized queries for all filters

## Code Implementation

The filtering is implemented in the repository layer:

```go
func (r *userRepositoryImpl) GetAll(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.User, error) {
    var users []*entities.User
    query := r.db.WithContext(ctx)

    // Apply search
    if queryParams.HasSearch() {
        searchTerm := "%" + queryParams.Search + "%"
        query = query.Where("name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
    }

    // Apply filters
    if queryParams.HasFilters() {
        for key, value := range queryParams.Filters {
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

    // Apply sorting with whitelist
    if queryParams.HasSort() {
        allowedSortColumns := map[string]bool{
            "name": true, "email": true, "role": true,
            "created_at": true, "updated_at": true,
        }

        if allowedSortColumns[queryParams.SortBy] {
            query = query.Order(queryParams.SortBy + " " + queryParams.SortType)
        }
    }

    // Apply pagination
    query = query.Limit(queryParams.Limit).Offset(queryParams.GetOffset())

    return query.Find(&users).Error
}
```

## Frontend Usage Example

```javascript
// JavaScript example for frontend
const fetchUsers = async (filters = {}) => {
    const params = new URLSearchParams();

    // Add pagination
    if (filters.page) params.append('page', filters.page);
    if (filters.limit) params.append('limit', filters.limit);

    // Add search
    if (filters.search) params.append('search', filters.search);

    // Add sorting
    if (filters.sortBy) params.append('sort_by', filters.sortBy);
    if (filters.sortType) params.append('sort_type', filters.sortType);

    // Add filters
    if (filters.role) params.append('role', filters.role);
    if (filters.createdAfter) params.append('created_after', filters.createdAfter);

    const response = await fetch(`/api/users?${params.toString()}`);
    return response.json();
};

// Usage examples
fetchUsers({ search: 'john', role: 'admin', sortBy: 'created_at', sortType: 'desc' });
fetchUsers({ role: 'user', page: 2, limit: 20 });
```
