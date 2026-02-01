package postgres

import (
	"context"
	"database/sql"
	"errors"

	"kasir-api/internal/model"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindByID(ctx context.Context, id int) (*model.Category, error) {
	query := `SELECT id, name, description FROM categories WHERE id = $1`

	var c model.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]model.Category, error) {
	query := `SELECT id, name, description FROM categories ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c model.Category) (*model.Category, error) {
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, c.Name, c.Description).Scan(&c.ID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id int, c model.Category) (*model.Category, error) {
	query := `UPDATE categories SET name = $1, description = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, c.Name, c.Description, id)
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

	c.ID = id
	return &c, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`

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
