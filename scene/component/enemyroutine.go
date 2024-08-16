package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// this abstract function will determine the behaviour of our NPC
// although it has Enemy, we can use this to determine the bahaviour
// of other characters
type EnemyRoutineFunc func(ecs *ecs.ECS, self *donburi.Entry)

type EnemyRoutineData struct {
	Routine EnemyRoutineFunc
	Memory  map[string]any
}

var EnemyRoutine = donburi.NewComponentType[EnemyRoutineData]()
