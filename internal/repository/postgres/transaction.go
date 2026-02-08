package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"kasir-api/internal/model"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, items []model.CheckoutItem) (*model.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Batch fetch products with FOR UPDATE to lock rows
	productIDs := make([]any, len(items))
	itemMap := make(map[int]int) // product_id -> quantity
	for i, item := range items {
		productIDs[i] = item.ProductID
		itemMap[item.ProductID] = item.Quantity
	}

	placeholders := ""
	for i := range productIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("SELECT id, name, price, stock FROM products WHERE id IN (%s) FOR UPDATE", placeholders)
	rows, err := tx.QueryContext(ctx, query, productIDs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type productInfo struct {
		id    int
		name  string
		price int
		stock int
	}
	products := make(map[int]productInfo)
	for rows.Next() {
		var p productInfo
		if err := rows.Scan(&p.id, &p.name, &p.price, &p.stock); err != nil {
			return nil, err
		}
		products[p.id] = p
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Validate all products exist and have sufficient stock
	totalAmount := 0
	details := make([]model.TransactionDetail, 0, len(items))

	for _, item := range items {
		product, exists := products[item.ProductID]
		if !exists {
			return nil, fmt.Errorf("%w: product id %d not found", model.ErrNotFound, item.ProductID)
		}

		if product.stock < item.Quantity {
			return nil, fmt.Errorf("%w: insufficient stock for product %s (available: %d, requested: %d)",
				model.ErrValidation, product.name, product.stock, item.Quantity)
		}

		subtotal := product.price * item.Quantity
		totalAmount += subtotal

		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: product.name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Batch update stock
	for productID, quantity := range itemMap {
		_, err = tx.ExecContext(ctx, "UPDATE products SET stock = stock - $1 WHERE id = $2", quantity, productID)
		if err != nil {
			return nil, err
		}
	}

	var transactionID int
	var createdAt sql.NullTime
	err = tx.QueryRowContext(ctx, "INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at", totalAmount).
		Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	// Batch insert transaction details with RETURNING
	if len(details) > 0 {
		query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
		args := make([]any, 0, len(details)*4)

		for i, detail := range details {
			if i > 0 {
				query += ", "
			}
			offset := i * 4
			query += fmt.Sprintf("($%d, $%d, $%d, $%d)", offset+1, offset+2, offset+3, offset+4)
			args = append(args, transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
			details[i].TransactionID = transactionID
		}
		query += " RETURNING id"

		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		i := 0
		for rows.Next() {
			if err := rows.Scan(&details[i].ID); err != nil {
				return nil, err
			}
			i++
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   createdAt.Time,
		Details:     details,
	}, nil
}
