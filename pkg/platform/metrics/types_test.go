package metrics_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"order-system/pkg/platform/metrics"
)

func TestMetricType(t *testing.T) {
	assert.Equal(t, metrics.MetricType(0), metrics.Counter)
	assert.Equal(t, metrics.MetricType(1), metrics.Gauge)
	assert.Equal(t, metrics.MetricType(2), metrics.Histogram)
}

func TestLabels(t *testing.T) {
	labels := metrics.Labels{
		"key1": "value1",
		"key2": "value2",
	}

	assert.Equal(t, "value1", labels["key1"])
	assert.Equal(t, "value2", labels["key2"])
	assert.Equal(t, 2, len(labels))

	// Test empty labels
	emptyLabels := metrics.Labels{}
	assert.Empty(t, emptyLabels)

	// Test nil labels
	var nilLabels metrics.Labels
	assert.Nil(t, nilLabels)
}

func TestMetric(t *testing.T) {
	now := time.Now()
	labels := metrics.Labels{"test": "value"}

	metric := metrics.Metric{
		Name:        "test_metric",
		Type:        metrics.Counter,
		Value:       42.0,
		Labels:      labels,
		Description: "Test metric",
		Timestamp:   now,
	}

	assert.Equal(t, "test_metric", metric.Name)
	assert.Equal(t, metrics.Counter, metric.Type)
	assert.Equal(t, 42.0, metric.Value)
	assert.Equal(t, labels, metric.Labels)
	assert.Equal(t, "Test metric", metric.Description)
	assert.Equal(t, now, metric.Timestamp)
}

func TestMetricWithEmptyLabels(t *testing.T) {
	metric := metrics.Metric{
		Name:        "test_metric",
		Type:        metrics.Gauge,
		Value:       1.0,
		Labels:      metrics.Labels{},
		Description: "Test metric",
		Timestamp:   time.Now(),
	}

	assert.NotNil(t, metric.Labels)
	assert.Empty(t, metric.Labels)
}

func TestMetricWithNilLabels(t *testing.T) {
	metric := metrics.Metric{
		Name:        "test_metric",
		Type:        metrics.Histogram,
		Value:       1.0,
		Description: "Test metric",
		Timestamp:   time.Now(),
	}

	assert.Nil(t, metric.Labels)
}

func TestMetricWithZeroValues(t *testing.T) {
	metric := metrics.Metric{}

	assert.Empty(t, metric.Name)
	assert.Equal(t, metrics.MetricType(0), metric.Type)
	assert.Equal(t, 0.0, metric.Value)
	assert.Nil(t, metric.Labels)
	assert.Empty(t, metric.Description)
	assert.True(t, metric.Timestamp.IsZero())
}

func TestMetricWithNegativeValue(t *testing.T) {
	metric := metrics.Metric{
		Name:  "negative_metric",
		Type:  metrics.Gauge,
		Value: -42.0,
	}

	assert.Equal(t, -42.0, metric.Value)
}

func TestMetricWithSpecialCharacters(t *testing.T) {
	metric := metrics.Metric{
		Name:        "test!@#$%^&*()_+",
		Description: "Special chars: !@#$%^&*()",
		Labels: metrics.Labels{
			"special!": "value@#",
		},
	}

	assert.Equal(t, "test!@#$%^&*()_+", metric.Name)
	assert.Equal(t, "Special chars: !@#$%^&*()", metric.Description)
	assert.Equal(t, "value@#", metric.Labels["special!"])
}
