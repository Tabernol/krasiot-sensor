package main

import (
	"context"
	"fmt"
	"github.com/Tabernol/krasiot-sensor/aws_sqs"
	"github.com/Tabernol/krasiot-sensor/handler"
	"github.com/Tabernol/krasiot-sensor/mqtt_broker"
	"github.com/Tabernol/krasiot-sensor/oracledb"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting Krasiot Sensor Subscriber...")

	// check location of oracle instant lib
	libDir := os.Getenv("ADB_LIB_DIR")
	fmt.Printf("Lib dir location is %s \n", libDir)

	// load and create configuration for MQTT
	cfg, err := mqtt_broker.LoadMqttConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load MQTT config: %v", err)
	}
	fmt.Printf("Port from config %d \n", cfg.Port)

	// load and create configuration for oracle ADB
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
	repo := oracledb.NewSensorRepository(db)

	// load configuration for SQS
	queueURL := os.Getenv("AWS_SQS_S_N_URL")
	if queueURL == "" {
		log.Fatal("❌ Environment variable AWS_SQS_S_N_URL is not set")
	}
	sqsNotifier := aws_sqs.NewSqsNotifier(context.Background(), queueURL)

	// create subscriber with all components
	subscriber := mqtt_broker.NewMqttSubscriberService(cfg, repo, sqsNotifier)
	go subscriber.ConnectAndSubscribe()

	// temporary endpoint
	router := mux.NewRouter()
	moistureHandler := handler.NewMoistureHandler(subscriber)
	router.HandleFunc("/krasiot/api/v1/sensors/moisture/latest", moistureHandler.GetLatestMoisture).Methods("GET")

	log.Println("HTTP server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

}
