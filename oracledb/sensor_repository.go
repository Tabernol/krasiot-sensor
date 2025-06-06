package oracledb

import (
	"database/sql"
	"fmt"
	"github.com/Tabernol/krasiot-sensor/model"
	"log"
	"time"
)

type SensorRepository struct {
	db *sql.DB
}

func NewSensorRepository(db *sql.DB) *SensorRepository {
	return &SensorRepository{db: db}
}

func (r *SensorRepository) InsertSensorData(data model.EnrichedSensorData) error {
	parsedTime, err := time.Parse(time.RFC3339, data.MeasuredAtUTC)
	if err != nil {
		return fmt.Errorf("failed to parse MeasuredAtUTC: %w", err)
	}

	query := `
		INSERT INTO sensor_raw (
			measured_at_utc,
			hardware_uid,
			ip,
			firmware_version,
			adc_resolution,
			battery_voltage,
			soil_moisture,
			moisture_category,
			moisture_percent
		) VALUES (
			:1, :2, :3, :4, :5, :6, :7, :8, :9
		)
	`

	_, err = r.db.Exec(query,
		parsedTime,
		data.HardwareUID,
		data.IP,
		data.FirmwareVersion,
		data.ADCResolution,
		data.BatteryVoltage,
		data.SoilMoisture,
		string(data.MoistureCategory),
		data.MoisturePercent,
	)

	if err != nil {
		return fmt.Errorf("DB insert failed: %w", err)
	}

	log.Println("✅ Sensor data inserted into Oracle DB")
	return nil
}
