package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/mayswind/lab/pkg/settings"
)

func TestClone(t *testing.T) {
	expectedObject := &settings.Config{
		AppName: "lab",
		Mode:    settings.MODE_PRODUCTION,
		DatabaseConfig: &settings.DatabaseConfig{
			DatabaseType:          settings.MySqlDbType,
			DatabaseHost:          "localhost",
			ConnectionMaxLifeTime: 60,
		},
		LogModes: []string{"console", "file"},
	}
	actualObject := &settings.Config{}
	err := Clone(expectedObject, actualObject)

	assert.Equal(t, nil, err)
	assert.EqualValues(t, expectedObject, actualObject)
}
