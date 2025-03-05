package data

import (
	"strings"

	"slices"

	"greenlight.mazavrbazavr.ru/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string // holds the supported sort values
}

// Checks that the client-provided Sort field matches one of the entries in our safelist
// and if it does, extracts the column name from the Sort field by stripping the leading
// hyphen character (if one exists).
func (f Filters) sortColumn() string {
	if slices.Contains(f.SortSafelist, f.Sort) {
		return strings.TrimPrefix(f.Sort, "-")
	}

	// Panic if the client-provided Sort value doesn’t match one of the entries
	// in the safelist. In theory this shouldn’t happen — the Sort value should have
	// already been checked by calling the ValidateFilters() function — but this is
	// a sensible failsafe to help stop a SQL injection attack occurring.
	panic("unsafe sort parameter: " + f.Sort)
}

// Returns the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

// Get the DB LIMIT value.
func (f Filters) limit() int {
	return f.PageSize
}

// Get the DB OFFSET value.
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// Checks that the Filters struct contains valid values.
func ValidateFilters(v *validator.Validator, f Filters) {
	// Check that the page and page_size parameters contain sensible values.
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	// Check that the sort parameter matches a value in the safelist.
	v.Check(
		validator.PermittedValue(f.Sort, f.SortSafelist...),
		"sort",
		"invalid sort value",
	)
}
