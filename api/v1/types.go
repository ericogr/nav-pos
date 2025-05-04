package v1

import (
	"github.com/ericogr/nav-pos/internal/app"
)

const HEADER_SESSION_ID = "X-Session-ID"

type HandleRequests struct {
	TelemetryService app.TelemetryService
	RadarService     app.RadarService
	TileMapService   app.TileMapService
	SessionId        int
}
