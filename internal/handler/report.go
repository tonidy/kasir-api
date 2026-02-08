package handler

import (
	"context"
	"net/http"

	"kasir-api/internal/model"
	"kasir-api/pkg/httputil"
)

type ReportService interface {
	GetTodayReport(ctx context.Context) (*model.ReportSummary, error)
	GetReportByDateRange(ctx context.Context, startDate, endDate string) (*model.ReportSummary, error)
}

type ReportHandler struct {
	svc ReportService
}

func NewReportHandler(svc ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) Today(w http.ResponseWriter, r *http.Request) {
	report, err := h.svc.GetTodayReport(r.Context())
	if err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, report)
}

func (h *ReportHandler) ByDateRange(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "start_date and end_date are required",
		})
		return
	}

	report, err := h.svc.GetReportByDateRange(r.Context(), startDate, endDate)
	if err != nil {
		httputil.WriteError(w, httputil.ErrorStatus(err), err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, report)
}
