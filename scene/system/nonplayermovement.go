package system

import (
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/component"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type npMoveSystem struct {
	query *donburi.Query
}

var NPMoveSystem = &npMoveSystem{
	query: donburi.NewQuery(
		filter.And(
			filter.Contains(
				mycomponent.Speed,
				mycomponent.GridPos,
				mycomponent.ScreenPos,
				mycomponent.Sprite,
			),
			filter.Not(
				filter.Contains(
					archetype.PlayerTag,
				),
			),
		),
	),
}

func (s *npMoveSystem) Update(ecs *ecs.ECS) {
	s.query.Each(ecs.World, func(e *donburi.Entry) {
		v := mycomponent.Speed.Get(e)
		gridPos := mycomponent.GridPos.Get(e)
		screenPos := mycomponent.ScreenPos.Get(e)
		if screenPos.X == 0 && screenPos.Y == 0 {
			screenPos.X = TileStartX + float64(gridPos.Col)*float64(tileWidth)
			screenPos.Y = TileStartY + float64(gridPos.Row)*float64(tileHeight)
		}
		screenPos.X += v.Vx
		screenPos.Y += v.Vy
		mycomponent.ScreenPos.Set(e, screenPos)
		col, row := Coord2Grid(screenPos.X, screenPos.Y)
		comp := component.GridPos.Get(e)
		comp.Col = col
		comp.Row = row
		if col < 0 || col > 7 || row < 0 || row > 3 {
			ecs.World.Remove(e.Entity())
		}
	})
}
