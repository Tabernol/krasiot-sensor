package sensor_service

import "github.com/Tabernol/krasiot-sensor/model"

func EnrichSensorData(raw model.SensorData) model.EnrichedSensorData {
	percent, category := classifySoilMoisture(raw.SoilMoisture)

	return model.EnrichedSensorData{
		MeasuredAtUTC:    raw.MeasuredAtUTC,
		HardwareUID:      raw.HardwareUID,
		IP:               raw.IP,
		FirmwareVersion:  raw.FirmwareVersion,
		ADCResolution:    raw.ADCResolution,
		BatteryVoltage:   raw.BatteryVoltage,
		SoilMoisture:     raw.SoilMoisture,
		MoisturePercent:  percent,
		MoistureCategory: category,
	}
}

func classifySoilMoisture(adc int) (int, model.MoistureCategory) {
	if adc < 0 {
		return -1, model.MoistureSensorError
	}

	percent := toPercent(adc)

	switch {
	case percent >= 80:
		return percent, model.MoistureWet
	case percent >= 60:
		return percent, model.MoistureMoist
	case percent >= 40:
		return percent, model.MoistureOptimal
	case percent >= 20:
		return percent, model.MoistureLow
	default:
		return percent, model.MoistureDry
	}
}

const (
	adcWet = 5300 // ADC value in water
	adcDry = 8000 // ADC value in dry air
)

func toPercent(adc int) int {
	if adc <= adcWet {
		return 100
	}
	if adc >= adcDry {
		return 0
	}

	// Linear mapping from [adcWet, adcDry] to [100, 0]
	p := 100 - int(float64(adc-adcWet)/float64(adcDry-adcWet)*100)

	if p < 0 {
		return 0
	}
	if p > 100 {
		return 100
	}
	return p
}
