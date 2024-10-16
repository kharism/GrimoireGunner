package attack

import (
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func JoinOnAtkHit(f1, f2 component.OnAtkHit) component.OnAtkHit {
	return func(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
		f1(ecs, projectile, receiver)
		f2(ecs, projectile, receiver)
	}
}
