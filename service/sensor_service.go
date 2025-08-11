package sensor_service

import "github.com/Tabernol/krasiot-sensor/model"

func EnrichSensorData(raw model.SensorData) model.EnrichedSensorData {
	percent, category := classifySoilMoisture(raw.ADCResolution, raw.SoilMoisture)

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

func classifySoilMoisture(adc int, raw int) (int, model.MoistureCategory) {
	if adc < 0 {
		return -1, model.MoistureSensorError
	}

	percent := toPercent(adc, raw)

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
	adcWet13 = 5300 // ADC 13 value in water
	adcDry13 = 8000 // ADC 13 value in dry air
	adcWet12 = 1600 // ADC 12 value in water
	adcDry12 = 3300 // ADC 12 value in dry air
)

func toPercent(adc int, raw int) int {
	if adc == 12 {
		return calc12(raw)
	}
	if adc == 13 {
		return calc13(raw)
	}
	return -1
}

func calc13(raw int) int {
	if raw <= adcWet13 {
		return 100
	}
	if raw >= adcDry13 {
		return 0
	}

	// Linear mapping from [adcWet, adcDry] to [100, 0]
	p := 100 - int(float64(raw-adcWet13)/float64(adcDry13-adcWet13)*100)

	if p < 0 {
		return 0
	}
	if p > 100 {
		return 100
	}
	return p
}

func calc12(raw int) int {
	if raw <= adcWet12 {
		return 100
	}
	if raw >= adcDry12 {
		return 0
	}

	// Linear mapping from [adcWet, adcDry] to [100, 0]
	p := 100 - int(float64(raw-adcWet12)/float64(adcDry12-adcWet12)*100)

	if p < 0 {
		return 0
	}
	if p > 100 {
		return 100
	}
	return p
}
