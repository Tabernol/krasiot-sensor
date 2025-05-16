package mqtt

import (
	"crypto/tls"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"sync"
	"time"
)

type MqttSubscriberService struct {
	client        mqtt.Client
	cfg           *MqttConfig
	mu            sync.RWMutex
	latestMessage []byte
}

func NewMqttSubscriberService(cfg *MqttConfig) *MqttSubscriberService {
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
		client: client,
		cfg:    cfg,
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

	//var dto models.MoistureSensorDTO
	//err := json.Unmarshal(payload, &dto)
	//if err != nil {
	//	log.Printf("❌ Failed to unmarshal message: %v", err)
	//	return
	//}
	//
	////err = services.SaveMoisture(dto)
	//if err != nil {
	//	log.Printf("❌ Failed to save moisture data: %v", err)
	//} else {
	//	log.Println("✅ Moisture value saved to DB")
	//}
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
