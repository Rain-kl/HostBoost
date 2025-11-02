package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Server.Port != ":15920" {
		t.Errorf("expected port :15920, got %s", cfg.Server.Port)
	}

	if cfg.Data.HostFile != "hosts.json" {
		t.Errorf("expected host file hosts.json, got %s", cfg.Data.HostFile)
	}

	if len(cfg.CORS.AllowOrigins) == 0 {
		t.Error("expected CORS allow origins to be set")
	}
}

func TestNormalizePort(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		expected string
	}{
		{"with colon", ":8080", ":8080"},
		{"without colon", "8080", ":8080"},
		{"empty", "", ":15920"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &ServerConfig{Port: tt.port}
			result := cfg.NormalizePort()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetMaxAge(t *testing.T) {
	tests := []struct {
		name     string
		maxAge   string
		expected time.Duration
	}{
		{"hours", "12h", 12 * time.Hour},
		{"minutes", "30m", 30 * time.Minute},
		{"invalid", "invalid", 12 * time.Hour}, // 默认值
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &CORSConfig{MaxAge: tt.maxAge}
			result := cfg.GetMaxAge()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLoadAndSave(t *testing.T) {
	// 创建临时配置文件
	tmpFile := "test_config.yaml"
	defer os.Remove(tmpFile)

	// 测试自动创建配置文件
	cfg, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Server.Port != ":15920" {
		t.Errorf("expected default port, got %s", cfg.Server.Port)
	}

	// 验证文件是否已创建
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("config file was not created")
	}

	// 测试加载已存在的配置文件
	cfg2, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("failed to reload config: %v", err)
	}

	if cfg2.Server.Port != cfg.Server.Port {
		t.Error("reloaded config does not match")
	}
}
