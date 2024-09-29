package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/system/attack"
)

type RestScene struct {
	data          *SceneData
	sm            *stagehand.SceneDirector[*SceneData]
	selectedIndex int
}

var RestSceneInstance = &RestScene{}

func (r *RestScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) && r.selectedIndex == 1 {
		r.selectedIndex = (r.selectedIndex - 1) % 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) && r.selectedIndex == 0 {
		r.selectedIndex = (r.selectedIndex + 1) % 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if r.selectedIndex == 0 {
			r.data.PlayerHP += 300
			if r.data.PlayerHP > r.data.PlayerMaxHP {
				r.data.PlayerHP = r.data.PlayerMaxHP
			}
		} else {
			t := attack.NewHealCaster()
			r.data.Inventory = append(r.data.Inventory, t)
		}
		r.sm.ProcessTrigger(TriggerToCombat)
	}
	return nil
}

var restOptionStartX = 71.0
var restOptionStartY = 107.0

func (r *RestScene) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.BgRest, &ebiten.DrawImageOptions{})
	geom := ebiten.GeoM{}
	geom.Translate(196, 67)
	text.Draw(screen, "Pick one", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: geom,
		},
	})
	bg2 := ebiten.NewImage(900, 70)
	bg2.Fill(color.RGBA{R: 0x9b, G: 0x55, B: 0x22, A: 255})
	geom.Reset()
	rect := assets.CardPick.Bounds()
	geom.Scale(940/float64(rect.Dx()), 75/float64(rect.Dy()))

	geom.Translate(restOptionStartX-20, restOptionStartY+float64(r.selectedIndex)*90)
	screen.DrawImage(assets.CardPick, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
	geom.Reset()
	geom.Translate(restOptionStartX, restOptionStartY)
	screen.DrawImage(bg2, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
	geom.Translate(0, 10)
	text.Draw(screen, "Heal 300HP", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: geom,
		},
	})
	geom.Translate(0, 80)
	screen.DrawImage(bg2, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
	geom.Translate(0, 10)
	text.Draw(screen, "Get Healcaster", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: geom,
		},
	})
}
func (r *RestScene) Load(state *SceneData, manager stagehand.SceneController[*SceneData]) {
	r.data = state
	r.sm = manager.(*stagehand.SceneDirector[*SceneData]) // This type assertion is important
	r.selectedIndex = 0
}
func (s *RestScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *RestScene) Unload() *SceneData {
	return s.data
}
