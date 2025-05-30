package model

type SensorData struct {
	TimestampUTC    string  `json:"timestamp_utc"`
	DeviceID        string  `json:"device_id"`
	IP              string  `json:"ip"`
	FirmwareVersion string  `json:"firmware_version"`
	ADCResolution   int     `json:"adc_resolution"`
	BatteryVoltage  float64 `json:"battery_voltage"`
	SoilMoisture    int     `json:"soil_moisture"`
}

type MoistureCategory string

const (
	MoistureDry         MoistureCategory = "Dry"
	MoistureLow         MoistureCategory = "Low Moisture"
	MoistureOptimal     MoistureCategory = "Optimal"
	MoistureMoist       MoistureCategory = "Moist"
	MoistureWet         MoistureCategory = "Wet"
	MoistureSensorError MoistureCategory = "SensorError"
)

// EnrichedSensorData adds computed fields
type EnrichedSensorData struct {
	TimestampUTC     string           `json:"timestamp_utc"`
	DeviceID         string           `json:"device_id"`
	IP               string           `json:"ip"`
	FirmwareVersion  string           `json:"firmware_version"`
	ADCResolution    int              `json:"adc_resolution"`
	BatteryVoltage   float64          `json:"battery_voltage"`
	SoilMoisture     int              `json:"soil_moisture"`
	MoisturePercent  int              `json:"moisture_percent"`
	MoistureCategory MoistureCategory `json:"moisture_category"`
}
