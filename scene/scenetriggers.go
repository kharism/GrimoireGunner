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
	TriggerToPostLv1Story
	TriggerToPostLv2Story
	TriggerToPostLv3Story
)
