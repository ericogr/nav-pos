package app

import "context"

const OPENSKY_NETWORK_API_STATES_ALL = "https://opensky-network.org/api/states/all"
const HEADER_X_RATE_LIMIT_REMAINING = "X-Rate-Limit-Remaining"

type Telemetry interface {
	GetTelemetryData() TelemetryData
	Initialize(ctx context.Context) error
	Close() error
}

type TelemetryData struct {
	Latitude         float64 `json:"lat"`
	Longitude        float64 `json:"lon"`
	Altitude         float64 `json:"alt"`
	Heading          float64 `json:"hdg"`
	GroundSpeed      float64 `json:"spd"`
	NumSatellites    uint8   `json:"sat"`
	BatteryVoltage   float64 `json:"bat"`
	BatteryCurrent   float64 `json:"cur"`
	BatteryConsumed  uint16  `json:"mah"`
	BatteryRemaining uint8   `json:"rem"`
	GPSValid         bool    `json:"gps"`
}

type Aircraft interface {
	GetAircrafts(bbox BoundingBox) ([]AircraftData, error)
}

type BoundingBox struct {
	MinLatitude  float64 `json:"min_latitude"`
	MinLongitude float64 `json:"min_longitude"`
	MaxLatitude  float64 `json:"max_latitude"`
	MaxLongitude float64 `json:"max_longitude"`
}

type AircraftData struct {
	Icao24         string  `json:"icao24"`
	Callsign       string  `json:"callsign"`
	OriginCountry  string  `json:"origin_country"`
	TimePosition   int64   `json:"time_position"`
	LastContact    int64   `json:"last_contact"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	BaroAltitude   float64 `json:"baro_altitude"`
	OnGround       bool    `json:"on_ground"`
	Velocity       float64 `json:"velocity"`
	TrueTrack      float64 `json:"true_track"`
	VerticalRate   float64 `json:"vertical_rate"`
	Sensors        []int   `json:"sensors"`
	GeoAltitude    float64 `json:"geo_altitude"`
	Squawk         string  `json:"squawk"`
	Spi            bool    `json:"spi"`
	PositionSource int     `json:"position_source"`
	Category       int     `json:"category"`
}

type OpenSkyResponse struct {
	Time   int64           `json:"time"`
	States [][]interface{} `json:"states"`
}
