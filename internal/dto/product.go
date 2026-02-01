package dto

// ProductResponse represents product data with category information for API responses
type ProductResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Price        int     `json:"price"`
	Stock        int     `json:"stock"`
	CategoryID   *int    `json:"category_id,omitempty"`
	CategoryName *string `json:"category_name,omitempty"`
}
