package service

import (
	"context"

	"kasir-api/internal/model"
	"kasir-api/internal/repository"
)

type TransactionService struct {
	writer repository.TransactionWriter
}

func NewTransactionService(writer repository.TransactionWriter) *TransactionService {
	return &TransactionService{writer: writer}
}

func (s *TransactionService) Checkout(ctx context.Context, req model.CheckoutRequest) (*model.Transaction, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.writer.CreateTransaction(ctx, req.Items)
}
