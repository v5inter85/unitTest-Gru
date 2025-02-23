package metrics

import (
	"order-system/pkg/infra/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("metrics disabled", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Metrics.Enabled = false

		collector, err := New(cfg)
		assert.Error(t, err)
		assert.Nil(t, collector)
	})

	t.Run("metrics enabled", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Metrics.Enabled = true

		collector, err := New(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, collector)
	})
}

func TestRegister(t *testing.T) {
	cfg := &config.Config{}
	cfg.Metrics.Enabled = true
	collector, _ := New(cfg)

	t.Run("register new counter", func(t *testing.T) {
		err := collector.Register("test_counter", Counter, "test counter")
		assert.NoError(t, err)
	})

	t.Run("register new gauge", func(t *testing.T) {
		err := collector.Register("test_gauge", Gauge, "test gauge")
		assert.NoError(t, err)
	})

	t.Run("register new histogram", func(t *testing.T) {
		err := collector.Register("test_histogram", Histogram, "test histogram")
		assert.NoError(t, err)
	})

	t.Run("register duplicate metric", func(t *testing.T) {
		err := collector.Register("test_counter", Counter, "duplicate counter")
		assert.Error(t, err)
	})
}

func TestSetGauge(t *testing.T) {
	cfg := &config.Config{}
	cfg.Metrics.Enabled = true
	collector, _ := New(cfg)

	t.Run("set and get gauge", func(t *testing.T) {
		name := "test_gauge"
		labels := Labels{"label1": "value1"}

		err := collector.Register(name, Gauge, "test gauge")
		assert.NoError(t, err)

		collector.SetGauge(name, 10.0, labels)
		value := collector.GetGauge(name, labels)
		assert.Equal(t, 10.0, value)
	})

	t.Run("set unregistered gauge", func(t *testing.T) {
		name := "unregistered_gauge"
		labels := Labels{"label1": "value1"}

		collector.SetGauge(name, 10.0, labels)
		value := collector.GetGauge(name, labels)
		assert.Equal(t, 0.0, value)
	})
}

func TestObserveHistogram(t *testing.T) {
	cfg := &config.Config{}
	cfg.Metrics.Enabled = true
	collector, _ := New(cfg)

	t.Run("observe and get histogram", func(t *testing.T) {
		name := "test_histogram"
		labels := Labels{"label1": "value1"}

		err := collector.Register(name, Histogram, "test histogram")
		assert.NoError(t, err)

		collector.ObserveHistogram(name, 1.0, labels)
		collector.ObserveHistogram(name, 2.0, labels)

		values := collector.GetHistogram(name, labels)
		assert.Equal(t, []float64{1.0, 2.0}, values)
	})

	t.Run("observe unregistered histogram", func(t *testing.T) {
		name := "unregistered_histogram"
		labels := Labels{"label1": "value1"}

		collector.ObserveHistogram(name, 1.0, labels)
		values := collector.GetHistogram(name, labels)
		assert.Nil(t, values)
	})
}

func TestCollect(t *testing.T) {
	cfg := &config.Config{}
	cfg.Metrics.Enabled = true
	collector, _ := New(cfg)

	t.Run("collect all metric types", func(t *testing.T) {
		labels := Labels{"service": "test"}

		counterName := "test_counter"
		err := collector.Register(counterName, Counter, "test counter")
		assert.NoError(t, err)
		collector.IncrementCounter(counterName, 1.0, labels)

		gaugeName := "test_gauge"
		err = collector.Register(gaugeName, Gauge, "test gauge")
		assert.NoError(t, err)
		collector.SetGauge(gaugeName, 10.0, labels)

		histogramName := "test_histogram"
		err = collector.Register(histogramName, Histogram, "test histogram")
		assert.NoError(t, err)
		collector.ObserveHistogram(histogramName, 5.0, labels)

		metrics := collector.Collect()

		expectedLabels := labelsToString(labels)

		assert.Len(t, metrics, 3)
		for _, metric := range metrics {
			assert.Equal(t, metric.Labels, stringToLabels(expectedLabels))
			assert.NotEmpty(t, metric.Description)
			assert.NotZero(t, metric.Timestamp)

			switch metric.Name {
			case "test_counter":
				assert.Equal(t, Counter, metric.Type)
				assert.Equal(t, 1.0, metric.Value)
			case "test_gauge":
				assert.Equal(t, Gauge, metric.Type)
				assert.Equal(t, 10.0, metric.Value)
			case "test_histogram":
				assert.Equal(t, Histogram, metric.Type)
				assert.Equal(t, 5.0, metric.Value)
			}
		}
	})
}

func TestLabelsConversion(t *testing.T) {
	t.Run("empty labels", func(t *testing.T) {
		labels := Labels{}
		str := labelsToString(labels)
		assert.Empty(t, str)

		convertedLabels := stringToLabels(str)
		assert.Equal(t, labels, convertedLabels)
	})
}
