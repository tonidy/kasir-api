package service

import (
	"context"

	"kasir-api/internal/model"
	"kasir-api/internal/repository"
	errorsPkg "kasir-api/pkg/errors"
	"kasir-api/pkg/tracing"
)

type TransactionService struct {
	writer repository.TransactionWriter
}

func NewTransactionService(writer repository.TransactionWriter) *TransactionService {
	return &TransactionService{writer: writer}
}

func (s *TransactionService) Checkout(ctx context.Context, req model.CheckoutRequest) (*model.Transaction, error) {
	ctx, spanEnd := tracing.TraceRequest(ctx, "TransactionService.Checkout", req)
	defer spanEnd(nil, nil)

	if err := req.Validate(); err != nil {
		spanEnd(nil, err)
		return nil, err
	}

	transaction, err := s.writer.CreateTransaction(ctx, req.Items)
	if err != nil {
		spanEnd(nil, err)
		return nil, errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to create transaction")
	}

	spanEnd(transaction, nil)
	return transaction, nil
}
