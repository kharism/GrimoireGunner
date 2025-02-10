package assets

import (
	"bytes"
	_ "embed"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func ReadAndFillBg(raw []byte) *ebiten.Image {
	reader := bytes.NewReader(raw)
	h, _, _ := ebitenutil.NewImageFromReader(reader)
	jj := ebiten.NewImage(64, 64)
	jj.Fill(color.RGBA{R: 0x4f, G: 0x8f, B: 0xba, A: 255})
	jj.DrawImage(h, nil)
	return jj
}

//go:embed images/portrait/sven.png
var sven []byte

var Sven *ebiten.Image

//go:embed images/portrait/sven_2.png
var sven2 []byte

var Sven2 *ebiten.Image

//go:embed images/portrait/jack.png
var jack []byte

var Jack *ebiten.Image

//go:embed images/portrait/shizuku.png
var shizuku []byte

var Shizuku *ebiten.Image

//go:embed images/bedroom.png
var bedroom []byte

var Bedroom *ebiten.Image

//go:embed images/bedroomdoor.png
var bedroomdoor []byte

var BedroomDoor *ebiten.Image

//go:embed images/workshop.png
var workshop_1 []byte

var Workshop1 *ebiten.Image

//go:embed images/workshop_2.png
var workshop_2 []byte

var Workshop2 *ebiten.Image

func init() {
	Sven = ReadAndFillBg(sven)
	Sven2 = ReadAndFillBg(sven2)
	Shizuku = ReadAndFillBg(shizuku)
	Jack = ReadAndFillBg(jack)

	if Bedroom == nil {
		imgReader := bytes.NewReader(bedroom)
		Bedroom, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}

	if BedroomDoor == nil {
		imgReader := bytes.NewReader(bedroomdoor)
		BedroomDoor, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Workshop1 == nil {
		imgReader := bytes.NewReader(workshop_1)
		Workshop1, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Workshop2 == nil {
		imgReader := bytes.NewReader(workshop_2)
		Workshop2, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
}
