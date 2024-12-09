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

func NewPoacher(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Poacher)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 700, MaxHP: 700, Name: "Poacher"})

	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: PoacherRoutine, Memory: data})
}
func PoacherRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	gridPos := component.GridPos.Get(entity)
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "MOVE"
		}
	}
	if memory[CURRENT_STRATEGY] == "MOVE" {
		component.Sprite.Set(entity, &component.SpriteData{Image: assets.Poacher})
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			moveCount := memory[MOVE_COUNT].(int)
			if moveCount < 2 {

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
				component.Sprite.Set(entity, &component.SpriteData{Image: assets.PoacherCooldown})
				memory[WARM_UP] = time.Now().Add(600 * time.Millisecond)
			}
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK1" {
		// throw 2 column of traps
		sourceScrX, sourceSrcY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
		sourceScrX -= 50
		sourceSrcY -= 50

		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			if memory[MOVE_COUNT].(int) < 2 {
				col := memory[MOVE_COUNT].(int)
				for i := 0; i < 4; i++ {
					targetScrX, targetScrY := assets.GridCoord2Screen(i, memory[MOVE_COUNT].(int))
					targetScrX -= 50
					targetScrY -= 100
					row := i
					movableImg := core.NewMovableImage(assets.Projectile1, core.NewMovableImageParams().WithMoveParam(
						core.MoveParam{Sx: sourceScrX, Sy: sourceSrcY, Tx: targetScrX, Ty: targetScrY, Speed: 10}))
					fx := ecs.World.Create(component.Fx)
					entry := ecs.World.Entry(fx)
					movableImg.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{Tx: targetScrX, Ty: targetScrY, Speed: 10}))
					movableImg.Done = func() {
						ecs.World.Remove(fx)
						entity := ecs.World.Create(component.Sprite, component.Damage, component.GridPos, component.Transient, component.OnHit, component.ScreenPos)
						entry := ecs.World.Entry(entity)
						component.Damage.Set(entry, &component.DamageData{Damage: dmg})
						component.Sprite.Set(entry, &component.SpriteData{Image: assets.BearTrap})
						component.GridPos.Set(entry, &component.GridPosComponentData{Col: col, Row: row})
						component.Transient.Set(entry, &component.TransientData{Start: time.Now(), Duration: 3 * time.Second})
						component.OnHit.SetValue(entry, attack.SingleHitProjectile)
					}
					component.Fx.Set(entry, &component.FxData{Animation: movableImg})
				}
				memory[MOVE_COUNT] = memory[MOVE_COUNT].(int) + 1
			} else {
				memory[MOVE_COUNT] = 0
				memory[CURRENT_STRATEGY] = "ATTACK2"
				component.Sprite.Set(entity, &component.SpriteData{Image: assets.Poacher})
				memory[WARM_UP] = time.Now().Add(300 * time.Millisecond)
			}

		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK2" {
		// move in front of player
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			playerPos, _ := attack.GetPlayerGridPos(ecs)
			if playerPos == nil {
				return
			}
			scrPos := component.ScreenPos.Get(entity)
			component.Sprite.Get(entity).Image = assets.PoacherWarmup
			newCol := playerPos.Col + 1
			newRow := playerPos.Row
			if validMove(ecs, newRow, newCol) || (gridPos.Col == playerPos.Col+1 && gridPos.Row == playerPos.Row) {
				gridPos.Row = newRow
				gridPos.Col = newCol
				scrPos.Y = 0
				scrPos.X = 0
				now := time.Now()
				memory[CURRENT_STRATEGY] = "ATTACK3"

				memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
				target2 := ecs.World.Create(component.GridPos, component.GridTarget, component.Transient)
				entry2 := ecs.World.Entry(target2)
				component.Transient.Set(entry2, &component.TransientData{Start: now, Duration: 745 * time.Millisecond})
				component.GridPos.Set(entry2, &component.GridPosComponentData{Col: gridPos.Col - 1, Row: gridPos.Row})
			}
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK3" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Set(entity, &component.SpriteData{Image: assets.PoacherCooldown})
			// reaperScreenX, reaperScreenY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
			component.Sprite.Get(entity).Image = assets.PoacherCooldown
			hitbox2 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Transient)
			var entry2 *donburi.Entry
			entry2 = ecs.World.Entry(hitbox2)

			component.Damage.Set(entry2, &component.DamageData{Damage: dmg + 20})
			component.GridPos.Set(entry2, &component.GridPosComponentData{Row: gridPos.Row, Col: gridPos.Col - 1})
			component.OnHit.SetValue(entry2, attack.SingleHitProjectile)
			component.Transient.Set(entry2, &component.TransientData{Start: time.Now(), Duration: 200 * time.Millisecond})
			memory[CURRENT_STRATEGY] = "RETURN"
			memory[WARM_UP] = time.Now().Add(200 * time.Millisecond)
		}
	}
	if memory[CURRENT_STRATEGY] == "RETURN" {
		scrPos := component.ScreenPos.Get(entity)
	checkMove2:
		for {
			// rndMove := rand.Int() % 4
			newCol := 4 + rand.Int()%4
			newRow := rand.Int() % 4
			if validMove(ecs, newRow, newCol) {
				gridPos.Row = newRow
				gridPos.Col = newCol
				scrPos.Y = 0
				scrPos.X = 0
				memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)

				break checkMove2
			}
		}
		memory[CURRENT_STRATEGY] = "MOVE"

	}
}
