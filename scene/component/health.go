package component

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type HealthData struct {
	HP           int
	Name         string
	InvisTime    time.Time
	OnTakeDamage OnTakeDamageFunc
}
type DamageDetail struct {
	Amount int
}

// this function is called whenever anything with Health component takes damage.
// changing sprite, add IFrame, whatever
// DONOT take damage here. Implement it on OnAtkHit instead
// self is the entity who takes damage
type OnTakeDamageFunc func(ecs *ecs.ECS, self *donburi.Entry, detail DamageDetail)

func AddIFrame(ecs *ecs.ECS, self *donburi.Entry, detail DamageDetail) {
	Health.Get(self).InvisTime = time.Now().Add(1000 * time.Millisecond)
}

var Health = donburi.NewComponentType[HealthData]()
