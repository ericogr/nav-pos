package app

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type TileServiceMbtiles struct {
	databaseLocation string
}

func (a *TileServiceMbtiles) GetName() string {
	return "mbtiles"
}

func (a *TileServiceMbtiles) GetContentType() string {
	return "application/x-protobuf"
}

func (a *TileServiceMbtiles) GetEncoding() string {
	return "gzip"
}

func (m *TileServiceMbtiles) GetTile(x, y, z int) ([]byte, error) {
	db, err := sql.Open("sqlite3", m.databaseLocation)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// MBTiles uses TMS (y inverted)
	tileRow := (1 << z) - 1 - y

	query := `SELECT tile_data FROM tiles WHERE zoom_level = ? AND tile_column = ? AND tile_row = ? LIMIT 1`
	var data []byte
	err = db.QueryRow(query, z, x, tileRow).Scan(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
