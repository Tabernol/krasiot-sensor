package main

import (
	"fmt"
	"github.com/Tabernol/krasiot-sensor/handler"
	"github.com/Tabernol/krasiot-sensor/mqtt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting Krasiot Sensor Subscriber...")
	cfg, err := mqtt.LoadMqttConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	fmt.Println(cfg)

	subscriber := mqtt.NewMqttSubscriberService(cfg)
	go subscriber.ConnectAndSubscribe()

	router := mux.NewRouter()
	moistureHandler := handler.NewMoistureHandler(subscriber)
	router.HandleFunc("/krasiot/api/v1/sensors/moisture/latest", moistureHandler.GetLatestMoisture).Methods("GET")

	log.Println("HTTP server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

}
