package repository

import (
	"context"
	"fmt"
	"github.com/Tabernol/krasiot-sensor/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoRepository(tableName string) (*DynamoRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	return &DynamoRepository{
		client:    client,
		tableName: tableName,
	}, nil
}

func (r *DynamoRepository) SaveSensorData(data model.SensorData) error {
	item := map[string]types.AttributeValue{
		"device_id":        &types.AttributeValueMemberS{Value: data.DeviceID},
		"timestamp_utc":    &types.AttributeValueMemberS{Value: data.TimestampUTC},
		"ip":               &types.AttributeValueMemberS{Value: data.IP},
		"firmware_version": &types.AttributeValueMemberS{Value: data.FirmwareVersion},
		"adc_resolution":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", data.ADCResolution)},
		"battery_voltage":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", data.BatteryVoltage)},
		"soil_moisture":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", data.SoilMoisture)},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName), // or hardcode table name for testing
		Item:      item,
	}

	_, err := r.client.PutItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}
	return nil
}
