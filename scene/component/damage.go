package component

import (
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type DamageData struct {
	Damage    int
	OnHitAnim *core.AnimatedImage
}

// inflict damage here.
type OnAtkHit func(ecs *ecs.ECS, projectile, receiver *donburi.Entry)

var Damage = donburi.NewComponentType[DamageData]()
var OnHit = donburi.NewComponentType[OnAtkHit]()
