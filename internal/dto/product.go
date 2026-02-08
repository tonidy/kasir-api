package dto

// CategoryResponse represents category data for API responses
type CategoryResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProductResponse represents product data with category information for API responses
type ProductResponse struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	Price    int               `json:"price"`
	Stock    int               `json:"stock"`
	Active   bool              `json:"active"`
	Category *CategoryResponse `json:"category,omitempty"`
}
