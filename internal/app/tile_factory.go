package app

import "fmt"

func CreateTileMapService(kind string, params map[string]string) (TileMapService, error) {
	switch kind {
	case "fake":
		return &TileServiceFake{}, nil
	case "openstreetmap":
		return &TileServiceOpenStreetMap{}, nil
	case "mbtiles":
		databaseLocation := params["databaseLocation"]
		if databaseLocation == "" {
			databaseLocation = "local.mbtiles"
		}

		return &TileServiceMbtiles{
			databaseLocation: databaseLocation,
		}, nil
	default:
		return nil, fmt.Errorf("unknown tile map service: %s", kind)
	}
}
