package component

import (
	"time"

	"github.com/yohamta/donburi/ecs"
)

type eventQueue struct {
	Queue []Event
}

type Event interface {
	Execute(ecs *ecs.ECS)
	GetTime() time.Time
}

// eventqueue to generate/do things after certain period of time is reached
// do not mix this with donburi/feature/events.
// this does not process that events whatsoever
var EventQueue = eventQueue{Queue: []Event{}}

func (eq *eventQueue) AddEvent(ev Event) {
	j := []Event{}
	if len(eq.Queue) > 0 {
		for i := 0; i < len(eq.Queue); i++ {
			if eq.Queue[i].GetTime().After(ev.GetTime()) {
				j = append(j, ev)
				j = append(j, eq.Queue[i:len(eq.Queue)]...)
				break
			} else {
				j = append(j, eq.Queue[i])
			}
		}
		if j[len(j)-1].GetTime().Before(ev.GetTime()) || j[len(j)-1].GetTime() == ev.GetTime() {
			j = append(j, ev)
		}
	} else {
		j = append(j, ev)
	}

	eq.Queue = j
}
