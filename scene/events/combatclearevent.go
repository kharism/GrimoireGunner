package events

import "github.com/yohamta/donburi/features/events"

type CombatClearData struct{}

var CombatClearEvent = events.NewEventType[CombatClearData]()
