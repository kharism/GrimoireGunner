package component

import (
	"time"

	"github.com/yohamta/donburi"
)

// this store how long an entity will stay on the field.
// once we reach over start+duration then the entity
// should be removed by system
type TransientData struct {
	Start    time.Time
	Duration time.Duration
}

var Transient = donburi.NewComponentType[TransientData]()
