package mqtt_broker

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Tabernol/krasiot-sensor/aws_sqs"
	"github.com/Tabernol/krasiot-sensor/model"
	"github.com/Tabernol/krasiot-sensor/oracledb"
	sensor_service "github.com/Tabernol/krasiot-sensor/service"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"sync"
	"time"
)

type MqttSubscriberService struct {
	client        mqtt.Client
	cfg           *MqttConfig
	repo          *oracledb.SensorRepository
	notifier      *aws_sqs.SqsNotifier
	mu            sync.RWMutex
	latestMessage []byte
}

func NewMqttSubscriberService(cfg *MqttConfig, repo *oracledb.SensorRepository, notifier *aws_sqs.SqsNotifier) *MqttSubscriberService {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tls://%s:%d", cfg.Host, cfg.Port)).
		SetClientID(cfg.ClientID).
		SetUsername(cfg.Username).
		SetPassword(cfg.Password).
		SetTLSConfig(tlsConfig).
		SetAutoReconnect(true).
		SetCleanSession(false)

	client := mqtt.NewClient(opts)

	return &MqttSubscriberService{
		client:   client,
		cfg:      cfg,
		repo:     repo,
		notifier: notifier,
	}
}

func (s *MqttSubscriberService) ConnectAndSubscribe() {
	if token := s.client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Failed to connect to MQTT broker: %v", token.Error())
	}
	log.Println("✅ MQTT client connected")

	// Register message handler
	if token := s.client.Subscribe(s.cfg.Topic, 1, s.handleMessage); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Subscription failed: %v", token.Error())
	}
	log.Printf("✅ Subscribed to topic: %s", s.cfg.Topic)

	// Keep the subscriber running
	for {
		time.Sleep(1 * time.Second)
	}
}

func (s *MqttSubscriberService) handleMessage(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	log.Printf("📥 Received message: %s -> %s", msg.Topic(), string(payload))

	s.mu.Lock()
	s.latestMessage = make([]byte, len(payload))
	copy(s.latestMessage, payload)
	s.mu.Unlock()

	var rawData model.SensorData
	err := json.Unmarshal(payload, &rawData)
	if err != nil {
		log.Printf("❌ Failed to unmarshal message: %v", err)
		return
	}

	enriched := sensor_service.EnrichSensorData(rawData)

	if s.repo != nil {
		if err := s.repo.InsertSensorData(enriched); err != nil {
			log.Printf("❌ DB insert failed: %v", err)
		}
	}

	if s.notifier != nil {
		err := s.notifier.SendEnrichedSensorData(context.Background(), enriched)
		if err != nil {
			log.Printf("❌ Failed to send to SQS: %v", err)
		} else {
			log.Printf("📤 Sent message to SQS queue")
		}
	}
}

// temporary method
func (s *MqttSubscriberService) GetLatestMessage() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.latestMessage == nil {
		return nil
	}
	// Return a copy for safety
	msgCopy := make([]byte, len(s.latestMessage))
	copy(msgCopy, s.latestMessage)
	return msgCopy
}
