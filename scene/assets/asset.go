package assets

import (
	"bytes"
	_ "embed"
	_ "image/jpeg"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed images/tile_blue.png
var blueTilePng []byte

//go:embed images/tile_red.png
var redTilePng []byte

//go:embed images/basicsprite.png
var player1Stand []byte

//go:embed images/attacksprite.png
var player1attack []byte

//go:embed images/magibullet.png
var projectile1 []byte

//go:embed images/boulder.png
var boulder []byte

//go:embed fonts/PixelOperator8-bold.ttf
var PixelFontTTF []byte

//go:embed images/bg_forest/bg.png
var bg_forest []byte

var BlueTile *ebiten.Image
var RedTile *ebiten.Image
var Bg *ebiten.Image
var BgForrest *ebiten.Image
var Player1Stand *ebiten.Image
var Player1Attack *ebiten.Image
var Projectile1 *ebiten.Image
var Boulder *ebiten.Image
var PixelFont *text.GoTextFaceSource
var FontFace *text.GoTextFace

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(PixelFontTTF))
	if err != nil {
		log.Fatal(err)
	}
	PixelFont = s
	FontFace = &text.GoTextFace{
		Source: PixelFont,
		Size:   15,
	}
	if BlueTile == nil {
		imgReader := bytes.NewReader(blueTilePng)
		BlueTile, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if RedTile == nil {
		imgReader := bytes.NewReader(redTilePng)
		RedTile, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}

	if Player1Stand == nil {
		imgReader := bytes.NewReader(player1Stand)
		Player1Stand, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Player1Attack == nil {
		imgReader := bytes.NewReader(player1attack)
		Player1Attack, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Projectile1 == nil {
		imgReader := bytes.NewReader(projectile1)
		Projectile1, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Boulder == nil {
		imgReader := bytes.NewReader(boulder)
		Boulder, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BgForrest == nil {
		imgReader := bytes.NewReader(bg_forest)
		BgForrest, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
}
