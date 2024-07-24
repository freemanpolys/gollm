// File: config.go

package goal

import (
	"time"

	"github.com/teilomillet/goal/internal/llm"
)

// LogLevel represents the level of logging
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type Config struct {
	Provider    string
	Model       string
	Temperature float64
	MaxTokens   int
	Timeout     time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
	APIKey      string
	DebugLevel  LogLevel
}

func (c *Config) toInternalConfig() *llm.Config {
	return &llm.Config{
		Provider:    c.Provider,
		Model:       c.Model,
		Temperature: c.Temperature,
		MaxTokens:   c.MaxTokens,
		Timeout:     c.Timeout,
		MaxRetries:  c.MaxRetries,
		RetryDelay:  c.RetryDelay,
		APIKeys:     map[string]string{c.Provider: c.APIKey},
		LogLevel:    llm.LogLevel(c.DebugLevel),
	}
}

// Convert llm.LogLevel to goal.LogLevel
func convertLogLevel(level llm.LogLevel) LogLevel {
	switch level {
	case llm.LogLevelDebug:
		return LogLevelDebug
	case llm.LogLevelInfo:
		return LogLevelInfo
	case llm.LogLevelWarn:
		return LogLevelWarn
	case llm.LogLevelError:
		return LogLevelError
	default:
		return LogLevelWarn // Default to Warn if unknown
	}
}

type ConfigOption func(*Config)

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	internalConfig, err := llm.LoadConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{
		Provider:    internalConfig.Provider,
		Model:       internalConfig.Model,
		Temperature: internalConfig.Temperature,
		MaxTokens:   internalConfig.MaxTokens,
		Timeout:     internalConfig.Timeout,
		MaxRetries:  internalConfig.MaxRetries,
		RetryDelay:  internalConfig.RetryDelay,
		DebugLevel:  convertLogLevel(internalConfig.LogLevel),
	}

	// Set the API key for the default provider
	if apiKey, exists := internalConfig.APIKeys[internalConfig.Provider]; exists {
		config.APIKey = apiKey
	}

	return config, nil
}

// SetProvider sets the provider in the Config
func SetProvider(provider string) ConfigOption {
	return func(c *Config) {
		c.Provider = provider
	}
}

// SetModel sets the model in the Config
func SetModel(model string) ConfigOption {
	return func(c *Config) {
		c.Model = model
	}
}

// SetTemperature sets the temperature in the Config
func SetTemperature(temperature float64) ConfigOption {
	return func(c *Config) {
		c.Temperature = temperature
	}
}

// SetMaxTokens sets the max tokens in the Config
func SetMaxTokens(maxTokens int) ConfigOption {
	return func(c *Config) {
		if maxTokens < 1 {
			maxTokens = 1
		}
		c.MaxTokens = maxTokens
	}
}

// SetTimeout sets the timeout in the Config
func SetTimeout(timeout time.Duration) ConfigOption {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// SetAPIKey sets the API key for the current provider in the Config
func SetAPIKey(apiKey string) ConfigOption {
	return func(c *Config) {
		c.APIKey = apiKey
	}
}

// SetMaxRetries sets the maximum number of retries in the Config
func SetMaxRetries(maxRetries int) ConfigOption {
	return func(c *Config) {
		c.MaxRetries = maxRetries
	}
}

// SetRetryDelay sets the delay between retries in the Config
func SetRetryDelay(retryDelay time.Duration) ConfigOption {
	return func(c *Config) {
		c.RetryDelay = retryDelay
	}
}

func SetDebugLevel(level LogLevel) ConfigOption {
	return func(c *Config) {
		c.DebugLevel = level
	}
}
