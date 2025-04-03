package usecases

import (
	"context"
	"go-metrics/internal/converters"
	"go-metrics/internal/domain"
	"strings"
)

type MetricListHTMLService interface {
	List(ctx context.Context) ([]*domain.Metric, error)
}

type MetricListHTMLUsecase struct {
	svc MetricListHTMLService
}

func NewMetricListHTMLUsecase(svc MetricListHTMLService) *MetricListHTMLUsecase {
	return &MetricListHTMLUsecase{svc: svc}
}

func (uc *MetricListHTMLUsecase) Execute(
	ctx context.Context,
) (*MetricListHTMLResponse, error) {
	metrics, err := uc.svc.List(ctx)
	if err != nil {
		return nil, err
	}
	return NewMetricListHTMLResponse(metrics), nil
}

type MetricListHTMLResponse struct {
	HTML string
}

func NewMetricListHTMLResponse(metrics []*domain.Metric) *MetricListHTMLResponse {
	var sb strings.Builder
	sb.WriteString("<html><head><title>Metrics</title></head><body>")
	sb.WriteString("<h1>Metrics List</h1>")
	sb.WriteString("<table border='1' style='border-collapse: collapse; width: 50%;'>")
	sb.WriteString("<tr><th>ID</th><th>Value</th></tr>")
	for _, metric := range metrics {
		sb.WriteString("<tr>")
		sb.WriteString("<td>" + metric.ID + "</td>")
		sb.WriteString("<td>")

		if metric.Type == domain.Counter && metric.Delta != nil {
			sb.WriteString(converters.FormatInt64(*metric.Delta))
		} else if metric.Type == domain.Gauge && metric.Value != nil {
			sb.WriteString(converters.FormatFloat64(*metric.Value))
		} else {
			sb.WriteString("N/A")
		}

		sb.WriteString("</td></tr>")
	}
	sb.WriteString("</table></body></html>")
	return &MetricListHTMLResponse{HTML: sb.String()}
}
