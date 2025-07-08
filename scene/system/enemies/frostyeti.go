package enemies

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func NewFrostYeti(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Yeti)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	entry.AddComponent(component.Shader)
	component.Health.Set(entry, &component.HealthData{HP: 1000, MaxHP: 1000, Name: "Yeti", Element: component.WATER})
	component.Shader.Set(entry, assets.IcyShader)

	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: FrostYetiRoutine, Memory: data})
}
func FrostYetiRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	gridPos := component.GridPos.Get(entity)
	if memory[CURRENT_STRATEGY] == "WAIT" {
		component.Sprite.Get(entity).Image = assets.Yeti
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			// memory[CURRENT_STRATEGY] = "MOVE"
			memory[CURRENT_STRATEGY] = "MOVE"
			memory[WARM_UP] = time.Now().Add(300 * time.Millisecond)
			// if gridPos.Col != 4 {
			// 	memory[CURRENT_STRATEGY] = "MOVE"
			// } else {
			// 	memory[CURRENT_STRATEGY] = "SUMMON_CONSTRUCT"
			// 	memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
			// }
		}
	}
	if memory[CURRENT_STRATEGY] == "MOVE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			hp := component.Health.Get(entity).HP
			playerGrid, _ := attack.GetPlayerGridPos(ecs)
			if playerGrid == nil {
				return
			}
			demonPos := component.GridPos.Get(entity)
			tempRow := playerGrid.Row
			// default to melee
			tempCol := 4
			if hp <= 300 {
				// go into ranged atk if HP is low
				tempCol = 6
			}
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

			if playerGrid.Col == gridPos.Col-1 {

				memory[CURRENT_STRATEGY] = "WARM_UP_MELEE"
				component.Sprite.Get(entity).Image = assets.YetiWarmup
				memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)

			} else {
				if hp > 200 {
					memory[CURRENT_STRATEGY] = "SUMMON_CONSTRUCT"
					// component.Sprite.Get(entity).Image = assets.YetiWarmup
					memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
				} else {
					memory[CURRENT_STRATEGY] = "WARM_UP_AVALANCHE"
					component.Sprite.Get(entity).Image = assets.YetiWarmup2
					memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
					filter := donburi.NewQuery(
						filter.Contains(
							component.EnemyTag,
						),
					)
					enemyCount := filter.Count(ecs.World)
					if enemyCount == 1 {
						NewHealslime(ecs, 4, 0)
					}

				}
			}

		}
	}
	if memory[CURRENT_STRATEGY] == "SUMMON_CONSTRUCT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			if validMove(ecs, gridPos.Row, gridPos.Col-1) {
				rock := archetype.NewConstruct(ecs.World, assets.Wall)
				rockEnt := ecs.World.Entry(*rock)
				rockGridPos := component.GridPos.Get(rockEnt)
				rockGridPos.Col = gridPos.Col - 1
				rockGridPos.Row = gridPos.Row
				memory[CURRENT_STRATEGY] = "WARM_UP_MELEE"
				component.Sprite.Get(entity).Image = assets.YetiWarmup
				memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
			} else {
				memory[CURRENT_STRATEGY] = "MOVE"
			}

		}
	}
	if memory[CURRENT_STRATEGY] == "WARM_UP_AVALANCHE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.YetiCooldown2
			playerPos, _ := attack.GetPlayerGridPos(ecs)
			if playerPos == nil {
				return
			}
			now := time.Now()
			for col := playerPos.Col - 1; col <= playerPos.Col+1; col++ {
				for row := playerPos.Row - 1; row <= playerPos.Row+1; row++ {
					if col < 0 || col > 7 || row < 0 || row > 3 {
						continue
					}
					target := ecs.World.Create(component.GridPos, component.Damage, component.GridTarget, component.Transient)
					t := ecs.World.Entry(target)
					component.GridPos.Set(t, &component.GridPosComponentData{Col: col, Row: row})
					component.Damage.Set(t, &component.DamageData{Damage: dmg})
					component.Transient.Set(t, &component.TransientData{
						Start:            now,
						Duration:         200 * time.Millisecond,
						OnRemoveCallback: CreateAvalance,
					})
					// component.GridPos.Set(t,&component.GridPosComponentData{Col:col,Row: row})
				}
			}
			memory[CURRENT_STRATEGY] = "MOVE"
			memory[WARM_UP] = time.Now().Add(800 * time.Millisecond)
		}
	}
	if memory[CURRENT_STRATEGY] == "WARM_UP_MELEE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.YetiCooldown
			entityPunch := ecs.World.Create(component.Damage, component.OnHit, component.GridPos, component.Transient)
			entryPunch := ecs.World.Entry(entityPunch)
			component.Damage.Set(entryPunch, &component.DamageData{Damage: dmg})
			component.GridPos.Set(entryPunch, &component.GridPosComponentData{Row: gridPos.Row, Col: gridPos.Col - 1})
			component.Transient.Set(entryPunch, &component.TransientData{Start: time.Now(), Duration: 300 * time.Millisecond})
			component.OnHit.SetValue(entryPunch, YetiOnPunchHit)
			// punchAnim := core.NewMovableImage(assets.Fist, core.NewMovableImageParams())
			// punchAnim.ScaleParam = &core.ScaleParam{Sx: -1, Sy: 1}
			// punchAnim.
			memory[CURRENT_STRATEGY] = "COOLDOWN_MELEE"
			memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
		}
	}
	if memory[CURRENT_STRATEGY] == "COOLDOWN_MELEE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "WAIT"
			component.Sprite.Get(entity).Image = assets.Yeti
			memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
		}
	}
}

func CreateAvalance(ecs *ecs.ECS, entry *donburi.Entry) {
	gridPos := component.GridPos.Get(entry)
	scrX, scrY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
	startY := scrY - 300
	dmg := component.Damage.Get(entry).Damage
	anim := core.NewMovableImage(
		assets.Icicle,
		core.NewMovableImageParams().WithMoveParam(core.MoveParam{
			Sx: scrX,
			Sy: startY,
		}).WithScale(&core.ScaleParam{Sx: 1, Sy: -1}),
	)
	anim.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{Tx: scrX, Ty: scrY, Speed: 5}))
	fxentity := ecs.World.Create(component.Fx)
	fx := ecs.World.Entry(fxentity)
	component.Fx.Set(fx, &component.FxData{
		Animation: anim,
	})
	anim.Done = func() {
		ecs.World.Remove(fxentity)
		jj := ecs.World.Create(component.GridPos, component.Damage, component.Elements, component.Transient, component.OnHit)
		gridDmg := ecs.World.Entry(jj)
		component.GridPos.Set(gridDmg, gridPos)
		component.Damage.Set(gridDmg, &component.DamageData{
			Damage: dmg,
		})
		component.Elements.SetValue(gridDmg, component.WATER)
		component.Transient.Set(gridDmg, &component.TransientData{
			Start:    time.Now(),
			Duration: 100 * time.Millisecond,
		})
	}
}
