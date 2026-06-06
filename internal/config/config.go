package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Providers []ProviderConfig `yaml:"providers"`
	Routing   RoutingConfig   `yaml:"routing"`
	Cost      CostConfig      `yaml:"cost"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ProviderConfig struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	APIKey  string `yaml:"api_key"`
	BaseURL string `yaml:"base_url,omitempty"`
	Models  []string `yaml:"models"`
}

type RoutingConfig struct {
	DefaultProvider string `yaml:"default_provider"`
	DefaultModel   string `yaml:"default_model"`
	Rules          []Rule `yaml:"rules"`
}

type Rule struct {
	Pattern  string `yaml:"pattern"`
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
}

type CostConfig struct {
	MaxBudgetPerMonth float64 `yaml:"max_budget_per_month"`
	AlertThreshold    float64 `yaml:"alert_threshold"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	data = expandEnv(data)

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func expandEnv(data []byte) []byte {
	s := string(data)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			s = strings.ReplaceAll(s, "${"+parts[0]+"}", parts[1])
		}
	}
	return []byte(s)
}

func (c *Config) validate() error {
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if len(c.Providers) == 0 {
		return fmt.Errorf("at least one provider is required")
	}
	if c.Routing.DefaultProvider == "" {
		c.Routing.DefaultProvider = c.Providers[0].Name
	}
	if c.Routing.DefaultModel == "" && len(c.Providers[0].Models) > 0 {
		c.Routing.DefaultModel = c.Providers[0].Models[0]
	}
	return nil
}

func (c *Config) FindProvider(name string) *ProviderConfig {
	for _, p := range c.Providers {
		if p.Name == name {
			return &p
		}
	}
	return nil
}
