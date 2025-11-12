package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Agent AgentConfig `yaml:"agent"`
	LLM   LLMConfig   `yaml:"llm"`
}

// AgentConfig contains agent-specific configuration
type AgentConfig struct {
	DecisionInterval time.Duration `yaml:"decision_interval"`
	MaxEventChain    int           `yaml:"max_event_chain"`
}

// LLMConfig contains LLM client configuration
type LLMConfig struct {
	APIEndpoint    string  `yaml:"api_endpoint"`
	APIKey         string  `yaml:"api_key"`
	Model          string  `yaml:"model"`
	MaxTokens      int     `yaml:"max_tokens"`
	Temperature    float64 `yaml:"temperature"`
	TimeoutSeconds int     `yaml:"timeout_seconds"`
}

// Load loads configuration from a YAML file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	cfg.applyDefaults()

	// Validate configuration
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// applyDefaults sets default values for missing configuration
func (c *Config) applyDefaults() {
	if c.Agent.DecisionInterval == 0 {
		c.Agent.DecisionInterval = 30 * time.Second
	}
	if c.Agent.MaxEventChain == 0 {
		c.Agent.MaxEventChain = 100
	}
	if c.LLM.MaxTokens == 0 {
		c.LLM.MaxTokens = 150
	}
	if c.LLM.Temperature == 0 {
		c.LLM.Temperature = 0.7
	}
	if c.LLM.TimeoutSeconds == 0 {
		c.LLM.TimeoutSeconds = 30
	}
	if c.LLM.Model == "" {
		c.LLM.Model = "llama3.2"
	}
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	if c.LLM.APIEndpoint == "" {
		return fmt.Errorf("llm.api_endpoint is required")
	}
	if c.Agent.DecisionInterval < time.Second {
		return fmt.Errorf("agent.decision_interval must be at least 1 second")
	}
	if c.LLM.MaxTokens < 10 {
		return fmt.Errorf("llm.max_tokens must be at least 10")
	}
	if c.LLM.Temperature < 0 || c.LLM.Temperature > 2 {
		return fmt.Errorf("llm.temperature must be between 0 and 2")
	}
	return nil
}

// Example returns an example configuration
func Example() *Config {
	return &Config{
		Agent: AgentConfig{
			DecisionInterval: 30 * time.Second,
			MaxEventChain:    100,
		},
		LLM: LLMConfig{
			APIEndpoint:    "http://localhost:11434/v1/completions",
			APIKey:         "",
			Model:          "llama3.2",
			MaxTokens:      150,
			Temperature:    0.7,
			TimeoutSeconds: 30,
		},
	}
}