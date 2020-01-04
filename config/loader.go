package config

import (
	"encoding/json"
	"os"
)

const envConfigKey = "ENV_APP_CONFIG"

func LoadConfig(config interface{}, defaultConfig string) error {
	rawConfig := envOrDefaultString(envConfigKey, defaultConfig)

	if err := json.Unmarshal([]byte(rawConfig), &config); err != nil {
		return err
	}
	return nil
}

func envOrDefaultString(envVar string, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}
