package assets

import (
	"bytes"
	_ "embed"
	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed images/tile_blue.png
var blueTilePng []byte

//go:embed images/tile_red.png
var redTilePng []byte

//go:embed images/bg.jpg
var bg []byte

//go:embed images/basicsprite.png
var player1Stand []byte

//go:embed images/boulder.png
var boulder []byte

var BlueTile *ebiten.Image
var RedTile *ebiten.Image
var Bg *ebiten.Image
var Player1Stand *ebiten.Image
var Boulder *ebiten.Image

func init() {
	if BlueTile == nil {
		imgReader := bytes.NewReader(blueTilePng)
		BlueTile, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if RedTile == nil {
		imgReader := bytes.NewReader(redTilePng)
		RedTile, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Bg == nil {
		imgReader := bytes.NewReader(bg)
		Bg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Player1Stand == nil {
		imgReader := bytes.NewReader(player1Stand)
		Player1Stand, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Boulder == nil {
		imgReader := bytes.NewReader(boulder)
		Boulder, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
}
