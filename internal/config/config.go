// Package config provides configuration loading and management for the application
// Package config provides functionality to load application configuration
package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig   `mapstructure:"app"`
	Kafka KafkaConfig `mapstructure:"kafka"`
}

type AppConfig struct {
	SampleDownloadURI string `mapstructure:"sample_download_uri"`
	ExportDownloadURI string `mapstructure:"export_download_uri"`
	ScheduleAPI       string `mapstructure:"schedule_api"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

func InitConfig() *Config {
	viper.SetConfigName("extra_config") // Name of your file (config.yaml)
	viper.SetConfigType("yml")
	viper.AddConfigPath(".") // Look in the current directory (project folder)

	// Enable environment variable overrides
	// Example: export APP_PORT=9000 will override the YAML
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil
	}

	return &config
}
