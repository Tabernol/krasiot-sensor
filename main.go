package main

import (
	"fmt"
	"github.com/Tabernol/krasiot-sensor/handler"
	"github.com/Tabernol/krasiot-sensor/mqtt_broker"
	"github.com/Tabernol/krasiot-sensor/repository"
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
	//fmt.Println(cfg)

	adbCfg, err := repository.LoadDBConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load ADB config: %v", err)
	}
	fmt.Println("==================")
	fmt.Println(adbCfg.Username)
	fmt.Println("==================")

	//oracleRepo, err := repository.NewOracleRepository("ADMIN", "Ironbike=3862", "https://G34BA1A39372B52-KRASIOT.adb.us-phoenix-1.oraclecloudapps.com/ords/apex")
	//if err != nil {
	//	log.Fatalf("❌ Failed to init Oracle repository: %v", err)
	//}
	//mqtt_broker.SetOracleRepository(oracleRepo)

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
