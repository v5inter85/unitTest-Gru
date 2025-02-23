package config_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"order-system/pkg/infra/config"
)

func TestConfigUnmarshalEmptyJSON(t *testing.T) {
	emptyJSON := `{}`

	var cfg config.Config
	err := json.Unmarshal([]byte(emptyJSON), &cfg)
	assert.NoError(t, err)
}
