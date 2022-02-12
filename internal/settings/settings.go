package settings

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type R53u2Settings struct {
	AWS           AWSSettings `yaml:"aws"`
	IPProvider    string      `yaml:"ip-provider"`
	CheckInterval string      `yaml:"check-interval"`
	Domains       []string    `yaml:"domains"`
}

type AWSSettings struct {
	AWSAccessKeyId     string `yaml:"aws-access-key-id"`
	AWSAccessKeySecret string `yaml:"aws-access-key-secret"`
	AWSDefaultRegion   string `yaml:"aws-default-region"`
}

func InitSettings(logger *zap.Logger, f string) *R53u2Settings {
	var r53u2Settings R53u2Settings

	fileSettings, err := os.ReadFile(f)
	if err != nil {
		logger.Fatal("Error loading settings.yaml file")
	}

	err = yaml.Unmarshal(fileSettings, &r53u2Settings)
	if err != nil {
		logger.Fatal("YAML failed to unmarshal to SheeshSettings", zap.Error(err))
	}

	return &r53u2Settings
}
