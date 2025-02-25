package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type PlayerMoveSystem struct {
	PlayerIndex *donburi.Entity
	isAnim      bool
}

func NewPlayerMoveSystem(player *donburi.Entity) *PlayerMoveSystem {
	return &PlayerMoveSystem{PlayerIndex: player}
}

// make this multiple of 5
var DefaultSpeed = 10.0
var QueryHP = donburi.NewQuery(
	filter.Contains(
		component.Health,
		component.GridPos,
	),
)

// check whether there are obstacle on row-col grid
// for now it checks another character
func ValidMove(ecs *ecs.ECS, row, col int) bool {
	ObstacleExist := false
	if row < 0 || row > 7 || col < 0 || col > 4 {
		return false
	}
	QueryHP.Each(ecs.World, func(e *donburi.Entry) {
		pos := component.GridPos.Get(e)
		if pos.Col == col && pos.Row == row {
			ObstacleExist = true
		}
	})
	return !ObstacleExist
}
func (p *PlayerMoveSystem) Update(ecs *ecs.ECS) {
	playerEntry := ecs.World.Entry(*p.PlayerIndex)
	// var targetX, targetY float64
	if !playerEntry.HasComponent(component.GridPos) {
		return
	}
	if playerEntry.HasComponent(component.Root) {
		return
	}
	if playerEntry == nil {
		return
	}
	gridPos := component.GridPos.Get(playerEntry)
	if !p.isAnim && inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {

		if gridPos.Row > 0 {
			// gridPos.Row -= 1
			if !ValidMove(ecs, gridPos.Row-1, gridPos.Col) {
				return
			}
			gridPos.Row -= 1
			// targetX, targetY = assets.GridCoord2Screen(gridPos.Row-1, gridPos.Col)
		}

		// component.GridPos.Set(playerEntry, gridPos)
	}
	if !p.isAnim && inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		gridPos := component.GridPos.Get(playerEntry)
		if gridPos.Row < 3 {
			// gridPos.Row += 1
			if !ValidMove(ecs, gridPos.Row+1, gridPos.Col) {
				return
			}
			gridPos.Row += 1
			// targetX, targetY = assets.GridCoord2Screen(gridPos.Row+1, gridPos.Col)
		}

		// component.GridPos.Set(playerEntry, gridPos)
	}
	if !p.isAnim && inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		gridPos := component.GridPos.Get(playerEntry)
		if gridPos.Col > 0 {
			// gridPos.Row += 1
			if !ValidMove(ecs, gridPos.Row, gridPos.Col-1) {
				return
			}
			gridPos.Col -= 1
			// targetX, targetY = assets.GridCoord2Screen(gridPos.Row, gridPos.Col-1)
		}

		if gridPos.Col <= -1 {
			gridPos.Col = 0
		}
		// component.GridPos.Set(playerEntry, gridPos)
	}
	if !p.isAnim && inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		gridPos := component.GridPos.Get(playerEntry)
		if gridPos.Col < 3 {
			// gridPos.Row += 1
			if !ValidMove(ecs, gridPos.Row, gridPos.Col+1) {
				return
			}
			gridPos.Col += 1
			// targetX, targetY = assets.GridCoord2Screen(gridPos.Row, gridPos.Col+1)
		}

		// component.GridPos.Set(playerEntry, gridPos)
	}
	// targetX, targetY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col+1)
	if !playerEntry.HasComponent(component.ScreenPos) {
		return
	}
	scrPos := component.ScreenPos.Get(playerEntry)
	scrPos.X = 0
	scrPos.Y = 0
	// screenPos := component.ScreenPos.Get(playerEntry)
	// if targetX != 0 && targetY != 0 {
	// 	p.isAnim = true
	// 	component.TargetLocation.Set(playerEntry, &component.MoveTargetData{Tx: targetX, Ty: targetY})
	// 	vx := targetX - screenPos.X
	// 	vy := targetY - screenPos.Y
	// 	if vx != 0 || vy != 0 {
	// 		speedVector := csg.NewVector(vx, vy, 0)
	// 		speedVector = speedVector.Normalize().MultiplyScalar(DefaultSpeed)
	// 		component.Speed.Set(playerEntry, &component.SpeedData{Vx: speedVector.X, Vy: speedVector.Y})
	// 	}
	// 	component.TargetLocation.Set(playerEntry, &component.MoveTargetData{Tx: targetX, Ty: targetY})
	// }
	// targetLoc := component.TargetLocation.Get(playerEntry)
	// speedComponent := component.Speed.Get(playerEntry)
	// screenPos.X += speedComponent.Vx
	// screenPos.Y += speedComponent.Vy
	// component.ScreenPos.Set(playerEntry, screenPos)
	// // fmt.Println(screenPos)
	// if screenPos.X == targetLoc.Tx && screenPos.Y == targetLoc.Ty {
	// 	i := component.Speed.Get(playerEntry)
	// 	i.Vx = 0
	// 	i.Vy = 0
	// 	p.isAnim = false
	// 	component.Speed.Set(playerEntry, i)
	// 	if assets.TileHeight > 0 {
	// 		col, row := assets.Coord2Grid(targetLoc.Tx, targetLoc.Ty)
	// 		comp := component.GridPos.Get(playerEntry)
	// 		comp.Col = col
	// 		comp.Row = row
	// 	}

	// }

}
