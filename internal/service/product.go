package service

import (
	"context"

	"kasir-api/internal/model"
	"kasir-api/internal/repository"
	errorsPkg "kasir-api/pkg/errors"
	"kasir-api/pkg/tracing"
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
	ctx, spanEnd := tracing.TraceRequest(ctx, "ProductService.GetByID", map[string]interface{}{"id": id})
	defer spanEnd(nil, nil)

	product, err := s.reader.FindByID(ctx, id)
	if err != nil {
		spanEnd(nil, err)
		return nil, err
	}

	spanEnd(product, nil)
	return product, nil
}

func (s *ProductService) GetAll(ctx context.Context) ([]model.Product, error) {
	ctx, spanEnd := tracing.TraceRequest(ctx, "ProductService.GetAll", nil)
	defer spanEnd(nil, nil)

	products, err := s.reader.FindAll(ctx)
	if err != nil {
		spanEnd(nil, err)
		return nil, errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to get all products")
	}

	spanEnd(products, nil)
	return products, nil
}

func (s *ProductService) GetByFilters(ctx context.Context, name string, active *bool) ([]model.Product, error) {
	ctx, spanEnd := tracing.TraceRequest(ctx, "ProductService.GetByFilters", map[string]interface{}{"name": name, "active": active})
	defer spanEnd(nil, nil)

	products, err := s.reader.FindByFilters(ctx, name, active)
	if err != nil {
		spanEnd(nil, err)
		return nil, errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to get products by filters")
	}

	spanEnd(products, nil)
	return products, nil
}

func (s *ProductService) Create(ctx context.Context, p model.Product) (*model.Product, error) {
	ctx, spanEnd := tracing.TraceRequest(ctx, "ProductService.Create", p)
	defer spanEnd(nil, nil)

	if err := p.Validate(); err != nil {
		spanEnd(nil, err)
		return nil, err
	}

	created, err := s.writer.Create(ctx, p)
	if err != nil {
		spanEnd(nil, err)
		return nil, errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to create product")
	}

	spanEnd(created, nil)
	return created, nil
}

func (s *ProductService) Update(ctx context.Context, id int, p model.Product) (*model.Product, error) {
	ctx, spanEnd := tracing.TraceRequest(ctx, "ProductService.Update", map[string]interface{}{"id": id, "product": p})
	defer spanEnd(nil, nil)

	if err := p.Validate(); err != nil {
		spanEnd(nil, err)
		return nil, err
	}

	updated, err := s.writer.Update(ctx, id, p)
	if err != nil {
		spanEnd(nil, err)
		return nil, errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to update product")
	}

	spanEnd(updated, nil)
	return updated, nil
}

func (s *ProductService) Delete(ctx context.Context, id int) error {
	ctx, spanEnd := tracing.TraceRequest(ctx, "ProductService.Delete", map[string]interface{}{"id": id})
	defer spanEnd(nil, nil)

	err := s.writer.Delete(ctx, id)
	if err != nil {
		spanEnd(nil, err)
		return errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to delete product")
	}

	spanEnd(nil, nil)
	return nil
}
