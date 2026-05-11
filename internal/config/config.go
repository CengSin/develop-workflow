package config

import (
	"os"
	"strconv"
)

type Config struct {
	Temporal  TemporalConfig
	Redis     RedisConfig
	GitHub    GitHubConfig
	Anthropic AnthropicConfig
	App       AppConfig
}

type TemporalConfig struct {
	Address   string
	Namespace string
	TaskQueue string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type GitHubConfig struct {
	Token   string
	Owner   string
	Private bool
}

type AnthropicConfig struct {
	APIKey string
	Model  string
}

type AppConfig struct {
	Env      string
	LogLevel string
}

func Load() *Config {
	return &Config{
		Temporal: TemporalConfig{
			Address:   getEnv("TEMPORAL_ADDRESS", "localhost:7233"),
			Namespace: getEnv("TEMPORAL_NAMESPACE", "default"),
			TaskQueue: getEnv("TEMPORAL_TASK_QUEUE", "develop-workflow"),
		},
		Redis: RedisConfig{
			Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		GitHub: GitHubConfig{
			Token:   getEnv("GITHUB_TOKEN", ""),
			Owner:   getEnv("GITHUB_OWNER", ""),
			Private: getEnvAsBool("GITHUB_PRIVATE", true),
		},
		Anthropic: AnthropicConfig{
			APIKey: getEnv("ANTHROPIC_API_KEY", ""),
			Model:  getEnv("ANTHROPIC_MODEL", "claude-sonnet-4-20250514"),
		},
		App: AppConfig{
			Env:      getEnv("APP_ENV", "development"),
			LogLevel: getEnv("LOG_LEVEL", "debug"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
