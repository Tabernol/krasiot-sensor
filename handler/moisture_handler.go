package handler

import (
	"github.com/Tabernol/krasiot-sensor/mqtt_broker"
	"net/http"
)

type MoistureHandler struct {
	subscriber *mqtt_broker.MqttSubscriberService
}

func NewMoistureHandler(sub *mqtt_broker.MqttSubscriberService) *MoistureHandler {
	return &MoistureHandler{subscriber: sub}
}

func (h *MoistureHandler) GetLatestMoisture(w http.ResponseWriter, r *http.Request) {
	msg := h.subscriber.GetLatestMessage()
	if msg == nil {
		http.Error(w, "No moisture data available!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(msg)
}
