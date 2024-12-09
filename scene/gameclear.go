package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/grimoiregunner/scene/assets"
)

type Clear struct {
	data        *SceneData
	sm          *stagehand.SceneDirector[*SceneData]
	LevelLayout *Level
	musicPlayer *assets.AudioPlayer
	loopMusic   bool
}

func (r *Clear) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.BGClear, nil)
}

func (r *Clear) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		r.sm.ProcessTrigger(TriggerToMain)
	}
	return nil
}

var GameClearInstance = &Clear{}

func (r *Clear) Load(state *SceneData, manager stagehand.SceneController[*SceneData]) {
	r.sm = manager.(*stagehand.SceneDirector[*SceneData]) // This type assertion is important
	r.data = state
}

func (s *Clear) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *Clear) Unload() *SceneData {
	return s.data
}
