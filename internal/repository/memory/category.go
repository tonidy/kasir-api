package memory

import (
	"context"
	"sync"

	"kasir-api/internal/model"
)

type CategoryRepository struct {
	mu     sync.RWMutex
	data   []model.Category
	nextID int
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		data:   make([]model.Category, 0),
		nextID: 1,
	}
}

func (r *CategoryRepository) FindByID(ctx context.Context, id int) (*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, c := range r.data {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, model.ErrNotFound
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.data, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c model.Category) (*model.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	c.ID = r.nextID
	r.nextID++
	r.data = append(r.data, c)
	return &c, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id int, c model.Category) (*model.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.data {
		if r.data[i].ID == id {
			c.ID = id
			r.data[i] = c
			return &c, nil
		}
	}
	return nil, model.ErrNotFound
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, c := range r.data {
		if c.ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return model.ErrNotFound
}
