package app

import (
	"fmt"
	"io"
	"net/http"
)

type TileServiceOpenStreetMap struct {
}

func (m *TileServiceOpenStreetMap) GetTile(x, y, z int) ([]byte, error) {
	// OpenStreetMap tiles are served over HTTP, so we can use a simple HTTP GET request
	// to fetch the tile image. The URL format is:
	// https://tile.openstreetmap.org/{z}/{x}/{y}.png
	url := fmt.Sprintf("https://tile.openstreetmap.org/%d/%d/%d.png", z, x, y)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "navpos/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch tile: %s", resp.Status)
	}

	tile, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return tile, nil
}
