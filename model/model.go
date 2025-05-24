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
