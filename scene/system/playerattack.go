package system

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/events"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type playerAttackSystem struct {
	PlayerIndex     *donburi.Entity
	returnToStandby time.Time
	// State change the behaviour of action key press, the default ones will be combatState
	State func(ecs *ecs.ECS, s *playerAttackSystem)
}

//	func NewPlayerAttackSystem(player *donburi.Entity) *PlayerAttackSystem {
//		return &PlayerAttackSystem{PlayerIndex: player}
//	}
var PlayerAttackSystem = playerAttackSystem{State: CombatState}

var timerDelay = time.Now()

func DoNothingState(ecs *ecs.ECS, s *playerAttackSystem) {

}
func CombatState(ecs *ecs.ECS, s *playerAttackSystem) {
	if ebiten.IsKeyPressed(ebiten.KeyE) { //ebiten.IsKeyPressed(ebiten.KeyE) {
		if time.Now().Sub(timerDelay) > 200*time.Millisecond {
			playerId := ecs.World.Entry(*s.PlayerIndex)
			gridPos := component.GridPos.Get(playerId)
			component.Sprite.Set(playerId, &component.SpriteData{Image: assets.Player1Attack})
			s.returnToStandby = time.Now().Add(500 * time.Millisecond)
			attack.GenerateMagibullet(ecs, gridPos.Row, gridPos.Col+1, 25)
			timerDelay = time.Now()
		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		playerId := ecs.World.Entry(*s.PlayerIndex)
		// gridPos := component.GridPos.Get(playerId)
		// scrPos := component.ScreenPos.Get(playerId)
		component.Sprite.Set(playerId, &component.SpriteData{Image: assets.Player1Attack})
		s.returnToStandby = time.Now().Add(500 * time.Millisecond)
		// attack.NewLongSwordAttack(EnergySystem, ecs, *scrPos, *gridPos)
		if len(loadout.CurLoadOut) >= 2 && loadout.CurLoadOut[1] != nil {
			if !loadout.CurLoadOut[1].GetCooldown().IsZero() && loadout.CurLoadOut[1].GetCooldown().Before(time.Now()) {
				loadout.CurLoadOut[1].Cast(EnergySystem, ecs)
				if vv, ok := loadout.CurLoadOut[1].(Consumables); ok {
					if vv.GetCharge() == 0 {
						loadout.CurLoadOut[1] = nil
					}
				}
			}

		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		playerId := ecs.World.Entry(*s.PlayerIndex)
		// gridPos := component.GridPos.Get(playerId)
		component.Sprite.Set(playerId, &component.SpriteData{Image: assets.Player1Attack})
		s.returnToStandby = time.Now().Add(500 * time.Millisecond)
		// attack.NewLigtningAttack(ecs, attack.LightnigAtkParam{
		// 	StartRow:  gridPos.Row,
		// 	StartCol:  gridPos.Col + 1,
		// 	Direction: 1,
		// 	Actor:     playerId,
		// })
		if len(loadout.CurLoadOut) >= 1 && loadout.CurLoadOut[0] != nil {
			if !loadout.CurLoadOut[0].GetCooldown().IsZero() && loadout.CurLoadOut[0].GetCooldown().Before(time.Now()) {
				loadout.CurLoadOut[0].Cast(EnergySystem, ecs)
				if vv, ok := loadout.CurLoadOut[0].(Consumables); ok {
					if vv.GetCharge() == 0 {
						loadout.CurLoadOut[0] = nil
					}
				}
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		temp := loadout.CurLoadOut
		loadout.CurLoadOut = loadout.SubLoadOut1
		loadout.SubLoadOut1 = loadout.SubLoadOut2
		loadout.SubLoadOut2 = temp
	}
	if time.Now().After(s.returnToStandby) {
		playerId := ecs.World.Entry(*s.PlayerIndex)
		component.Sprite.Set(playerId, &component.SpriteData{Image: assets.Player1Stand})
	}
}
func CombatClearState(ecs *ecs.ECS, s *playerAttackSystem) {
	if inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		events.CombatClearEvent.Publish(ecs.World, events.CombatClearData{})
		events.CombatClearEvent.ProcessEvents(ecs.World)
	}

}
func (s *playerAttackSystem) Update(ecs *ecs.ECS) {
	s.State(ecs, s)
}
