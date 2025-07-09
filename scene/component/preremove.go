package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type PreRemoveCallback func(ecs *ecs.ECS, entity *donburi.Entry)

var Preremove = donburi.NewComponentType[PreRemoveCallback]()
