package mqtt_broker

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
	if err := godotenv.Load("../secrets/.env"); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	portStr := os.Getenv("MQTT_PORT")
	if portStr == "" {
		portStr = "8883" // default fallback
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid MQTT_PORT: %v", err)
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
