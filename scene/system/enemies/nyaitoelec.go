package enemies

import (
	"math/rand"
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func NewNyaaito(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Swordswomen)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	entry.AddComponent(component.Shader)
	component.Health.Set(entry, &component.HealthData{HP: 1000, Name: "Swordwoman", Element: component.ELEC})
	component.Shader.Set(entry, assets.Shocky2Shader)
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: NyaitoRoutine, Memory: data})
}
func NyaitoRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		component.Sprite.Set(entity, &component.SpriteData{Image: assets.Swordswomen})
		memory[WARM_UP] = time.Now().Add(2 * time.Second)
	}
	gridPos := component.GridPos.Get(entity)
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "MOVE"
		}
	}
	moveCount := memory[MOVE_COUNT].(int)
	if memory[CURRENT_STRATEGY] == "MOVE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			if moveCount < 4 {
				scrPos := component.ScreenPos.Get(entity)
			checkMove:
				for {
					// rndMove := rand.Int() % 4
					newCol := 4 + rand.Int()%4
					newRow := rand.Int() % 4
					if validMove(ecs, newRow, newCol) {
						gridPos.Row = newRow
						gridPos.Col = newCol
						scrPos.Y = 0
						scrPos.X = 0
						memory[MOVE_COUNT] = moveCount + 1
						memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)

						break checkMove
					}

				}
			} else {
				// warmup for attack
				memory[MOVE_COUNT] = 0
				moveCount = 0
				ii := 1 //rand.Int() % 2
				if ii == 0 {
					memory[CURRENT_STRATEGY] = "ATTACK_RANGED"
					component.Sprite.Set(entity, &component.SpriteData{Image: assets.SwordswomenShoot})
					// component.Sprite.Get(entity)
				} else {
					memory[CURRENT_STRATEGY] = "ATTACK_MELEE"
					component.Sprite.Set(entity, &component.SpriteData{Image: assets.SwordswomenWarmup})
				}
				memory[WARM_UP] = time.Now().Add(300 * time.Millisecond)

			}
		}

	}
	playerPos, _ := attack.GetPlayerGridPos(ecs)
	if playerPos == nil {
		return
	}
	if memory[CURRENT_STRATEGY] == "ATTACK_RANGED" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			if playerPos.Row == gridPos.Row {
				// single projectile
				archetype.NewProjectile(ecs.World, archetype.ProjectileParam{
					Vx:     -15,
					Vy:     0,
					Col:    gridPos.Col - 1,
					Row:    gridPos.Row,
					Damage: dmg - 30,
					Sprite: assets.Projectile1,
					OnHit:  attack.SingleHitProjectile,
				})
			} else {
				// swipe attack
				rows := []int{0, 1, 2, 3}
				rand.Shuffle(4, func(i, j int) {
					rows[i], rows[j] = rows[j], rows[i]
				})
				shootTime := time.Now().Add(1 * time.Second)

				for i := 0; i < 3; i++ {
					entity := archetype.NewProjectile(ecs.World, archetype.ProjectileParam{
						Vx:     0,
						Vy:     0,
						Col:    gridPos.Col - 1,
						Row:    rows[i],
						Damage: dmg - 30,
						Sprite: assets.Projectile2,
						OnHit:  attack.SingleHitProjectile,
					})
					jj := &moveDagger{entry: ecs.World.Entry(*entity), Time: shootTime}
					component.EventQueue.AddEvent(jj)
				}
			}
			memory[CURRENT_STRATEGY] = ""
			// memory[WARM_UP] = time.Now()
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK_MELEE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			if moveCount < 4 {
				component.Sprite.Set(entity, &component.SpriteData{Image: assets.SwordswomenWarmup})
				scrPos := component.ScreenPos.Get(entity)
				newCol := playerPos.Col + 1
				newRow := playerPos.Row
				if validMove(ecs, newRow, newCol) {
					gridPos.Row = newRow
					gridPos.Col = newCol
					scrPos.Y = 0
					scrPos.X = 0
					memory[MOVE_COUNT] = moveCount + 1
					memory[CURRENT_STRATEGY] = "ATTACK_MELEE_1"
					memory[WARM_UP] = time.Now().Add(750 * time.Millisecond)
					now := time.Now()
					if gridPos.Row > 1 {
						target1 := ecs.World.Create(component.GridPos, component.GridTarget, component.Transient)
						entry1 := ecs.World.Entry(target1)
						component.Transient.Set(entry1, &component.TransientData{Start: now, Duration: 745 * time.Millisecond})
						component.GridPos.Set(entry1, &component.GridPosComponentData{Col: gridPos.Col - 1, Row: gridPos.Row - 1})
					}
					target2 := ecs.World.Create(component.GridPos, component.GridTarget, component.Transient)
					entry2 := ecs.World.Entry(target2)
					component.Transient.Set(entry2, &component.TransientData{Start: now, Duration: 745 * time.Millisecond})
					component.GridPos.Set(entry2, &component.GridPosComponentData{Col: gridPos.Col - 1, Row: gridPos.Row})
					if gridPos.Row < 3 {
						target3 := ecs.World.Create(component.GridPos, component.GridTarget, component.Transient)
						entry3 := ecs.World.Entry(target3)
						component.Transient.Set(entry3, &component.TransientData{Start: now, Duration: 745 * time.Millisecond})
						component.GridPos.Set(entry3, &component.GridPosComponentData{Col: gridPos.Col - 1, Row: gridPos.Row + 1})
					}

				}

			} else {
				memory[CURRENT_STRATEGY] = "WAIT"
				memory[CUR_DMG] = dmg + 4
				memory[WARM_UP] = time.Now().Add(800 * time.Millisecond)
			}
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK_MELEE_1" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Set(entity, &component.SpriteData{Image: assets.SwordswomenCooldown})
			reaperScreenX, reaperScreenY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
			component.Sprite.Get(entity).Image = assets.SwordswomenCooldown
			var entry1 *donburi.Entry
			var entry2 *donburi.Entry
			var entry3 *donburi.Entry
			if gridPos.Row > 0 {
				hitbox1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Elements)
				entry1 = ecs.World.Entry(hitbox1)
				component.Damage.Set(entry1, &component.DamageData{Damage: dmg + 20})
				component.Elements.SetValue(entry1, component.ELEC)
				component.GridPos.Set(entry1, &component.GridPosComponentData{Row: gridPos.Row - 1, Col: gridPos.Col - 1})
				component.OnHit.SetValue(entry1, onReaperHit)
			}
			hitbox2 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Elements)
			entry2 = ecs.World.Entry(hitbox2)
			component.Damage.Set(entry2, &component.DamageData{Damage: dmg + 20})
			component.Elements.SetValue(entry1, component.ELEC)
			component.GridPos.Set(entry2, &component.GridPosComponentData{Row: gridPos.Row, Col: gridPos.Col - 1})
			component.OnHit.SetValue(entry2, onReaperHit)
			if gridPos.Row < 3 {
				hitbox3 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Elements)
				entry3 = ecs.World.Entry(hitbox3)
				component.Damage.Set(entry3, &component.DamageData{Damage: dmg + 20})
				component.Elements.SetValue(entry1, component.ELEC)
				component.GridPos.Set(entry3, &component.GridPosComponentData{Row: gridPos.Row + 1, Col: gridPos.Col - 1})
				component.OnHit.SetValue(entry3, onReaperHit)
			}
			fx := ecs.World.Create(component.Fx)
			fxEntry := ecs.World.Entry(fx)
			wideSlash := assets.NewWideSlashAtkAnim(assets.SpriteParam{
				ScreenX: reaperScreenX - float64(assets.TileWidth)*1.5,
				ScreenY: reaperScreenY - float64(assets.TileHeight)*2,
				Done: func() {
					if entry1 != nil {
						ecs.World.Remove(entry1.Entity())
					}

					ecs.World.Remove(entry2.Entity())
					if entry3 != nil {
						ecs.World.Remove(entry3.Entity())
					}
					ecs.World.Remove(fx)
				},
				Modulo: 3,
			})
			component.Fx.Set(fxEntry, &component.FxData{Animation: wideSlash})
			memory[CURRENT_STRATEGY] = "ATTACK_MELEE"
			memory[WARM_UP] = time.Now().Add(time.Second)
		}
	}
}
