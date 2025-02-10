package scene

import "github.com/joelschutz/stagehand"

const (
	TriggerToReward stagehand.SceneTransitionTrigger = iota
	TriggerToCombat
	TriggerToInventory
	TriggerToStageSelect
	TriggerToRest
	TriggerToShop
	TriggerToMain
	TriggerToClear
	TriggerToOpening
)
