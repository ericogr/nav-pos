package app

import (
	"fmt"
	"strconv"
)

func CreateTelemetryService(kind string, params map[string]string) (TelemetryService, error) {
	switch kind {
	case "fake":
		return &TelemetryFake{}, nil
	case "mavlinkserial":
		device := params["device"]
		if device == "" {
			device = "/dev/ttyUSB0"
		}
		baudRate := params["baudrate"]
		if baudRate == "" {
			baudRate = "115200"
		}
		// Convert baudRate to int
		baudRateInt, err := strconv.Atoi(baudRate)
		if err != nil {
			return nil, fmt.Errorf("invalid baud rate: %s", baudRate)
		}
		return &TelemetryMavLinkSerial{
			DevicePath: device,
			BaudRate:   baudRateInt,
		}, nil
	default:
		return nil, fmt.Errorf("unknown telemetry service: %s", kind)
	}
}
