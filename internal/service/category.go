package service

import (
	"context"

	"kasir-api/internal/model"
	"kasir-api/internal/repository"
)

type CategoryService struct {
	reader repository.CategoryReader
	writer repository.CategoryWriter
}

func NewCategoryService(reader repository.CategoryReader, writer repository.CategoryWriter) *CategoryService {
	return &CategoryService{
		reader: reader,
		writer: writer,
	}
}

func (s *CategoryService) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	return s.reader.FindByID(ctx, id)
}

func (s *CategoryService) GetAll(ctx context.Context) ([]domain.Category, error) {
	return s.reader.FindAll(ctx)
}

func (s *CategoryService) Create(ctx context.Context, c domain.Category) (*domain.Category, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return s.writer.Create(ctx, c)
}

func (s *CategoryService) Update(ctx context.Context, id int, c domain.Category) (*domain.Category, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return s.writer.Update(ctx, id, c)
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	return s.writer.Delete(ctx, id)
}
