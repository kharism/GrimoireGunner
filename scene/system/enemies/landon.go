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

func NewLandon(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Landon1)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	// component.Sprite.Get(entry).Scale = &core.ScaleParam{Sx: 1.5, Sy: 1.5}
	// entry.AddComponent(component.Shader)
	component.Health.Set(entry, &component.HealthData{HP: 1400, Name: "Landon"})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	// component.Shader.Set(entry, assets.DakkaShader)
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[MOVE_COUNT] = 0
	data[CURRENT_STRATEGY] = ""
	data[CUR_DMG] = 20
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: LandonRoutine, Memory: data})
}

const OLDCOL = "OLDCOL"
const OLDROW = "OLDROW"

func LandonRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "MOVE"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	gridPos := component.GridPos.Get(entity)
	if memory[CURRENT_STRATEGY] == "MOVE" {
		component.Sprite.Set(entity, &component.SpriteData{Image: assets.Landon1})
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
				moveCount = 0
				ii := rand.Int() % 3
				if ii == 0 {
					memory[CURRENT_STRATEGY] = "ATTACK_RANGED"
					component.Sprite.Set(entity, &component.SpriteData{Image: assets.LandonShoot})
					// component.Sprite.Get(entity)
					memory[WARM_UP] = time.Now().Add(300 * time.Millisecond)
				} else if ii == 1 {
					memory[CURRENT_STRATEGY] = "ATTACK_MELEE_1_WARMP"
					component.Sprite.Set(entity, &component.SpriteData{Image: assets.LandonWarmup1})
					memory[WARM_UP] = time.Now().Add(300 * time.Millisecond)
				} else if ii == 2 {
					memory[CURRENT_STRATEGY] = "ATTACK_MELEE_2_WARMP"
					component.Sprite.Set(entity, &component.SpriteData{Image: assets.LandonWarmup2})
					memory[WARM_UP] = time.Now().Add(300 * time.Millisecond)
				}

			}
		}
	}
	playerPos, _ := attack.GetPlayerGridPos(ecs)
	if playerPos == nil {
		return
	}
	if memory[CURRENT_STRATEGY] == "ATTACK_RANGED" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			attack.AtkSfxQueue.QueueSFX(assets.MagibulletFx)
			bulletEntity := archetype.NewProjectile(ecs.World, archetype.ProjectileParam{
				Vx:     -25,
				Vy:     0,
				Col:    gridPos.Col - 1,
				Row:    gridPos.Row,
				Damage: dmg,
				Sprite: assets.Projectile1,
				OnHit:  attack.SingleHitProjectile,
			})
			bulletEntry := ecs.World.Entry(*bulletEntity)
			bulletEntry.AddComponent(component.Preremove)
			component.Preremove.SetValue(bulletEntry, landonBulletRemove)
			memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
			memory[CURRENT_STRATEGY] = "MOVE"
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK_MELEE_1_WARMP" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			scrPos := component.ScreenPos.Get(entity)
			newCol := playerPos.Col + 1
			newRow := playerPos.Row
			if validMove(ecs, newRow, newCol) || (newCol == gridPos.Col && newRow == gridPos.Row) {
				memory[OLDCOL] = gridPos.Col
				memory[OLDROW] = gridPos.Row
				gridPos.Row = newRow
				gridPos.Col = newCol
				scrPos.Y = 0
				scrPos.X = 0
				now := time.Now()
				memory[CURRENT_STRATEGY] = "ATTACK_MELEE_1_COOLDOWN"
				startGridPosX := gridPos.Col - 2
				startGridPosY := gridPos.Row - 1

				for i := 0; i < 2; i++ {
					if startGridPosX+i < 0 {
						continue
					}
					for j := 0; j < 3; j++ {
						if startGridPosY+j < 0 || startGridPosY+j > 3 {
							continue
						}
						target2 := ecs.World.Create(component.GridPos, component.GridTarget, component.Transient)
						entry2 := ecs.World.Entry(target2)
						component.Transient.Set(entry2, &component.TransientData{Start: now, Duration: 745 * time.Millisecond})
						component.GridPos.Set(entry2, &component.GridPosComponentData{Row: startGridPosY + j, Col: startGridPosX + i})
					}
				}
				memory[WARM_UP] = time.Now().Add(600 * time.Millisecond)
			}
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK_MELEE_1_COOLDOWN" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Set(entity, &component.SpriteData{Image: assets.LandonCooldown1})
			// component.Sprite.Set(entity, &component.SpriteData{Image: assets.SwordswomenCooldown})
			reaperScreenX, reaperScreenY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
			// component.Sprite.Get(entity).Image = assets.SwordswomenCooldown
			startGridPosX := gridPos.Col - 2
			startGridPosY := gridPos.Row - 1
			damageGrids := []donburi.Entity{}
			for i := 0; i < 2; i++ {
				if startGridPosX+i < 0 {
					continue
				}
				for j := 0; j < 3; j++ {
					if startGridPosY+j < 0 || startGridPosY+j > 3 {
						continue
					}
					hitbox := ecs.World.Create(component.Damage, component.Elements, component.GridPos, component.OnHit)
					entry := ecs.World.Entry(hitbox)
					component.Elements.SetValue(entry, component.FIRE)
					component.Damage.Set(entry, &component.DamageData{Damage: dmg + 20})
					component.Elements.SetValue(entry, component.FIRE)
					// fmt.Println("DMG", gridPos.Row+i, gridPos.Col+j)
					component.GridPos.Set(entry, &component.GridPosComponentData{Row: startGridPosY + j, Col: startGridPosX + i})
					component.OnHit.SetValue(entry, onReaperHit)
					damageGrids = append(damageGrids, hitbox)
				}
			}
			// fmt.Println("=====")
			fx := ecs.World.Create(component.Fx)
			fxEntry := ecs.World.Entry(fx)
			wideSlash := assets.NewWiderSlashAtkAnim(assets.SpriteParam{
				ScreenX: reaperScreenX - float64(assets.TileWidth)*1.5,
				ScreenY: reaperScreenY - float64(assets.TileHeight)*2,
				Done: func() {
					for _, d := range damageGrids {
						ecs.World.Remove(d)
					}
					ecs.World.Remove(fx)
				},
				Modulo: 3,
			})
			component.Fx.Set(fxEntry, &component.FxData{Animation: wideSlash})
			memory[CURRENT_STRATEGY] = "RETURN"
			memory[WARM_UP] = time.Now().Add(300 * time.Millisecond)
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK_MELEE_2_WARMP" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			scrPos := component.ScreenPos.Get(entity)
			newCol := playerPos.Col + 1
			newRow := playerPos.Row
			if validMove(ecs, newRow, newCol) || (newCol == gridPos.Col && newRow == gridPos.Row) {
				memory[OLDCOL] = gridPos.Col
				memory[OLDROW] = gridPos.Row
				gridPos.Row = newRow
				gridPos.Col = newCol
				scrPos.Y = 0
				scrPos.X = 0
				memory[CURRENT_STRATEGY] = "ATTACK_MELEE_2_COOLDOWN"

				memory[WARM_UP] = time.Now().Add(600 * time.Millisecond)
			}
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK_MELEE_2_COOLDOWN" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			LightningAttack(ecs, entity)
			component.Sprite.Set(entity, &component.SpriteData{Image: assets.LandonCooldown2})
			memory[CURRENT_STRATEGY] = "RETURN"
			memory[WARM_UP] = time.Now().Add(300 * time.Millisecond)
		}
	}
	if memory[CURRENT_STRATEGY] == "RETURN" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			scrPos := component.ScreenPos.Get(entity)
			gridPos.Col = memory[OLDCOL].(int)
			gridPos.Row = memory[OLDROW].(int)
			scrPos.Y = 0
			scrPos.X = 0
			component.Sprite.Set(entity, &component.SpriteData{Image: assets.LandonCooldown1})
			memory[CURRENT_STRATEGY] = "MOVE"
			memory[CUR_DMG] = dmg + 10
		}
	}
}
func landonBulletRemove(ecs *ecs.ECS, entity *donburi.Entry) {
	gridPos := component.GridPos.Get(entity)
	posX := gridPos.Col
	posY := gridPos.Row
	now := time.Now()
	for i := -1; i <= 1; i++ {
		if posY+i < 0 || posY+i >= 4 {
			continue
		}
		//bamboo lance on player's back
		dmgtile := ecs.World.Create(component.Elements, component.Damage, component.GridPos, component.Transient, component.OnHit)
		dmgEntry := ecs.World.Entry(dmgtile)
		component.GridPos.Set(dmgEntry, &component.GridPosComponentData{
			Row: posY + i,
			Col: posX + 1,
		})
		component.Transient.Set(dmgEntry, &component.TransientData{
			Start:    now,
			Duration: 200 * time.Millisecond,
		})
		component.Elements.SetValue(dmgEntry, component.WOOD)
		component.Damage.Set(dmgEntry, &component.DamageData{
			Damage: 30,
		})
		component.OnHit.SetValue(dmgEntry, attack.SingleHitProjectile)
		fxEntity := ecs.World.Create(component.Fx, component.Transient)
		bambooLanceFx := ecs.World.Entry(fxEntity)
		sx, sy := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
		sx -= 50
		sy -= 100
		bounds := assets.BambooLance.Bounds()
		width := bounds.Dx()
		gg := core.NewMovableImageParams().WithMoveParam(
			core.MoveParam{Sx: sx + float64(width), Sy: sy},
		)
		fx := core.NewMovableImage(assets.BambooLance, gg)
		component.Fx.Set(bambooLanceFx, &component.FxData{Animation: fx})
		component.Transient.Set(bambooLanceFx, &component.TransientData{
			Start:    now,
			Duration: 200 * time.Millisecond,
		})
	}
}
