package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"order-system/pkg/infra/errors"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		message        string
		expectedCode   string
		expectedMsg    string
		expectedErrStr string
	}{
		{
			name:           "should create new error with code and message",
			code:           "ERR001",
			message:        "test error",
			expectedCode:   "ERR001",
			expectedMsg:    "test error",
			expectedErrStr: "ERR001: test error",
		},
		{
			name:           "should handle empty code and message",
			code:           "",
			message:        "",
			expectedCode:   "",
			expectedMsg:    "",
			expectedErrStr: ": ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.code, tt.message)
			assert.Equal(t, tt.expectedCode, err.Code)
			assert.Equal(t, tt.expectedMsg, err.Message)
			assert.Equal(t, tt.expectedErrStr, err.Error())
			assert.NotEmpty(t, err.Stack)
			assert.NotNil(t, err.Metadata)
		})
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name           string
		err           error
		code           string
		message        string
		expectedCode   string
		expectedMsg    string
		expectedErrStr string
		expectedNil    bool
	}{
		{
			name:           "should wrap existing error",
			err:           fmt.Errorf("original error"),
			code:          "ERR002",
			message:       "wrapped error",
			expectedCode:  "ERR002",
			expectedMsg:   "wrapped error",
			expectedErrStr: "ERR002: wrapped error: original error",
			expectedNil:   false,
		},
		{
			name:        "should return nil for nil error",
			err:         nil,
			code:        "ERR003",
			message:     "test message",
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.Wrap(tt.err, tt.code, tt.message)
			if tt.expectedNil {
				assert.Nil(t, err)
				return
			}
			assert.Equal(t, tt.expectedCode, err.Code)
			assert.Equal(t, tt.expectedMsg, err.Message)
			assert.Equal(t, tt.expectedErrStr, err.Error())
			assert.NotEmpty(t, err.Stack)
			assert.NotNil(t, err.Metadata)
		})
	}
}

func TestWithMetadata(t *testing.T) {
	err := errors.New("ERR004", "test error")
	err.WithMetadata("key1", "value1")
	err.WithMetadata("key2", 123)

	assert.Equal(t, "value1", err.Metadata["key1"])
	assert.Equal(t, 123, err.Metadata["key2"])
}

func TestGetStackTrace(t *testing.T) {
	err := errors.New("ERR005", "test error")
	assert.NotEmpty(t, err.Stack)
	assert.Contains(t, err.Stack, ".go:")
}
