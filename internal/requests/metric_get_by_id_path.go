package requests

import (
	"go-metrics/internal/domain"
	"go-metrics/internal/validation"
)

type MetricGetByIDPathRequest struct {
	mtype string
	name  string
}

func NewMetricGetByIDPathRequest(
	mtype string,
	name string,
) *MetricGetByIDPathRequest {
	return &MetricGetByIDPathRequest{
		mtype: mtype,
		name:  name,
	}
}

func (r *MetricGetByIDPathRequest) Validate() error {
	err := validation.ValidateType(r.mtype)
	if err != nil {
		return err
	}
	err = validation.ValidateName(r.name)
	if err != nil {
		return err
	}
	return nil
}

func (r *MetricGetByIDPathRequest) ToDomain() (*domain.MetricID, error) {
	return &domain.MetricID{
		ID:    r.name,
		MType: r.mtype,
	}, nil
}
