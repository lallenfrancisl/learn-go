package data

import "github.com/lallenfrancisl/greenlight-api/internal/validator"

type BaseFilter struct {
	Page         int      `json:"page"`
	PageSize     int      `json:"page_size"`
	Sort         string   `json:"sort"`
	SortSafelist []string `json:"-"`
}

func ValidateFilters(v *validator.Validator, f BaseFilter) {
	v.Check(validator.GreaterThan(f.Page, 0), "page", "must be greater than zero")
	v.Check(
		validator.Max(f.Page, 10_000_000), "page", "must be a maximum of 10 million",
	)

	v.Check(
		validator.GreaterThan(f.PageSize, 0), "page_size", "must be greater than zero",
	)
	v.Check(
		validator.Max(f.PageSize, 100), "page_size", "must be a maximum of 100",
	)

	v.Check(
		validator.In(f.Sort, f.SortSafelist...),
		"sort", "invalid sort value",
	)
}
