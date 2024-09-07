package system

import (
	"fmt"
	"time"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/yohamta/donburi/ecs"
)

type energySystem struct {
	CurrentEN int
	MaxEN     int
	ENRegen   int //en recovered per 0.1 sec

	lastRegen time.Time
}

// energy system is system that deals with player's energy (duh)
// the Current,Max and ENRegen is int by hundreds just to simplify calculation
// to display divide the number by one hundreds
var EnergySystem = &energySystem{CurrentEN: 300, MaxEN: 300, ENRegen: 20}

func (e *energySystem) Update(ecs *ecs.ECS) {
	if e.lastRegen.IsZero() {
		e.lastRegen = time.Now()
	}
	now := time.Now()
	dist := now.Sub(e.lastRegen)
	Enrecovery := int(dist.Milliseconds() * int64(e.ENRegen) / 100)
	// fmt.Println(int64(dist), float64(dist/time.Second))
	e.lastRegen = now
	e.CurrentEN = e.CurrentEN + Enrecovery
	if e.CurrentEN > e.MaxEN {
		e.CurrentEN = e.MaxEN
	}

}

var ENTextStartX = 200.0
var ENTextStartY = 50.0

var ENBarStartX = 250.0
var ENBarStartY = 50.0

func (e *energySystem) SetEn(val int) {
	e.CurrentEN = val

}

func (e *energySystem) GetEn() int {
	return e.CurrentEN
}
func (e *energySystem) GetMaxEn() int {
	return e.MaxEN
}

func (e *energySystem) DrawEnBar(ecs *ecs.ECS, screen *ebiten.Image) {
	textMes := fmt.Sprintf("%d/%d", e.CurrentEN/100, e.MaxEN/100)
	screenRect := screen.Bounds()
	blackBox := ebiten.NewImage(screenRect.Dx(), 80)
	blackBox.Fill(color.Black)
	screen.DrawImage(blackBox, &ebiten.DrawImageOptions{})
	// draw bar
	redBar := ebiten.NewImage(200, 10)
	redBar.Fill(color.RGBA{R: 255, A: 255})

	textTranslate := ebiten.GeoM{}
	textTranslate.Translate(ENTextStartX, ENTextStartY)
	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: textTranslate,
		},
	}
	text.Draw(screen, textMes, assets.FontFace, op)
	barTranslate := ebiten.GeoM{}
	barTranslate.Translate(ENBarStartX, ENBarStartY)
	screen.DrawImage(redBar, &ebiten.DrawImageOptions{GeoM: barTranslate})
	width := float64(e.CurrentEN) / float64(e.MaxEN) * 200.0
	if width > 0 {
		blueBar := ebiten.NewImage(int(width), 10)
		blueBar.Fill(color.CMYK{C: 255})
		screen.DrawImage(blueBar, &ebiten.DrawImageOptions{GeoM: barTranslate})
	}

}
