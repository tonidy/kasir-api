package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kasir-api/internal/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindByID(ctx context.Context, id int) (*model.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.active, p.category_id, c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1`

	var p model.Product
	var catID sql.NullInt64
	var catName, catDesc sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Active, &p.CategoryID, &catID, &catName, &catDesc)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	if catID.Valid {
		p.Category = &model.Category{
			ID:          int(catID.Int64),
			Name:        catName.String,
			Description: catDesc.String,
		}
	}

	return &p, nil
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.active, p.category_id, c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		var catID sql.NullInt64
		var catName, catDesc sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Active, &p.CategoryID, &catID, &catName, &catDesc); err != nil {
			return nil, err
		}

		if catID.Valid {
			p.Category = &model.Category{
				ID:          int(catID.Int64),
				Name:        catName.String,
				Description: catDesc.String,
			}
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) FindByFilters(ctx context.Context, name string, active *bool) ([]model.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.active, p.category_id, c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE 1=1`

	args := []interface{}{}
	argPos := 1

	if name != "" {
		query += fmt.Sprintf(" AND p.name ILIKE $%d", argPos)
		args = append(args, "%"+name+"%")
		argPos++
	}

	if active != nil {
		query += fmt.Sprintf(" AND p.active = $%d", argPos)
		args = append(args, *active)
	}

	query += " ORDER BY p.id"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		var catID sql.NullInt64
		var catName, catDesc sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Active, &p.CategoryID, &catID, &catName, &catDesc); err != nil {
			return nil, err
		}

		if catID.Valid {
			p.Category = &model.Category{
				ID:          int(catID.Int64),
				Name:        catName.String,
				Description: catDesc.String,
			}
		}

		products = append(products, p)
	}

	return products, rows.Err()
}

func (r *ProductRepository) Create(ctx context.Context, p model.Product) (*model.Product, error) {
	query := `INSERT INTO products (name, price, stock, active, category_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, p.Name, p.Price, p.Stock, p.Active, p.CategoryID).Scan(&p.ID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int, p model.Product) (*model.Product, error) {
	query := `UPDATE products SET name = $1, price = $2, stock = $3, active = $4, category_id = $5 WHERE id = $6`

	result, err := r.db.ExecContext(ctx, query, p.Name, p.Price, p.Stock, p.Active, p.CategoryID, id)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, model.ErrNotFound
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
		return model.ErrNotFound
	}

	return nil
}
