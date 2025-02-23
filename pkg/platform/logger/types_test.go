package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected string
	}{
		{
			name:     "Debug level",
			level:    Debug,
			expected: "DEBUG",
		},
		{
			name:     "Info level",
			level:    Info,
			expected: "INFO",
		},
		{
			name:     "Warn level",
			level:    Warn,
			expected: "WARN",
		},
		{
			name:     "Error level",
			level:    Error,
			expected: "ERROR",
		},
		{
			name:     "Unknown level",
			level:    Level(99),
			expected: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.String())
		})
	}
}
