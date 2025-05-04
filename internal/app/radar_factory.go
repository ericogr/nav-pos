package app

import "fmt"

func CreateRadarService(kind string, params map[string]string) (RadarService, error) {
	switch kind {
	case "fake":
		return &AircraftFake{}, nil
	case "opensky":
		return &AircraftOpenSky{
			Username: params["username"],
			Password: params["password"],
		}, nil
	default:
		return nil, fmt.Errorf("unknown radar service: %s", kind)
	}
}
