package data

import "github.com/ksolj/ongaku-api/internal/data/validator"

type Filters struct {
	Page     int
	PageSize int
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	// Potential risk of integer overflow
	// but it's mitigated by the validation rules in ValidateFilters()
	return (f.Page - 1) * f.PageSize
}
