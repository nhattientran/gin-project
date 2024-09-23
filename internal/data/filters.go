package data

import "gin-project/internal/validator"

type Filters struct {
	Page         int      `form:"page"`
	PageSize     int      `form:"page_size"`
	Sort         string   `form:"sort"`
	SortSafelist []string `json:"-"`
}

func NewFilters() Filters {
	return Filters{
		Page:     1,
		PageSize: 10,
		Sort:     "title",
	}
}
func ValidateFilters(v *validator.Validator, f Filters) {

	// Check that the page and page_size parameters contain sensible values.
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	// Check that the sort parameter matches a value in the safelist.
	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
