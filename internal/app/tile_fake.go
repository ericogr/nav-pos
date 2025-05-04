package app

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"github.com/aquilax/go-perlin"
)

type TileServiceFake struct {
}

const (
	tileSize    = 256
	perlinAlpha = 2.0
	perlinBeta  = 2.0
	perlinN     = 4
	scaleFactor = 99.0 // maior = mais variação
	seed        = 42
)

func (a *TileServiceFake) GetName() string {
	return "fake"
}

func (a *TileServiceFake) GetContentType() string {
	return "image/png"
}

func (a *TileServiceFake) GetEncoding() string {
	return ""
}

func (m *TileServiceFake) GetTile(x, y, z int) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	p := perlin.NewPerlin(perlinAlpha, perlinBeta, perlinN, int64(seed))

	worldSize := 1 << z

	for px := 0; px < tileSize; px++ {
		for py := 0; py < tileSize; py++ {
			globalX := float64(x*tileSize + px)
			globalY := float64(y*tileSize + py)

			nx := globalX * scaleFactor / float64(worldSize*tileSize)
			ny := globalY * scaleFactor / float64(worldSize*tileSize)

			value := 0.0
			amp := 1.0
			freq := 1.0
			for o := 0; o < 5; o++ {
				value += p.Noise2D(nx*freq, ny*freq) * amp
				amp *= 0.5
				freq *= 2
			}

			normalized := (value + 1.5) / 3.0

			// Define cor por faixa de altitude
			var c color.RGBA
			switch {
			case normalized < 0.4:
				c = color.RGBA{30, 60, 200, 255} // water
			case normalized < 0.45:
				c = color.RGBA{220, 210, 160, 255} // sand
			case normalized < 0.65:
				c = color.RGBA{80, 180, 80, 255} // grass
			case normalized < 0.85:
				c = color.RGBA{100, 80, 50, 255} // mountain
			default:
				c = color.RGBA{255, 255, 255, 255} // snow
			}

			img.Set(px, py, c)
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
