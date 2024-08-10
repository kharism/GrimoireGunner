package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	myComponent "github.com/kharism/mmbn_clone/scene/component"
)
type gridRenderer struct {
	query        *donburi.Query
	// orderedQuery *donburi.OrderedQuery[component.PositionData]
}
var GridRenderer = &gridRenderer{
	query:  donburi.NewQuery(
		filter.Contains(
			myComponent.GridPos,
		),
	),
}
func DrawGrid(ecs *ecs.ECS, screen *ebiten.Image) {
	// screen.Fill(color.RGBA{R: 41, G: 44, B: 45, A: 255})
	ecs.World.
}
