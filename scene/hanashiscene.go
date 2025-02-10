package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/hanashi/core"
)

type HanashiScene struct {
	scene    *core.Scene
	State    *SceneData
	director *stagehand.SceneDirector[*SceneData]
}

func (m *HanashiScene) Update() error {
	e := m.scene.Update()
	if e != nil {
		return e
	}

	return nil
}
func NewHanashiScene(hanashiScene *core.Scene) *HanashiScene {
	return &HanashiScene{scene: hanashiScene}
}
func (m *HanashiScene) Draw(screen *ebiten.Image) {
	m.scene.Draw(screen)
	// m.SkipButton.Draw(screen)
	// txt := "click to continue"
	// txtOpt := text.DrawOptions{}
	// txtOpt.ColorScale.ScaleWithColor(RED)
	// txtOpt.GeoM.Scale(0.5, 0.5)
	// text.Draw(screen, txt, face, &txtOpt)
}
func (m *HanashiScene) SetDoneFunc(done func()) {
	m.scene.Done = done
}
func (s *HanashiScene) Load(state *SceneData, director stagehand.SceneController[*SceneData]) {
	s.director = director.(*stagehand.SceneDirector[*SceneData]) // This type assertion is important
	s.scene.Events[0].Execute(s.scene)
	s.State = state

}
func (s *HanashiScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *HanashiScene) Unload() *SceneData {
	// your unload code
	// s.scene.Events[0].Execute(s.scene)
	s.scene.EventIndex = 0
	s.scene.CurCharName = ""
	s.scene.CurDialog = ""
	s.scene.ViewableCharacters = []*core.Character{}
	// s.scene.ViewableCharacters = []*core.Character{}
	s.scene.VisibleDialog = ""
	if s.State == nil {
		s.State = &SceneData{}
	}
	if s.State.HanashiChoices == nil {
		s.State.HanashiChoices = map[string]any{}
	}
	for key, val := range s.scene.SceneData {
		s.State.HanashiChoices[key] = val
	}

	return s.State
}
