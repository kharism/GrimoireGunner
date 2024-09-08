package system

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi/ecs"
)

type EventQueueSystem struct {
}

func (eq *EventQueueSystem) Update(ecs *ecs.ECS) {
	for {
		if len(component.EventQueue.Queue) == 0 {
			return
		}
		if component.EventQueue.Queue[0].GetTime().Before(time.Now()) {
			component.EventQueue.Queue[0].Execute(ecs)
			component.EventQueue.Queue = component.EventQueue.Queue[1:]
		} else {
			return
		}
	}
}
