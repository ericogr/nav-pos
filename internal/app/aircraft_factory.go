package app

import "fmt"

func CreateAircraft(kind string, params map[string]string) (Aircraft, error) {
	switch kind {
	case "fake":
		return &AircraftFake{}, nil
	case "opensky":
		return &AircraftOpenSky{
			Username: params["username"],
			Password: params["password"],
		}, nil
	default:
		return nil, fmt.Errorf("unknown aircraft provider: %s", kind)
	}
}
