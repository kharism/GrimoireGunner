package asset

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// go:embed ../../tile_blue.png
var blueTilePng []byte

// go:embed ../../tile_blue.png
var redTilePng []byte

var BlueTile *ebiten.Image
var RedTile *ebiten.Image

func init() {
	if BlueTile == nil {
		imgReader := bytes.NewReader(blueTilePng)
		BlueTile, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if RedTile == nil {
		imgReader := bytes.NewReader(redTilePng)
		RedTile, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
}
