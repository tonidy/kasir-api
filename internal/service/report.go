package service

import (
	"context"

	"kasir-api/internal/model"
	"kasir-api/internal/repository"
)

type ReportService struct {
	reader repository.ReportReader
}

func NewReportService(reader repository.ReportReader) *ReportService {
	return &ReportService{reader: reader}
}

func (s *ReportService) GetTodayReport(ctx context.Context) (*model.ReportSummary, error) {
	return s.reader.GetTodayReport(ctx)
}

func (s *ReportService) GetReportByDateRange(ctx context.Context, startDate, endDate string) (*model.ReportSummary, error) {
	return s.reader.GetReportByDateRange(ctx, startDate, endDate)
}
