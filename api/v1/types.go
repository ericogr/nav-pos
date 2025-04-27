package v1

import (
	"github.com/ericogr/nav-pos/internal/app"
)

type HandleRequests struct {
	Telemetry app.Telemetry
	Aircraft  app.Aircraft
}
