package system

import (
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
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
				mycomponent.TargetLocation,
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
			screenPos.X = assets.TileStartX + float64(gridPos.Col)*float64(assets.TileWidth)
			screenPos.Y = assets.TileStartY + float64(gridPos.Row)*float64(assets.TileHeight)
		}
		screenPos.X += v.Vx
		screenPos.Y += v.Vy
		targetPos := mycomponent.TargetLocation.Get(e)
		if screenPos.X == targetPos.Tx && screenPos.Y == targetPos.Ty {
			v.Vx = 0
			v.Vy = 0
		}
		mycomponent.ScreenPos.Set(e, screenPos)
		col, row := assets.Coord2Grid(screenPos.X, screenPos.Y)
		comp := component.GridPos.Get(e)
		comp.Col = col
		comp.Row = row
		if col < 0 || col > 7 || row < 0 || row > 3 {
			ecs.World.Remove(e.Entity())
		}
	})
}