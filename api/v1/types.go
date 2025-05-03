package v1

import (
	"github.com/ericogr/nav-pos/internal/app"
)

const HEADER_SESSION_ID = "X-Session-ID"

type HandleRequests struct {
	Telemetry  app.Telemetry
	Aircraft   app.Aircraft
	MapService app.TileService
	SessionId  int
}
