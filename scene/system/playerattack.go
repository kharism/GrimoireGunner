package system

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type PlayerAttackSystem struct {
	PlayerIndex     *donburi.Entity
	returnToStandby time.Time
}

func NewPlayerAttackSystem(player *donburi.Entity) *PlayerAttackSystem {
	return &PlayerAttackSystem{PlayerIndex: player}
}

// var timerDelay = time.Now()

func (s *PlayerAttackSystem) Update(ecs *ecs.ECS) {
	if inpututil.IsKeyJustReleased(ebiten.KeyE) { //ebiten.IsKeyPressed(ebiten.KeyE) {
		// if time.Now().Sub(timerDelay) > 500*time.Millisecond {
		playerId := ecs.World.Entry(*s.PlayerIndex)
		gridPos := component.GridPos.Get(playerId)
		component.Sprite.Set(playerId, &component.SpriteData{Image: assets.Player1Attack})
		s.returnToStandby = time.Now().Add(500 * time.Millisecond)
		attack.GenerateMagibullet(ecs, gridPos.Row, gridPos.Col+1, 15)
		// timerDelay = time.Now()
		// }

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		playerId := ecs.World.Entry(*s.PlayerIndex)
		gridPos := component.GridPos.Get(playerId)
		scrPos := component.ScreenPos.Get(playerId)
		component.Sprite.Set(playerId, &component.SpriteData{Image: assets.Player1Attack})
		s.returnToStandby = time.Now().Add(500 * time.Millisecond)
		attack.NewLongSwordAttack(ecs, *scrPos, *gridPos)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		playerId := ecs.World.Entry(*s.PlayerIndex)
		gridPos := component.GridPos.Get(playerId)
		component.Sprite.Set(playerId, &component.SpriteData{Image: assets.Player1Attack})
		s.returnToStandby = time.Now().Add(500 * time.Millisecond)
		attack.NewLigtningAttack(ecs, attack.LightnigAtkParam{
			StartRow:  gridPos.Row,
			StartCol:  gridPos.Col + 1,
			Direction: 1,
			Actor:     playerId,
		})
	}
	if time.Now().After(s.returnToStandby) {
		playerId := ecs.World.Entry(*s.PlayerIndex)
		component.Sprite.Set(playerId, &component.SpriteData{Image: assets.Player1Stand})
	}
}
