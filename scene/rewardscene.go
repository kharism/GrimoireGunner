package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/grimoiregunner/scene/assets"
)

type RewardScene struct {
	data SceneData
	sm   *stagehand.SceneManager[SceneData]
}

func (r *RewardScene) Update() error {
	return nil
}

var RewardSceneInstance = &RewardScene{}

func (r *RewardScene) Draw(screen *ebiten.Image) {
	bg := ebiten.NewImage(1024, 600)
	bg.Fill(color.RGBA{R: 0x21, G: 0x43, B: 0x58, A: 255})
	screen.DrawImage(bg, &ebiten.DrawImageOptions{})
	textTranslate := ebiten.GeoM{}
	textTranslate.Translate(512, 90)
	text.Draw(screen, "UPGRADE", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: textTranslate,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	})
}

func (r *RewardScene) Load(state SceneData, manager stagehand.SceneController[SceneData]) {
	r.sm = manager.(*stagehand.SceneManager[SceneData]) // This type assertion is important
	r.data = state
}
func (s *RewardScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *RewardScene) Unload() SceneData {
	// your unload code
	return s.data
}
