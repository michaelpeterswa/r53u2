package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type R53u2Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"error"`

	CheckIPProvider string        `env:"CHECK_IP_PROVIDER" envDefault:"https://checkip.amazonaws.com"`
	CheckIPTimeout  time.Duration `env:"CHECK_IP_TIMEOUT" envDefault:"5s"`
	Domains         []string      `env:"DOMAINS" envSeparator:","`

	AWSAccessKeyId     string `env:"AWS_ACCESS_KEY_ID"`
	AWSAccessKeySecret string `env:"AWS_ACCESS_KEY_SECRET"`
	AWSDefaultRegion   string `env:"AWS_DEFAULT_REGION"`

	CronSchedule string `env:"CRON_SCHEDULE" envDefault:"@every 30m"`

	MetricsEnabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	MetricsPort    int  `env:"METRICS_PORT" envDefault:"8081"`

	TracingEnabled    bool    `env:"TRACING_ENABLED" envDefault:"false"`
	TracingSampleRate float64 `env:"TRACING_SAMPLERATE" envDefault:"0.01"`
	TracingService    string  `env:"TRACING_SERVICE" envDefault:"r53u2"`
	TracingVersion    string  `env:"TRACING_VERSION"`
}

func NewConfig() (*R53u2Config, error) {
	var cfg R53u2Config

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
