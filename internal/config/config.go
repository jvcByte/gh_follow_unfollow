package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	GitHubToken    string
	GitHubUsername string
	WorkerCount    int
	QueueSize      int
	TimeDelay      int
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetEnvPrefix("")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("worker_count", 1)
	v.SetDefault("queue_size", 3)
	v.SetDefault("time_delay_ms", 2000)

	if _, err := os.Stat(".env"); err == nil {
		v.SetConfigFile(".env")
		v.SetConfigType("env")
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading .env: %w", err)
		}
	}

	cfg := &Config{
		GitHubToken:    v.GetString("GH_TOKEN"),
		GitHubUsername: v.GetString("GH_USERNAME"),
		WorkerCount:    getIntWithDefault(v, "WORKER_COUNT", 1),
		QueueSize:      getIntWithDefault(v, "QUEUE_SIZE", 3),
		TimeDelay:      getIntWithDefault(v, "TIME_DELAY_MS", 2000),
	}

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.GitHubToken == "" {
		return errors.New("github.token is required")
	}
	if cfg.GitHubUsername == "" {
		return errors.New("github.username is required")
	}
	if cfg.WorkerCount <= 0 {
		return errors.New("app.worker_count must be > 0")
	}
	if cfg.QueueSize <= 0 {
		return errors.New("app.queue_size must be > 0")
	}
	return nil
}

func getIntWithDefault(v *viper.Viper, key string, defaultVal int) int {
	if val := v.GetInt(key); val > 0 {
		return val
	}
	return defaultVal
}
