package main

// Product represents a product in the cashier system
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"nama"`
	Price int    `json:"harga"`
	Stock int    `json:"stok"`
}

// Category represents a product category
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"nama"`
	Description string `json:"deskripsi"`
}
