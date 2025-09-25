package model

type ProductSearchQuery struct {
	Query      string  `form:"q"`
	CategoryID uint    `form:"category"`
	PriceMin   float64 `form:"price_min"`
	PriceMax   float64 `form:"price_max"`
	Sort       string  `form:"sort"` // e.g., price_asc, price_desc
	Page       int     `form:"page"`
	PageSize   int     `form:"page_size"`
}
