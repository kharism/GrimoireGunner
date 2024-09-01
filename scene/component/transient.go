package component

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// this store how long an entity will stay on the field.
// once we reach over start+duration then the entity
// should be removed by system
type TransientData struct {
	Start    time.Time
	Duration time.Duration

	//this function is called before the entity is actually removed
	OnRemoveCallback func(ecs *ecs.ECS, entity *donburi.Entry)
}

var Transient = donburi.NewComponentType[TransientData]()
