package service

import (
	"context"

	"kasir-api/internal/model"
	"kasir-api/internal/repository"
)

type ProductService struct {
	reader repository.ProductReader
	writer repository.ProductWriter
}

func NewProductService(reader repository.ProductReader, writer repository.ProductWriter) *ProductService {
	return &ProductService{
		reader: reader,
		writer: writer,
	}
}

func (s *ProductService) GetByID(ctx context.Context, id int) (*model.Product, error) {
	return s.reader.FindByID(ctx, id)
}

func (s *ProductService) GetAll(ctx context.Context) ([]model.Product, error) {
	return s.reader.FindAll(ctx)
}

func (s *ProductService) GetByFilters(ctx context.Context, name string, active *bool) ([]model.Product, error) {
	return s.reader.FindByFilters(ctx, name, active)
}

func (s *ProductService) Create(ctx context.Context, p model.Product) (*model.Product, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return s.writer.Create(ctx, p)
}

func (s *ProductService) Update(ctx context.Context, id int, p model.Product) (*model.Product, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return s.writer.Update(ctx, id, p)
}

func (s *ProductService) Delete(ctx context.Context, id int) error {
	return s.writer.Delete(ctx, id)
}
