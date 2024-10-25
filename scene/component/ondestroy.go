package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type OnDestroyFunc func(ecs *ecs.ECS, entry *donburi.Entry)

var OnDestroy = donburi.NewComponentType[OnDestroyFunc]()
