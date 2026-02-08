package postgres

import (
	"context"
	"database/sql"

	"kasir-api/internal/model"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetTodayReport(ctx context.Context) (*model.ReportSummary, error) {
	var totalRevenue, totalTransaction int
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&totalRevenue, &totalTransaction)
	if err != nil {
		return nil, err
	}

	var topProduct *model.TopProduct
	var name sql.NullString
	var soldQty sql.NullInt64
	err = r.db.QueryRowContext(ctx, `
		SELECT p.name, SUM(td.quantity) as sold_qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.id, p.name
		ORDER BY sold_qty DESC
		LIMIT 1
	`).Scan(&name, &soldQty)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if name.Valid {
		topProduct = &model.TopProduct{
			Name:    name.String,
			SoldQty: int(soldQty.Int64),
		}
	}

	return &model.ReportSummary{
		TotalRevenue:     totalRevenue,
		TotalTransaction: totalTransaction,
		TopProduct:       topProduct,
	}, nil
}

func (r *ReportRepository) GetReportByDateRange(ctx context.Context, startDate, endDate string) (*model.ReportSummary, error) {
	var totalRevenue, totalTransaction int
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&totalRevenue, &totalTransaction)
	if err != nil {
		return nil, err
	}

	var topProduct *model.TopProduct
	var name sql.NullString
	var soldQty sql.NullInt64
	err = r.db.QueryRowContext(ctx, `
		SELECT p.name, SUM(td.quantity) as sold_qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) BETWEEN $1 AND $2
		GROUP BY p.id, p.name
		ORDER BY sold_qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(&name, &soldQty)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if name.Valid {
		topProduct = &model.TopProduct{
			Name:    name.String,
			SoldQty: int(soldQty.Int64),
		}
	}

	return &model.ReportSummary{
		TotalRevenue:     totalRevenue,
		TotalTransaction: totalTransaction,
		TopProduct:       topProduct,
	}, nil
}
