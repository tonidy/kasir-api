package service

import (
	"context"

	"kasir-api/internal/model"
	"kasir-api/internal/repository"
	errorsPkg "kasir-api/pkg/errors"
	"kasir-api/pkg/tracing"
)

type ReportService struct {
	reader repository.ReportReader
}

func NewReportService(reader repository.ReportReader) *ReportService {
	return &ReportService{reader: reader}
}

func (s *ReportService) GetTodayReport(ctx context.Context) (*model.ReportSummary, error) {
	ctx, spanEnd := tracing.TraceRequest(ctx, "ReportService.GetTodayReport", nil)
	defer spanEnd(nil, nil)

	report, err := s.reader.GetTodayReport(ctx)
	if err != nil {
		spanEnd(nil, err)
		return nil, errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to get today's report")
	}

	spanEnd(report, nil)
	return report, nil
}

func (s *ReportService) GetReportByDateRange(ctx context.Context, startDate, endDate string) (*model.ReportSummary, error) {
	ctx, spanEnd := tracing.TraceRequest(ctx, "ReportService.GetReportByDateRange", map[string]interface{}{"startDate": startDate, "endDate": endDate})
	defer spanEnd(nil, nil)

	report, err := s.reader.GetReportByDateRange(ctx, startDate, endDate)
	if err != nil {
		spanEnd(nil, err)
		return nil, errorsPkg.Wrap(err, errorsPkg.ErrorTypeInternal, "failed to get report by date range")
	}

	spanEnd(report, nil)
	return report, nil
}
