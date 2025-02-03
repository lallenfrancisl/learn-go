package data

import (
	"math"
	"strings"

	"github.com/lallenfrancisl/greenlight-api/internal/validator"
)

type BaseFilter struct {
	Page         int      `json:"page"`
	PageSize     int      `json:"page_size"`
	Sort         string   `json:"sort"`
	SortSafelist []string `json:"-"`
}

func (f *BaseFilter) SortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f BaseFilter) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f BaseFilter) Limit() int {
	return f.PageSize
}

func (f BaseFilter) Offset() int {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords int, page int, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
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
