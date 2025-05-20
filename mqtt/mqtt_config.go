package mqtt

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
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
	// Load from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	// Parse env variables
	port, err := strconv.Atoi(os.Getenv("MQTT_PORT"))
	if err != nil {
		return nil, err
	}

	return &MqttConfig{
		Host:     os.Getenv("MQTT_HOST"),
		Port:     port,
		Username: os.Getenv("MQTT_USERNAME"),
		Password: os.Getenv("MQTT_PASSWORD"),
		ClientID: os.Getenv("MQTT_CLIENT_ID"),
		Topic:    os.Getenv("MQTT_TOPIC"),
	}, nil
}
