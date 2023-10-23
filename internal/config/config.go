package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var AppName string = "franz"

type Configuration struct {
	KafkaBootstrapURLs     string `mapstructure:"KAFKA_BOOTSTRAP_URLS"`
	KafkaConsumerGroup     string `mapstructure:"KAFKA_CONSUMER_GROUP"`
	KafkaSchemaRegistryURL string `mapstructure:"KAFKA_SR_URL"`
	Debug                  bool   `mapstructure:"DEBUG"`
}

func NewConfiguration() (*Configuration, error) {
	// Initialize Viper
	viper.AutomaticEnv()
	viper.SetEnvPrefix(AppName)

	// Read configuration from environment variables or a config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Set default values for bound env
	viper.SetDefault("KAFKA_BOOTSTRAP_URLS", "localhost:29092")
	viper.SetDefault("KAFKA_CONSUMER_GROUP", AppName)
	viper.SetDefault("KAFKA_SR_URL", "http://localhost:8081")
	viper.SetDefault("DEBUG", false)

	// Create a new configuration instance
	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	fmt.Println(config)
	return &config, nil
}
