package logger

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"order-system/pkg/infra/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name: "valid stdout config",
			cfg: &config.Config{
				Logger: struct {
					Level      string `json:"level"`
					Format     string `json:"format"`
					Output     string `json:"output"`
					TimeFormat string `json:"timeFormat"`
				}{
					Level:  "info",
					Output: "stdout",
				},
			},
			wantErr: false,
		},
		{
			name: "valid stderr config",
			cfg: &config.Config{
				Logger: struct {
					Level      string `json:"level"`
					Format     string `json:"format"`
					Output     string `json:"output"`
					TimeFormat string `json:"timeFormat"`
				}{
					Level:  "debug",
					Output: "stderr",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			cfg: &config.Config{
				Logger: struct {
					Level      string `json:"level"`
					Format     string `json:"format"`
					Output     string `json:"output"`
					TimeFormat string `json:"timeFormat"`
				}{
					Level:  "invalid",
					Output: "stdout",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := New(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, l)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, l)
			}
		})
	}
}

func TestLogger_Levels(t *testing.T) {
	buf := &bytes.Buffer{}
	l := &defaultLogger{
		out:   buf,
		level: Debug,
	}

	ctx := context.Background()
	testErr := errors.New("test error")

	tests := []struct {
		name     string
		logFunc  func()
		level    string
		message  string
		hasError bool
	}{
		{
			name: "debug log",
			logFunc: func() {
				l.Debug(ctx, "debug message")
			},
			level:    "DEBUG",
			message:  "debug message",
			hasError: false,
		},
		{
			name: "info log",
			logFunc: func() {
				l.Info(ctx, "info message")
			},
			level:    "INFO",
			message:  "info message",
			hasError: false,
		},
		{
			name: "warn log",
			logFunc: func() {
				l.Warn(ctx, "warn message")
			},
			level:    "WARN",
			message:  "warn message",
			hasError: false,
		},
		{
			name: "error log",
			logFunc: func() {
				l.Error(ctx, "error message", testErr)
			},
			level:    "ERROR",
			message:  "error message",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			output := buf.String()
			assert.Contains(t, output, tt.level)
			assert.Contains(t, output, tt.message)
			if tt.hasError {
				assert.Contains(t, output, testErr.Error())
			}
		})
	}
}

func TestLogger_WithComponent(t *testing.T) {
	buf := &bytes.Buffer{}
	l := &defaultLogger{
		out:   buf,
		level: Info,
	}

	componentName := "test-component"
	l = l.WithComponent(componentName).(*defaultLogger)

	ctx := context.Background()
	l.Info(ctx, "test message")

	output := buf.String()
	assert.Contains(t, output, componentName)
}

func TestLogger_WithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	l := &defaultLogger{
		out:   buf,
		level: Info,
	}

	fields := []Field{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: 123},
	}

	l = l.WithFields(fields...).(*defaultLogger)

	ctx := context.Background()
	l.Info(ctx, "test message")

	output := buf.String()
	for _, f := range fields {
		assert.Contains(t, output, f.Key)
	}
}

func TestLogger_TraceContext(t *testing.T) {
	buf := &bytes.Buffer{}
	l := &defaultLogger{
		out:   buf,
		level: Info,
	}

	traceID := "trace-123"
	spanID := "span-456"
	ctx := context.WithValue(context.Background(), "trace_id", traceID)
	ctx = context.WithValue(ctx, "span_id", spanID)

	l.Info(ctx, "test message")

	output := buf.String()
	assert.Contains(t, output, traceID)
	assert.Contains(t, output, spanID)
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		want    Level
		wantErr bool
	}{
		{"debug level", "debug", Debug, false},
		{"info level", "info", Info, false},
		{"warn level", "warn", Warn, false},
		{"error level", "error", Error, false},
		{"invalid level", "invalid", Info, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLevel(tt.level)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
