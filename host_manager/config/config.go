package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 包含应用程序的所有配置
type Config struct {
	Server ServerConfig `yaml:"server"`
	Data   DataConfig   `yaml:"data"`
	CORS   CORSConfig   `yaml:"cors"`
}

// ServerConfig 服务器相关配置
type ServerConfig struct {
	Port string `yaml:"port"`
}

// DataConfig 数据存储相关配置
type DataConfig struct {
	HostFile string `yaml:"host_file"`
	OptFile  string `yaml:"opt_file"`
}

// CORSConfig CORS 相关配置
type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           string   `yaml:"max_age"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: ":15920",
		},
		Data: DataConfig{
			HostFile: "data/hosts.json",
			OptFile:  "data/opts.json",
		},
		CORS: CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: false,
			MaxAge:           "12h",
		},
	}
}

// Load 从文件加载配置，如果文件不存在则创建默认配置文件
func Load(path string) (*Config, error) {
	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 文件不存在，创建默认配置
		cfg := DefaultConfig()
		if err := cfg.Save(path); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return cfg, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// Save 保存配置到文件
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetMaxAge 解析并返回 MaxAge 的 time.Duration
func (c *CORSConfig) GetMaxAge() time.Duration {
	duration, err := time.ParseDuration(c.MaxAge)
	if err != nil {
		return 12 * time.Hour // 默认值
	}
	return duration
}

// NormalizePort 确保端口号格式正确（以 : 开头）
func (s *ServerConfig) NormalizePort() string {
	if s.Port == "" {
		return ":15920"
	}
	if s.Port[0] != ':' {
		return ":" + s.Port
	}
	return s.Port
}
