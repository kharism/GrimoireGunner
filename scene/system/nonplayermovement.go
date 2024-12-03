package system

import (
	csg "github.com/kharism/golang-csg/core"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
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
				component.Speed,
				component.GridPos,
				component.ScreenPos,
				// mycomponent.Sprite,
				component.TargetLocation,
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
		v := component.Speed.Get(e)
		gridPos := component.GridPos.Get(e)
		screenPos := component.ScreenPos.Get(e)
		targetPos := component.TargetLocation.Get(e)
		if screenPos.X == 0 && screenPos.Y == 0 {
			screenPos.X = assets.TileStartX + float64(gridPos.Col)*float64(assets.TileWidth)
			screenPos.Y = assets.TileStartY + float64(gridPos.Row)*float64(assets.TileHeight)
		}
		if v.Vx == 0 && v.Vy == 0 && v.V != 0 && targetPos.Tx != 0 && targetPos.Ty != 0 {
			distX := targetPos.Tx - screenPos.X
			distY := targetPos.Ty - screenPos.Y
			distVector := csg.NewVector(distX, distY, 0)
			distVector = distVector.Normalize().MultiplyScalar(v.V)
			v.Vx = distVector.X
			v.Vy = distVector.Y
		}
		screenPos.X += v.Vx
		screenPos.Y += v.Vy

		if screenPos.X == targetPos.Tx && screenPos.Y == targetPos.Ty {
			v.Vx = 0
			v.Vy = 0
			v.V = 0
		}
		component.ScreenPos.Set(e, screenPos)
		col, row := assets.Coord2Grid(screenPos.X, screenPos.Y)
		comp := component.GridPos.Get(e)
		comp.Col = col
		comp.Row = row
		if col < 0 || col > 7 || row < 0 || row > 3 {
			ecs.World.Remove(e.Entity())
		}
	})
}
