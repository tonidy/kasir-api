package memory

import (
	"context"
	"sync"

	"kasir-api/internal/model"
)

type ProductRepository struct {
	mu      sync.RWMutex
	data    []domain.Product
	nextID  int
	catRepo *CategoryRepository
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		data:   make([]domain.Product, 0),
		nextID: 1,
	}
}

func (r *ProductRepository) SetCategoryRepo(catRepo *CategoryRepository) {
	r.catRepo = catRepo
}

func (r *ProductRepository) FindByID(ctx context.Context, id int) (*domain.ProductWithCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.data {
		if p.ID == id {
			result := domain.ProductWithCategory{
				Product: p,
			}

			// Get category name if category_id exists
			if p.CategoryID != nil && r.catRepo != nil {
				if cat, err := r.catRepo.FindByID(ctx, *p.CategoryID); err == nil {
					result.CategoryName = &cat.Name
				}
			}

			return &result, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]domain.ProductWithCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]domain.ProductWithCategory, 0, len(r.data))
	for _, p := range r.data {
		result := domain.ProductWithCategory{
			Product: p,
		}

		// Get category name if category_id exists
		if p.CategoryID != nil && r.catRepo != nil {
			if cat, err := r.catRepo.FindByID(ctx, *p.CategoryID); err == nil {
				result.CategoryName = &cat.Name
			}
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *ProductRepository) Create(ctx context.Context, p domain.Product) (*domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p.ID = r.nextID
	r.nextID++
	r.data = append(r.data, p)
	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int, p domain.Product) (*domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.data {
		if r.data[i].ID == id {
			p.ID = id
			r.data[i] = p
			return &p, nil
		}
	}
	return nil, domain.ErrNotFound
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
	return domain.ErrNotFound
}
