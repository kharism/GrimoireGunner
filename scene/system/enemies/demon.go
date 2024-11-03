package enemies

import (
	"math"
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func NewDemon(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Demon)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 500, Name: "Demon"})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})

	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 0 //it actualy damage bonus
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: DemonRoutine, Memory: data})
}

func DemonRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "MOVE"

			memory[WARM_UP] = time.Now().Add(1 * time.Second)
		}
	}
	if memory[CURRENT_STRATEGY] == "MOVE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			playerGrid, _ := attack.GetPlayerGridPos(ecs)
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

			memory[CURRENT_STRATEGY] = "WARM_UP"
			component.Sprite.Get(entity).Image = assets.DemonWarmup
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
		}
	}
	if memory[CURRENT_STRATEGY] == "WARM_UP" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.DemonCooldown
			DemonAttack(ecs, entity)
			memory[CURRENT_STRATEGY] = "WAIT"
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
			memory[CUR_DMG] = memory[CUR_DMG].(int) + 10
		}
	}
}

// this is duplicate just check the opposite direction
func shockWaveOnAtkHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	damage := component.Damage.Get(projectile).Damage
	healthComp := component.Health.Get(receiver)
	healthComp.HP -= damage
	if !receiver.HasComponent(archetype.ConstructTag) {
		healthComp.InvisTime = time.Now().Add(400 * time.Millisecond)
	} else {
		projectileGridPos := component.GridPos.Get(projectile)
		projectileGridPos.Col -= 1
		scrPosProjectile := component.ScreenPos.Get(projectile)
		scrPosProjectile.X -= float64(assets.TileWidth)
	}

	receiverPos := component.GridPos.Get(receiver)
	if validMove(ecs, receiverPos.Row, receiverPos.Col-1) {
		receiverPos.Col -= 1
		scrPos := component.ScreenPos.Get(receiver)
		scrPos.X -= 100
	}
}
func DemonAttack(ecs *ecs.ECS, entry *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entry).Memory
	bonusDmg := memory[CUR_DMG].(int)
	demonPos := component.GridPos.Get(entry)
	demonScreenPosX, demonScreenPosY := assets.GridCoord2Screen(demonPos.Row, demonPos.Col)
	playerPos, _ := attack.GetPlayerGridPos(ecs)
	playerCol := playerPos.Col
	if math.Abs(float64(playerPos.Col-demonPos.Col)) <= 1 {
		if math.Abs(float64(playerPos.Row-demonPos.Row)) <= 1 {
			// wide swing
			DAMAGE := 50 + bonusDmg
			var entry1 *donburi.Entry
			var entry2 *donburi.Entry
			var entry3 *donburi.Entry
			if demonPos.Row > 0 {
				hitbox1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
				entry1 = ecs.World.Entry(hitbox1)
				component.Damage.Set(entry1, &component.DamageData{Damage: DAMAGE})
				component.GridPos.Set(entry1, &component.GridPosComponentData{Row: demonPos.Row - 1, Col: demonPos.Col - 1})
				component.OnHit.SetValue(entry1, onReaperHit)
			}

			hitbox2 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
			entry2 = ecs.World.Entry(hitbox2)
			component.Damage.Set(entry2, &component.DamageData{Damage: DAMAGE})
			component.GridPos.Set(entry2, &component.GridPosComponentData{Row: demonPos.Row, Col: demonPos.Col - 1})
			component.OnHit.SetValue(entry2, onReaperHit)

			if demonPos.Row < 3 {
				hitbox3 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
				entry3 = ecs.World.Entry(hitbox3)
				component.Damage.Set(entry3, &component.DamageData{Damage: DAMAGE})
				component.GridPos.Set(entry3, &component.GridPosComponentData{Row: demonPos.Row + 1, Col: demonPos.Col - 1})
				component.OnHit.SetValue(entry3, onReaperHit)
			}

			fx := ecs.World.Create(component.Fx)
			fxEntry := ecs.World.Entry(fx)
			wideSlash := assets.NewWideSlashAtkAnim(assets.SpriteParam{
				ScreenX: demonScreenPosX - float64(assets.TileWidth)*1.5,
				ScreenY: demonScreenPosY - float64(assets.TileHeight)*2,
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
		} else {
			// column wide attack
			DAMAGE := 40 + bonusDmg
			for i := 0; i < 4; i++ {
				targetScrX, targetScrY := assets.GridCoord2Screen(i, playerCol)
				targetScrX -= 50
				targetScrY -= 100
				entity := ecs.World.Create(component.Burner, component.Damage, component.GridPos, component.Transient, component.Fx)
				entry := ecs.World.Entry(entity)
				component.Burner.Set(entry, &component.BurnerData{
					Damage: DAMAGE,
				})
				component.Damage.Set(entry, &component.DamageData{Damage: DAMAGE})
				component.GridPos.Set(entry, &component.GridPosComponentData{Col: playerCol, Row: i})
				component.Transient.Set(entry, &component.TransientData{Start: time.Now(), Duration: 5 * time.Second})
				// component.OnHit.SetValue(entry, attack.OnTowerHit)
				flameTower := core.NewMovableImage(assets.FlametowerRaw, core.NewMovableImageParams().
					WithMoveParam(core.MoveParam{Sx: targetScrX, Sy: targetScrY, Speed: 3}))
				component.Fx.Set(entry, &component.FxData{Animation: flameTower})
			}
		}

		memory[WARM_UP] = time.Now().Add(time.Second)
	} else {
		// ranged atk
		DAMAGE := 30 + bonusDmg
		shockwave := ecs.World.Create(
			component.GridPos,
			component.ScreenPos,
			component.Speed,
			component.Damage,
			// component.Sprite,
			component.OnHit,
			component.Fx,
			component.TargetLocation,
		)
		shockwaveEntry := ecs.World.Entry(shockwave)
		component.GridPos.Set(shockwaveEntry, &component.GridPosComponentData{
			Row: demonPos.Row,
			Col: demonPos.Col - 1,
		})
		screenTargetX, screenTargetY := assets.GridCoord2Screen(demonPos.Row, -1)
		component.TargetLocation.Set(shockwaveEntry, &component.MoveTargetData{
			Tx: screenTargetX,
			Ty: screenTargetY,
		})
		component.OnHit.SetValue(shockwaveEntry, shockWaveOnAtkHit)
		SPEED := -5.5
		component.Speed.Set(shockwaveEntry, &component.SpeedData{Vx: SPEED, Vy: 0})
		component.Damage.Set(shockwaveEntry, &component.DamageData{Damage: DAMAGE})
		screenX, screenY := assets.GridCoord2Screen(demonPos.Row, demonPos.Col-1)
		screenX = screenX - 50
		screenY = screenY - 100
		shockwaveAnim := assets.NewShockwaveAnim(assets.SpriteParam{
			ScreenX: screenX,
			ScreenY: screenY,
			Modulo:  10,

			Done: func() {
				// ecs.World.Remove(shockwave)
			},
		})
		shockwaveAnim.FlipHorizontal = true
		shockwaveAnim.AddAnimation(
			core.NewMoveAnimationFromParam(
				core.MoveParam{
					Tx:    screenTargetX,
					Ty:    screenY,
					Speed: 5,
				}),
		)
		component.Fx.Set(shockwaveEntry, &component.FxData{Animation: shockwaveAnim})
	}
}
