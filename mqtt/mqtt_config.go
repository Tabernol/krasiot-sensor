package mqtt

import (
	"github.com/spf13/viper"
)

type MqttConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	ClientID string `mapstructure:"client-id"`
	Topic    string
}

func LoadMqttConfig() (*MqttConfig, error) {
	viper.SetConfigName("application") // application.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // current dir

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config MqttConfig
	err := viper.Sub("mqtt").Unmarshal(&config)
	return &config, err
}
