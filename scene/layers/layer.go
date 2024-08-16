package layers

import "github.com/yohamta/donburi/ecs"

const (
	LayerBackground ecs.LayerID = iota
	LayerGrid
	LayerCharacter
	LayerFx
	LayerHP
	LayerDebug
	LayerUI
)
