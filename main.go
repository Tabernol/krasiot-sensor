package main

import (
	"fmt"
	"github.com/Tabernol/krasiot-sensor/handler"
	"github.com/Tabernol/krasiot-sensor/mqtt_broker"
	"github.com/Tabernol/krasiot-sensor/oracledb"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting Krasiot Sensor Subscriber...")

	cfg, err := mqtt_broker.LoadMqttConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load MQTT config: %v", err)
	}
	fmt.Println(cfg)

	oracleCfg, err := oracledb.LoadOracleConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load Oracle DB config: %v", err)
	}

	db, err := oracledb.InitOracle(oracleCfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Oracle DB: %v", err)
	}

	fmt.Println("CONNECTED to ADB")
	defer db.Close()

	//repo := oracledb.NewOracleRepository(db)
	//subscriber := mqtt_broker.NewMqttSubscriberService(cfg, repo)
	//go subscriber.ConnectAndSubscribe()

	subscriber := mqtt_broker.NewMqttSubscriberService(cfg)
	go subscriber.ConnectAndSubscribe()

	router := mux.NewRouter()
	moistureHandler := handler.NewMoistureHandler(subscriber)
	router.HandleFunc("/krasiot/api/v1/sensors/moisture/latest", moistureHandler.GetLatestMoisture).Methods("GET")

	log.Println("HTTP server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

}
