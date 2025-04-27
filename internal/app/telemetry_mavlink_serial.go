package app

import (
	"context"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

const MAVLINK_SYSTEM_ID = 67

type TelemetryMavLinkSerial struct {
	node gomavlib.Node

	DevicePath    string
	BaudRate      int
	TelemetryData TelemetryData
}

func (t *TelemetryMavLinkSerial) Initialize(ctx context.Context) error {
	t.node = gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: t.DevicePath,
				Baud:   t.BaudRate,
			},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: MAVLINK_SYSTEM_ID,
	}

	err := t.node.Initialize()
	if err != nil {
		return err
	}

	go t.process(ctx)

	return nil
}

func (t *TelemetryMavLinkSerial) GetTelemetryData() TelemetryData {
	return t.TelemetryData
}

func (t *TelemetryMavLinkSerial) Close() error {
	t.node.Close()

	return nil
}

func (t *TelemetryMavLinkSerial) process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Process incoming messages
			for evt := range t.node.Events() {
				frm, ok := evt.(*gomavlib.EventFrame)

				if !ok || frm == nil {
					continue
				}

				t.processEvent(frm)
			}
		}
	}
}

func (t *TelemetryMavLinkSerial) processEvent(evt *gomavlib.EventFrame) {
	switch msg := evt.Message().(type) {
	case *ardupilotmega.MessageGlobalPositionInt:
		t.TelemetryData.Latitude = float64(msg.Lat) / 1e7
		t.TelemetryData.Longitude = float64(msg.Lon) / 1e7
		t.TelemetryData.Altitude = float64(msg.Alt) / 1000.0 // t para metros
		t.TelemetryData.Heading = float64(msg.Hdg) / 100.0   // centigraus para graus
	case *ardupilotmega.MessageVfrHud:
		t.TelemetryData.GroundSpeed = float64(msg.Groundspeed) * 3.6 // m/s para km/h
		t.TelemetryData.Altitude = float64(msg.Alt)                  // metros
	case *ardupilotmega.MessageSysStatus:
		t.TelemetryData.BatteryVoltage = float64(msg.VoltageBattery) / 1000.0 // Converter de mV para V
		t.TelemetryData.BatteryCurrent = float64(msg.CurrentBattery) / 100.0  // Converter de cA para A
		t.TelemetryData.BatteryRemaining = uint8(msg.BatteryRemaining)        // Porcentagem de bateria restante
	case *ardupilotmega.MessageGpsRawInt:
		t.TelemetryData.NumSatellites = msg.SatellitesVisible
		t.TelemetryData.GPSValid = msg.FixType >= 2 // 2 = 2D fix, 3 = 3D fix
	}
}

func NewTelemetryMavLink(devicePath string, bauldRate int) Telemetry {
	return &TelemetryMavLinkSerial{
		DevicePath: devicePath,
		BaudRate:   bauldRate,
	}
}
