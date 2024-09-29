package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/grimoiregunner/scene/assets"
)

type CombatSceneNextStage struct {
	decorator CombatSceneDecorator
}

func (c *CombatSceneNextStage) DecorSceneData(data *SceneData) {
	data.SceneDecor = c.decorator
}

func (c *CombatSceneNextStage) GetNextStageTrigger() stagehand.SceneTransitionTrigger {
	return TriggerToCombat
}
func (c *CombatSceneNextStage) GetIcon() *ebiten.Image {
	return assets.BattleIcon
}

func NewCombatNextStage(decorator CombatSceneDecorator) *CombatSceneNextStage {
	if decorator != nil {
		return &CombatSceneNextStage{decorator: decorator}
	}
	return &CombatSceneNextStage{decorator: RandCombatDecorator()}
}

type RestSceneNextStage struct {
}

func (c *RestSceneNextStage) DecorSceneData(data *SceneData) {

}

func (c *RestSceneNextStage) GetNextStageTrigger() stagehand.SceneTransitionTrigger {
	return TriggerToRest
}
func (c *RestSceneNextStage) GetIcon() *ebiten.Image {
	return assets.RestIcon
}
