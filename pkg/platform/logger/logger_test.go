package logger

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			log, err := New(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, log)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, log)
			}
		})
	}
}

func TestLogger_LogLevels(t *testing.T) {
	var buf bytes.Buffer
	l := &defaultLogger{
		out:   &buf,
		level: Debug,
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "test-trace")
	ctx = context.WithValue(ctx, "span_id", "test-span")

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
				l.Error(ctx, "error message", errors.New("test error"))
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
			assert.Contains(t, output, "test-trace")
			assert.Contains(t, output, "test-span")
			if tt.hasError {
				assert.Contains(t, output, "test error")
			}
		})
	}
}

func TestLogger_WithComponent(t *testing.T) {
	var buf bytes.Buffer
	l := &defaultLogger{
		out:   &buf,
		level: Info,
	}

	componentLogger := l.WithComponent("test-component")
	componentLogger.Info(context.Background(), "component test")

	output := buf.String()
	assert.Contains(t, output, "test-component")
	assert.Contains(t, output, "component test")
}

func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	l := &defaultLogger{
		out:   &buf,
		level: Info,
	}

	fields := []Field{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: 123},
	}

	fieldLogger := l.WithFields(fields...)
	fieldLogger.Info(context.Background(), "field test")

	output := buf.String()
	assert.Contains(t, output, "key1")
	assert.Contains(t, output, "value1")
	assert.Contains(t, output, "key2")
	assert.Contains(t, output, "123")
}

func TestLogger_FileOutput(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "log-test-*.log")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	cfg := &config.Config{
		Logger: struct {
			Level      string `json:"level"`
			Format     string `json:"format"`
			Output     string `json:"output"`
			TimeFormat string `json:"timeFormat"`
		}{
			Level:  "info",
			Output: tmpFile.Name(),
		},
	}

	log, err := New(cfg)
	require.NoError(t, err)

	log.Info(context.Background(), "file test")

	content, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	assert.Contains(t, string(content), "file test")
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

func TestFieldsToMap(t *testing.T) {
	fields := []Field{
		{Key: "string", Value: "value"},
		{Key: "int", Value: 123},
		{Key: "bool", Value: true},
	}

	result := fieldsToMap(fields)

	assert.Equal(t, "value", result["string"])
	assert.Equal(t, 123, result["int"])
	assert.Equal(t, true, result["bool"])
}

func TestErrorToString(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"nil error", nil, ""},
		{"with error", errors.New("test error"), "test error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errorToString(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
