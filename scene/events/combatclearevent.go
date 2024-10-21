package events

import "github.com/yohamta/donburi/features/events"

type CombatClearData struct {
	IsGameOver bool
}

var CombatClearEvent = events.NewEventType[CombatClearData]()
