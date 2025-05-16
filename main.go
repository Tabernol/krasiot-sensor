package main

import (
	"fmt"
	"github.com/Tabernol/krasiot-sensor/mqtt"
	"log"
)

func main() {
	fmt.Println("Starting Krasiot Sensor Subscriber...")
	cfg, err := mqtt.LoadMqttConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	fmt.Println(cfg)

	subscriber := mqtt.NewMqttSubscriberService(cfg)
	subscriber.ConnectAndSubscribe()

}
