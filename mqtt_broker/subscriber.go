package mqtt_broker

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
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
	mu            sync.RWMutex
	latestMessage []byte
}

func NewMqttSubscriberService(cfg *MqttConfig, repo *oracledb.SensorRepository) *MqttSubscriberService {
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
		repo:   repo,
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

	//// Encode the enriched data to JSON
	//jsonData, err := json.Marshal(enriched)
	//if err != nil {
	//	log.Printf("❌ Failed to marshal enriched data: %v", err)
	//	return
	//}

	if s.repo != nil {
		if err := s.repo.InsertSensorData(enriched); err != nil {
			log.Printf("❌ DB insert failed: %v", err)
		}
	}

	//// Send HTTP POST to ORDS endpoint
	//url := "https://g34ba1a39372b52-krasiot.adb.us-phoenix-1.oraclecloudapps.com/ords/admin/api/sensors"
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	//if err != nil {
	//	log.Printf("❌ Failed to create POST request: %v", err)
	//	return
	//}
	//req.Header.Set("Content-Type", "application/json")
	//
	//clientHTTP := &http.Client{Timeout: 10 * time.Second}
	//resp, err := clientHTTP.Do(req)
	//if err != nil {
	//	log.Printf("❌ HTTP POST request failed: %v", err)
	//	return
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode >= 200 && resp.StatusCode < 300 {
	//	log.Println("✅ Data posted successfully to ORDS")
	//} else {
	//	log.Printf("❌ Failed to post data to ORDS. Status: %s", resp.Status)
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
