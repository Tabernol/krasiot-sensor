package main

import (
	"fmt"
	"github.com/Tabernol/krasiot-sensor/mqtt_broker"
	"log"
)

func main() {
	fmt.Println("Starting Krasiot Sensor Subscriber...")

	cfg, err := mqtt_broker.LoadMqttConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load MQTT config: %v", err)
	}
	fmt.Println(cfg)

	subscriber := mqtt_broker.NewMqttSubscriberService(cfg)
	go subscriber.ConnectAndSubscribe()

	//router := mux.NewRouter()
	//moistureHandler := handler.NewMoistureHandler(subscriber)
	//router.HandleFunc("/krasiot/api/v1/sensors/moisture/latest", moistureHandler.GetLatestMoisture).Methods("GET")
	//
	//log.Println("HTTP server starting on :8080")
	//if err := http.ListenAndServe(":8080", router); err != nil {
	//	log.Fatalf("Failed to start HTTP server: %v", err)
	//}

}
