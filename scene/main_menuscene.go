package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/grimoiregunner/scene/assets"
)

type MainMenuScene struct {
	sm           *stagehand.SceneDirector[*SceneData]
	data         *SceneData
	selectedMenu int
}

var menus = []string{
	"New Game",
	"Exit",
}

func (r *MainMenuScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		r.selectedMenu += 1
		if r.selectedMenu == len(menus) {
			r.selectedMenu -= 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		r.selectedMenu -= 1
		if r.selectedMenu == -1 {
			r.selectedMenu += 1
		}
	}

	return nil
}

var MainMenuInstance = &MainMenuScene{}

func (r *MainMenuScene) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.BgOpening, &ebiten.DrawImageOptions{})
	// buttonBg := ebiten.NewImage(248, 50)
	// buttonBg.Fill(color.RGBA{R: 0x72, G: 0x72, B: 0x72, A: 0xFF})
	textColor := ebiten.ColorScale{}
	textColor.Scale(0, 0, 0, 1)
	for idx, i := range menus {
		pos := ebiten.GeoM{}
		if idx == r.selectedMenu {
			pos.Scale(1.5, 1)
		}
		pos.Translate(0, float64(165+55*idx))

		screen.DrawImage(assets.MenuButtonBg, &ebiten.DrawImageOptions{GeoM: pos})
		pos.Reset()
		pos.Scale(1.6, 1.6)
		pos.Translate(50, float64(165+55*idx))

		text.Draw(screen, i, assets.UnispaceFace, &text.DrawOptions{

			DrawImageOptions: ebiten.DrawImageOptions{GeoM: pos, ColorScale: textColor},
		})
	}
}
func init() {
	assets.UnispaceFace = &text.GoTextFace{
		Source: assets.UnispaceFont,
		Size:   15,
	}
}
func (r *MainMenuScene) Load(state *SceneData, manager stagehand.SceneController[*SceneData]) {

}
func (s *MainMenuScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *MainMenuScene) Unload() *SceneData {
	return s.data
}
