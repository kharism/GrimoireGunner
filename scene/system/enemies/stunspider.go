package enemies

import (
	"math/rand"
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func NewStunSpider(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.StunSpider)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 600, Name: "StunSpider"})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: StunSpiderRoutine, Memory: data})
}
func StunSpiderRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		component.Sprite.Set(entity, &component.SpriteData{Image: assets.StunSpider})
		memory[WARM_UP] = time.Now().Add(2 * time.Second)
	}
	gridPos := component.GridPos.Get(entity)
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "MOVE1"
		}
	}
	if memory[CURRENT_STRATEGY] == "MOVE1" {
		component.Sprite.Set(entity, &component.SpriteData{Image: assets.StunSpider})
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			moveCount := memory[MOVE_COUNT].(int)
			if moveCount < 4 {

				scrPos := component.ScreenPos.Get(entity)
			checkMove:
				for {
					rndMove := rand.Int() % 4
					switch rndMove {
					case 0:
						if validMove(ecs, gridPos.Row-1, gridPos.Col) {
							gridPos.Row -= 1
							scrPos.Y -= float64(assets.TileHeight)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					case 1:
						if validMove(ecs, gridPos.Row, gridPos.Col-1) && gridPos.Col-1 >= 4 {
							gridPos.Col -= 1
							scrPos.X -= float64(assets.TileWidth)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					case 2:
						if validMove(ecs, gridPos.Row+1, gridPos.Col) {
							gridPos.Row += 1
							scrPos.Y += float64(assets.TileHeight)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					case 3:
						if validMove(ecs, gridPos.Row, gridPos.Col+1) {
							gridPos.Col += 1
							scrPos.X += float64(assets.TileWidth)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					}
				}
			} else {
				memory[MOVE_COUNT] = 0
				memory[CURRENT_STRATEGY] = "ATTACK1"
				component.Sprite.Set(entity, &component.SpriteData{Image: assets.StunSpiderWarmup2})
				memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
			}
		}
	}
	playerZones := [16][2]int{
		{0, 0}, {0, 1}, {0, 2}, {0, 3},
		{1, 0}, {1, 1}, {1, 2}, {1, 3},
		{2, 0}, {2, 1}, {2, 2}, {2, 3},
		{3, 0}, {3, 1}, {3, 2}, {3, 3},
	}
	if memory[CURRENT_STRATEGY] == "ATTACK1" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Set(entity, &component.SpriteData{Image: assets.StunSpider})
			rand.Shuffle(16, func(i, j int) {
				playerZones[i], playerZones[j] = playerZones[j], playerZones[i]
			})
			targets := playerZones[0:3]
			for _, val := range targets {
				createTargetWeb(ecs, val[0], val[1], dmg)
			}
			memory[CURRENT_STRATEGY] = "MOVE2"
		}
	}
	if memory[CURRENT_STRATEGY] == "MOVE2" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			playerGrid, _ := attack.GetPlayerGridPos(ecs)
			if playerGrid == nil {
				return
			}
			demonPos := component.GridPos.Get(entity)
			tempRow := playerGrid.Row
			tempCol := 4
			for {
				if validMove(ecs, tempRow, tempCol) {
					demonPos.Row = tempRow
					demonPos.Col = tempCol
					scrPos := component.ScreenPos.Get(entity)
					scrPos.X = 0
					scrPos.Y = 0
					break
				} else if tempCol < 8 {
					tempCol += 1
				} else {
					tempRow += 1
					tempCol = 0
				}
			}

			memory[CURRENT_STRATEGY] = "ATTACK2"
			component.Sprite.Get(entity).Image = assets.StunSpiderWarmup2
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK2" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			demonPos := component.GridPos.Get(entity)
			component.Sprite.Get(entity).Image = assets.StunSpider
			shootWeb(ecs, demonPos.Row, demonPos.Col-1, dmg+10)
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
			memory[CURRENT_STRATEGY] = "WAIT"
		}
	}

}
func createTargetWeb(_ecs *ecs.ECS, row, col, dmg int) {
	target := _ecs.World.Create(component.GridPos, component.GridTarget, component.Transient)
	entry := _ecs.World.Entry(target)
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.Transient.Set(entry, &component.TransientData{
		Start:    time.Now(),
		Duration: 400 * time.Millisecond,
		OnRemoveCallback: func(ecs *ecs.ECS, entity *donburi.Entry) {
			createWeb(ecs, row, col, dmg)
		},
	})
}
func lightingWebHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	attack.SingleHitProjectile(ecs, projectile, receiver)
	receiver.AddComponent(component.Root)
	component.EventQueue.AddEvent(&clearRootEvent{Target: receiver, Time: time.Now().Add(1 * time.Second)})
}

type clearRootEvent struct {
	Target *donburi.Entry
	Time   time.Time
}

func (cc *clearRootEvent) Execute(ecs *ecs.ECS) {
	cc.Target.RemoveComponent(component.Root)
}
func (cc *clearRootEvent) GetTime() time.Time {
	return cc.Time
}
func shootWeb(_ecs *ecs.ECS, row, col, dmg int) {
	projectile := archetype.NewProjectile(
		_ecs.World, archetype.ProjectileParam{
			Col:    col,
			Row:    row,
			Damage: dmg,
			Vx:     -15,
			Vy:     0,
			Sprite: assets.ElecWeb,
			OnHit:  lightingWebHit,
		},
	)
	entry := _ecs.World.Entry(*projectile)
	entry.AddComponent(component.Elements)
	component.Elements.SetValue(entry, component.ELEC)
}
func createWeb(_ecs *ecs.ECS, row, col, dmg int) {
	target := _ecs.World.Create(component.GridPos, component.OnHit, component.Elements, component.Damage, component.Fx, component.Transient)
	entry := _ecs.World.Entry(target)
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.Elements.SetValue(entry, component.ELEC)
	component.Damage.Set(entry, &component.DamageData{Damage: dmg})
	scrX, scrY := assets.GridCoord2Screen(row, col)
	scrX -= 50
	scrY -= 100
	component.Fx.Set(entry, &component.FxData{
		Animation: core.NewMovableImage(assets.Web,
			core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: scrX, Sy: scrY}),
		),
	})
	component.OnHit.SetValue(entry, lightingWebHit)
	component.Transient.Set(entry, &component.TransientData{
		Start:            time.Now(),
		Duration:         2000 * time.Millisecond,
		OnRemoveCallback: nil,
	})
}
