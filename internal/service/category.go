package service

import (
	"context"

	"kasir-api/internal/model"
	"kasir-api/internal/repository"
	errorsPkg "kasir-api/pkg/errors"
	"kasir-api/pkg/tracing"
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

func (s *CategoryService) GetByID(ctx context.Context, id int) (*model.Category, error) {
	ctx, spanEnd := tracing.TraceRequest(ctx, "CategoryService.GetByID", map[string]interface{}{"id": id})
	defer spanEnd(nil, nil)

	category, err := s.reader.FindByID(ctx, id)
	if err != nil {
		spanEnd(nil, err)
		return nil, err
	}

	spanEnd(category, nil)
	return category, nil
}

func (s *CategoryService) GetAll(ctx context.Context) ([]model.Category, error) {
	ctx, spanEnd := tracing.TraceRequest(ctx, "CategoryService.GetAll", nil)
	defer spanEnd(nil, nil)

	categories, err := s.reader.FindAll(ctx)
	if err != nil {
		spanEnd(nil, err)
		return nil, errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to get all categories")
	}

	spanEnd(categories, nil)
	return categories, nil
}

func (s *CategoryService) Create(ctx context.Context, c model.Category) (*model.Category, error) {
	ctx, spanEnd := tracing.TraceRequest(ctx, "CategoryService.Create", c)
	defer spanEnd(nil, nil)

	if err := c.Validate(); err != nil {
		spanEnd(nil, err)
		return nil, err
	}

	created, err := s.writer.Create(ctx, c)
	if err != nil {
		spanEnd(nil, err)
		return nil, errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to create category")
	}

	spanEnd(created, nil)
	return created, nil
}

func (s *CategoryService) Update(ctx context.Context, id int, c model.Category) (*model.Category, error) {
	ctx, spanEnd := tracing.TraceRequest(ctx, "CategoryService.Update", map[string]interface{}{"id": id, "category": c})
	defer spanEnd(nil, nil)

	if err := c.Validate(); err != nil {
		spanEnd(nil, err)
		return nil, err
	}

	updated, err := s.writer.Update(ctx, id, c)
	if err != nil {
		spanEnd(nil, err)
		return nil, errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to update category")
	}

	spanEnd(updated, nil)
	return updated, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	ctx, spanEnd := tracing.TraceRequest(ctx, "CategoryService.Delete", map[string]interface{}{"id": id})
	defer spanEnd(nil, nil)

	err := s.writer.Delete(ctx, id)
	if err != nil {
		spanEnd(nil, err)
		return errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to delete category")
	}

	spanEnd(nil, nil)
	return nil
}
