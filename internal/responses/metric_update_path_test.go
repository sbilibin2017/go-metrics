package responses

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricUpdatePathResponse_ToResponse(t *testing.T) {
	r := MetricUpdatePathResponse{}
	expected := []byte("Metric updated successfully")
	actual := r.ToResponse()
	assert.Equal(t, expected, actual)
}
