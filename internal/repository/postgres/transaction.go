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

	totalAmount := 0
	details := make([]model.TransactionDetail, 0, len(items))

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRowContext(ctx, "SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).
			Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: product id %d not found", model.ErrNotFound, item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("%w: insufficient stock for product %s (available: %d, requested: %d)",
				model.ErrValidation, productName, stock, item.Quantity)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.ExecContext(ctx, "UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	var createdAt sql.NullTime
	err = tx.QueryRowContext(ctx, "INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at", totalAmount).
		Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		var detailID int
		err = tx.QueryRowContext(ctx,
			"INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID
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
