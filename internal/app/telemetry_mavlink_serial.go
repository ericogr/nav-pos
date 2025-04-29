package app

import (
	"context"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

const MAVLINK_SYSTEM_ID = 67
const TELEMETRY_MESSAGE_OK = "OK"

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
				switch e := evt.(type) {
				case *gomavlib.EventFrame:
					t.processEvent(e)
				case *gomavlib.EventChannelClose:
					t.processEventClose(e)
				}

				frm, ok := evt.(*gomavlib.EventFrame)

				if !ok || frm == nil {
					continue
				}

				t.processEvent(frm)
			}
		}
	}
}

func (t *TelemetryMavLinkSerial) processEventClose(evt *gomavlib.EventChannelClose) {
	t.TelemetryData.Valid = false
	t.TelemetryData.Message = evt.Error.Error()
}

func (t *TelemetryMavLinkSerial) processEvent(evt *gomavlib.EventFrame) {
	switch msg := evt.Message().(type) {
	case *ardupilotmega.MessageGlobalPositionInt:
		t.TelemetryData.Latitude = float64(msg.Lat) / 1e7
		t.TelemetryData.Longitude = float64(msg.Lon) / 1e7
		t.TelemetryData.Altitude = float64(msg.Alt) / 1000.0
		t.TelemetryData.Heading = float64(msg.Hdg) / 100.0
	case *ardupilotmega.MessageVfrHud:
		t.TelemetryData.GroundSpeed = float64(msg.Groundspeed) / 100 * 3.6
		t.TelemetryData.Altitude = float64(msg.Alt)
	case *ardupilotmega.MessageSysStatus:
		t.TelemetryData.BatteryVoltage = float64(msg.VoltageBattery) / 1000.0
		t.TelemetryData.BatteryCurrent = float64(msg.CurrentBattery) / 100.0
		t.TelemetryData.BatteryRemaining = uint8(msg.BatteryRemaining)
	case *ardupilotmega.MessageGpsRawInt:
		t.TelemetryData.NumSatellites = msg.SatellitesVisible
		t.TelemetryData.GroundSpeed = float64(msg.Vel) / 100 * 3.6
	}

	t.TelemetryData.Valid = true
	t.TelemetryData.Message = TELEMETRY_MESSAGE_OK
}

func NewTelemetryMavLink(devicePath string, bauldRate int) Telemetry {
	return &TelemetryMavLinkSerial{
		DevicePath: devicePath,
		BaudRate:   bauldRate,
	}
}
