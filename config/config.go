package config

import "os"

// TODO add config entries as the app grows
type Config struct {
	Port    string
	Token   string
	HookUrl string
}

// NewConfig creates a new Config struct
func NewConfig() *Config {
	return &Config{
		Port:    getEnv("TH_PORT", "8080"),
		Token:   getEnv("TH_TOKEN", "123"),
		HookUrl: getEnv("TH_HOOK_URL", "http://localhost:8080/hook"),
	}
}

// getEnv returns the value of an environment variable or a default value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
