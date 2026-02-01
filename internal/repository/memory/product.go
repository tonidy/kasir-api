package memory

import (
	"context"
	"sync"

	"kasir-api/internal/model"
)

type ProductRepository struct {
	mu      sync.RWMutex
	data    []model.Product
	nextID  int
	catRepo *CategoryRepository
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		data:   make([]model.Product, 0),
		nextID: 1,
	}
}

func (r *ProductRepository) SetCategoryRepo(catRepo *CategoryRepository) {
	r.catRepo = catRepo
}

func (r *ProductRepository) FindByID(ctx context.Context, id int) (*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.data {
		if p.ID == id {
			result := p

			// Get category if category_id exists
			if p.CategoryID != nil && r.catRepo != nil {
				if cat, err := r.catRepo.FindByID(ctx, *p.CategoryID); err == nil {
					result.Category = cat
				}
			}

			return &result, nil
		}
	}
	return nil, model.ErrNotFound
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]model.Product, 0, len(r.data))
	for _, p := range r.data {
		result := p

		// Get category if category_id exists
		if p.CategoryID != nil && r.catRepo != nil {
			if cat, err := r.catRepo.FindByID(ctx, *p.CategoryID); err == nil {
				result.Category = cat
			}
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *ProductRepository) Create(ctx context.Context, p model.Product) (*model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p.ID = r.nextID
	r.nextID++
	r.data = append(r.data, p)
	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int, p model.Product) (*model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.data {
		if r.data[i].ID == id {
			p.ID = id
			r.data[i] = p
			return &p, nil
		}
	}
	return nil, model.ErrNotFound
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.data {
		if p.ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return model.ErrNotFound
}
