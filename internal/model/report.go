package model

type ReportSummary struct {
	TotalRevenue     int         `json:"total_revenue"`
	TotalTransaction int         `json:"total_transaction"`
	TopProduct       *TopProduct `json:"top_product"`
}

type TopProduct struct {
	Name    string `json:"name"`
	SoldQty int    `json:"sold_qty"`
}
