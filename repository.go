package main

// Repository interface for data operations
type Repository interface {
	FindByID(id int) (*Product, bool)
	FindAll() []Product
	Create(p Product) Product
	Update(id int, p Product) (*Product, bool)
	Delete(id int) bool
}

// ProductRepository in-memory implementation
type ProductRepository struct {
	data []Product
}

func (r *ProductRepository) FindByID(id int) (*Product, bool) {
	for _, p := range r.data {
		if p.ID == id {
			return &p, true
		}
	}
	return nil, false
}

func (r *ProductRepository) FindAll() []Product {
	return r.data
}

func (r *ProductRepository) Create(p Product) Product {
	p.ID = len(r.data) + InitialIDOffset
	r.data = append(r.data, p)
	return p
}

func (r *ProductRepository) Update(id int, p Product) (*Product, bool) {
	for i := range r.data {
		if r.data[i].ID == id {
			p.ID = id
			r.data[i] = p
			return &p, true
		}
	}
	return nil, false
}

func (r *ProductRepository) Delete(id int) bool {
	for i, p := range r.data {
		if p.ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return true
		}
	}
	return false
}

// In-memory storage
var repo Repository = &ProductRepository{
	data: []Product{
		{ID: 1, Name: "Indomie Godog", Price: 3500, Stock: 10},
		{ID: 2, Name: "Vit 1L", Price: 3000, Stock: 40},
		{ID: 3, Name: "Kecap Bango", Price: 12000, Stock: 20},
	},
}

// CategoryRepository interface
type CategoryRepository interface {
	FindByID(id int) (*Category, bool)
	FindAll() []Category
	Create(c Category) Category
	Update(id int, c Category) (*Category, bool)
	Delete(id int) bool
}

// InMemoryCategoryRepository implementation
type InMemoryCategoryRepository struct {
	data []Category
}

func (r *InMemoryCategoryRepository) FindByID(id int) (*Category, bool) {
	for _, c := range r.data {
		if c.ID == id {
			return &c, true
		}
	}
	return nil, false
}

func (r *InMemoryCategoryRepository) FindAll() []Category {
	return r.data
}

func (r *InMemoryCategoryRepository) Create(c Category) Category {
	c.ID = len(r.data) + InitialIDOffset
	r.data = append(r.data, c)
	return c
}

func (r *InMemoryCategoryRepository) Update(id int, c Category) (*Category, bool) {
	for i := range r.data {
		if r.data[i].ID == id {
			c.ID = id
			r.data[i] = c
			return &c, true
		}
	}
	return nil, false
}

func (r *InMemoryCategoryRepository) Delete(id int) bool {
	for i, c := range r.data {
		if c.ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return true
		}
	}
	return false
}

var categoryRepo CategoryRepository = &InMemoryCategoryRepository{
	data: []Category{
		{ID: 1, Name: "Makanan", Description: "Produk makanan"},
		{ID: 2, Name: "Minuman", Description: "Produk minuman"},
	},
}
