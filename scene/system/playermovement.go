package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	csg "github.com/kharism/golang-csg/core"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type PlayerMoveSystem struct {
	PlayerIndex *donburi.Entity
	isAnim      bool
}

func NewPlayerMoveSystem(player *donburi.Entity) *PlayerMoveSystem {
	return &PlayerMoveSystem{PlayerIndex: player}
}

// make this multiple of 5
var DefaultSpeed = 5.0

func (p *PlayerMoveSystem) Update(ecs *ecs.ECS) {
	playerEntry := ecs.World.Entry(*p.PlayerIndex)
	var targetX, targetY float64
	if !p.isAnim && inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		gridPos := component.GridPos.Get(playerEntry)
		if gridPos.Row > 0 {
			// gridPos.Row -= 1
			targetX, targetY = GridCoord2Screen(gridPos.Row-1, gridPos.Col)
		}
		component.GridPos.Set(playerEntry, gridPos)
	}
	if !p.isAnim && inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		gridPos := component.GridPos.Get(playerEntry)
		if gridPos.Row < 3 {
			// gridPos.Row += 1
			targetX, targetY = GridCoord2Screen(gridPos.Row+1, gridPos.Col)
		}
		component.GridPos.Set(playerEntry, gridPos)
	}
	if !p.isAnim && inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		gridPos := component.GridPos.Get(playerEntry)
		if gridPos.Col > 0 {
			// gridPos.Row += 1
			targetX, targetY = GridCoord2Screen(gridPos.Row, gridPos.Col-1)
		}
		component.GridPos.Set(playerEntry, gridPos)
	}
	if !p.isAnim && inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		gridPos := component.GridPos.Get(playerEntry)
		if gridPos.Col < 3 {
			// gridPos.Row += 1
			targetX, targetY = GridCoord2Screen(gridPos.Row, gridPos.Col+1)
		}
		component.GridPos.Set(playerEntry, gridPos)
	}
	screenPos := component.ScreenPos.Get(playerEntry)
	if targetX != 0 && targetY != 0 {
		p.isAnim = true
		component.TargetLocation.Set(playerEntry, &component.MoveTargetData{Tx: targetX, Ty: targetY})
		vx := targetX - screenPos.X
		vy := targetY - screenPos.Y
		if vx != 0 || vy != 0 {
			speedVector := csg.NewVector(vx, vy, 0)
			speedVector = speedVector.Normalize().MultiplyScalar(DefaultSpeed)
			component.Speed.Set(playerEntry, &component.SpeedData{Vx: speedVector.X, Vy: speedVector.Y})
		}
		component.TargetLocation.Set(playerEntry, &component.MoveTargetData{Tx: targetX, Ty: targetY})
	}
	targetLoc := component.TargetLocation.Get(playerEntry)
	speedComponent := component.Speed.Get(playerEntry)
	screenPos.X += speedComponent.Vx
	screenPos.Y += speedComponent.Vy
	component.ScreenPos.Set(playerEntry, screenPos)
	// fmt.Println(screenPos)
	if screenPos.X == targetLoc.Tx && screenPos.Y == targetLoc.Ty {
		i := component.Speed.Get(playerEntry)
		i.Vx = 0
		i.Vy = 0
		p.isAnim = false
		component.Speed.Set(playerEntry, i)
		if assets.TileHeight > 0 {
			col, row := Coord2Grid(targetLoc.Tx, targetLoc.Ty)
			comp := component.GridPos.Get(playerEntry)
			comp.Col = col
			comp.Row = row
		}

	}

}
