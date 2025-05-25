package repository

//
//import (
//	"database/sql"
//	"fmt"
//	"github.com/Tabernol/krasiot-sensor/model"
//	"github.com/godror/godror"
//)
//
//type OracleRepository struct {
//	db *sql.DB
//}
//
//func NewOracleRepository(username, password, connectString string) (*OracleRepository, error) {
//	dsn := fmt.Sprintf(`user="%s" password="%s" connectString="%s"`, username, password, connectString)
//
//	db, err := sql.Open("godror", dsn)
//	if err != nil {
//		return nil, fmt.Errorf("failed to connect to Oracle DB: %w", err)
//	}
//
//	// Optional: ping to verify connection
//	if err := db.Ping(); err != nil {
//		return nil, fmt.Errorf("ping failed: %w", err)
//	}
//
//	return &OracleRepository{db: db}, nil
//}
//
//func (r *OracleRepository) SaveSensorData(data model.EnrichedSensorData) error {
//	query := `
//		INSERT INTO sensor_raw (
//			timestamp_utc,
//			device_id,
//			ip,
//			firmware_version,
//			adc_resolution,
//			battery_voltage,
//			soil_moisture,
//			moisture_category,
//			moisture_percent
//		)
//		VALUES (
//			TO_TIMESTAMP_TZ(:1, 'YYYY-MM-DD"T"HH24:MI:SS.FFTZH:TZM'),
//			:2, :3, :4, :5, :6, :7, :8, :9
//		)
//	`
//
//	_, err := r.db.Exec(
//		query,
//		data.TimestampUTC,
//		data.DeviceID,
//		data.IP,
//		data.FirmwareVersion,
//		data.ADCResolution,
//		data.BatteryVoltage,
//		data.SoilMoisture,
//		data.MoistureCategory,
//		data.MoisturePercent,
//	)
//
//	if err != nil {
//		return fmt.Errorf("failed to insert sensor data: %w", err)
//	}
//	return nil
//}
