package postgres

import (
	"context"
	"database/sql"
	"errors"

	"kasir-api/internal/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindByID(ctx context.Context, id int) (*domain.ProductWithCategory, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1`

	var p domain.ProductWithCategory
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]domain.ProductWithCategory, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.ProductWithCategory
	for rows.Next() {
		var p domain.ProductWithCategory
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Create(ctx context.Context, p domain.Product) (*domain.Product, error) {
	query := `INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, p.Name, p.Price, p.Stock, p.CategoryID).Scan(&p.ID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int, p domain.Product) (*domain.Product, error) {
	query := `UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5`

	result, err := r.db.ExecContext(ctx, query, p.Name, p.Price, p.Stock, p.CategoryID, id)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, domain.ErrNotFound
	}

	p.ID = id
	return &p, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
