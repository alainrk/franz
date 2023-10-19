package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	KafkaBootstrapURLs string `mapstructure:"KAFKA_BOOTSTRAP_URLS"`
}

func NewConfiguration() (*Configuration, error) {
	// Initialize Viper
	viper.AutomaticEnv()
	viper.SetEnvPrefix("franz")

	// Read configuration from environment variables or a config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Set default values for bound env
	viper.SetDefault("KAFKA_BOOTSTRAP_URLS", "localhost:29092")

	// Create a new configuration instance
	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
