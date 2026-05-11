package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// 设置环境变量
	os.Setenv("TEMPORAL_ADDRESS", "test-temporal:7233")
	os.Setenv("REDIS_ADDRESS", "test-redis:6379")
	os.Setenv("GITHUB_TOKEN", "test-token")
	os.Setenv("ANTHROPIC_API_KEY", "test-key")

	// 加载配置
	cfg := Load()

	// 验证配置
	if cfg.Temporal.Address != "test-temporal:7233" {
		t.Errorf("Expected Temporal address 'test-temporal:7233', got '%s'", cfg.Temporal.Address)
	}

	if cfg.Redis.Address != "test-redis:6379" {
		t.Errorf("Expected Redis address 'test-redis:6379', got '%s'", cfg.Redis.Address)
	}

	if cfg.GitHub.Token != "test-token" {
		t.Errorf("Expected GitHub token 'test-token', got '%s'", cfg.GitHub.Token)
	}

	if cfg.Anthropic.APIKey != "test-key" {
		t.Errorf("Expected Anthropic API key 'test-key', got '%s'", cfg.Anthropic.APIKey)
	}

	// 清理环境变量
	os.Unsetenv("TEMPORAL_ADDRESS")
	os.Unsetenv("REDIS_ADDRESS")
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("ANTHROPIC_API_KEY")
}

func TestLoadDefaults(t *testing.T) {
	// 清除所有环境变量
	os.Unsetenv("TEMPORAL_ADDRESS")
	os.Unsetenv("REDIS_ADDRESS")
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("ANTHROPIC_API_KEY")

	// 加载配置
	cfg := Load()

	// 验证默认值
	if cfg.Temporal.Address != "localhost:7233" {
		t.Errorf("Expected default Temporal address 'localhost:7233', got '%s'", cfg.Temporal.Address)
	}

	if cfg.Redis.Address != "localhost:6379" {
		t.Errorf("Expected default Redis address 'localhost:6379', got '%s'", cfg.Redis.Address)
	}

	if cfg.Temporal.TaskQueue != "develop-workflow" {
		t.Errorf("Expected default task queue 'develop-workflow', got '%s'", cfg.Temporal.TaskQueue)
	}
}
