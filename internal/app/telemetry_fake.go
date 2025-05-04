package app

import "context"

type TelemetryFake struct {
}

func (t TelemetryFake) GetName() string {
	return "fake"
}

func (t TelemetryFake) GetTelemetryData() TelemetryData {
	return TelemetryData{
		Latitude:         -23.5431,
		Longitude:        -46.7889,
		Altitude:         1000,
		Heading:          90,
		GroundSpeed:      50,
		NumSatellites:    5,
		BatteryVoltage:   12.5,
		BatteryCurrent:   1.5,
		BatteryConsumed:  1000,
		BatteryRemaining: 80,
		Valid:            true,
		Message:          TELEMETRY_MESSAGE_OK,
	}
}

func (t TelemetryFake) Initialize(ctx context.Context) error {
	return nil
}

func (t TelemetryFake) Close() error {
	return nil
}

func NewFakeTelemetry() TelemetryService {
	return &TelemetryFake{}
}
