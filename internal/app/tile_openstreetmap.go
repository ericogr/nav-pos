package app

import (
	"fmt"
	"io"
	"net/http"
)

type TileServiceOpenStreetMap struct {
}

func (a *TileServiceOpenStreetMap) GetName() string {
	return "openstreetmap"
}

func (a *TileServiceOpenStreetMap) GetContentType() string {
	return "image/png"
}

func (a *TileServiceOpenStreetMap) GetEncoding() string {
	return ""
}

func (m *TileServiceOpenStreetMap) GetTile(x, y, z int) ([]byte, error) {
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
