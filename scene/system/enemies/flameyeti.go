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
)

func NewFlameYeti(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Yeti)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	entry.AddComponent(component.Shader)
	component.Health.Set(entry, &component.HealthData{HP: 1000, MaxHP: 1000, Name: "Yeti", Element: component.FIRE})
	component.Shader.Set(entry, assets.DakkaShader)

	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[PUNCHCOUNT] = 0
	data[CUR_DMG] = 50
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: FlameYetiRoutine, Memory: data})
}

var PUNCHCOUNT = "PunchCount"

func FlameYetiRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory

	dmg := memory[CUR_DMG].(int)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	punchcount := memory[PUNCHCOUNT].(int)
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
	playerGrid, _ := attack.GetPlayerGridPos(ecs)
	if memory[CURRENT_STRATEGY] == "MOVE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			if playerGrid == nil {
				return
			}
			if punchcount%5 == 0 && punchcount > 0 {
				// throw mines at the back
				memory[CURRENT_STRATEGY] = "THROW_MINE"
				component.Sprite.Get(entity).Image = assets.YetiWarmup2

			} else {
				memory[CURRENT_STRATEGY] = "WARM_UP_MELEE"
				component.Sprite.Get(entity).Image = assets.YetiWarmup
			}
			demonPos := component.GridPos.Get(entity)
			tempRow := playerGrid.Row
			// default to melee
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

			memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
		}
	}
	if memory[CURRENT_STRATEGY] == "THROW_MINE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			sourceScrX, sourceSrcY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
			sourceScrX -= 50
			sourceSrcY -= 50
			col := playerGrid.Col //memory[MOVE_COUNT].(int)
			component.Sprite.Get(entity).Image = assets.YetiCooldown2
			memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
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
			memory[CURRENT_STRATEGY] = "WARM_UP_MELEE"
		}
	}
	if memory[CURRENT_STRATEGY] == "WARM_UP_MELEE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.YetiCooldown
			flyingPunch := archetype.NewProjectile(ecs.World, archetype.ProjectileParam{
				Vx:             -10,
				Vy:             0,
				Col:            gridPos.Col - 1,
				Row:            gridPos.Row,
				Sprite:         assets.Fist,
				FlipHorizontal: true,
				OnHit:          attack.SingleHitProjectile,
				Damage:         dmg,
			})
			flyingPunchEntry := ecs.World.Entry(*flyingPunch)
			flyingPunchEntry.AddComponent(component.Preremove)
			component.Preremove.SetValue(flyingPunchEntry, LeaveFlameTower)
			memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
			punchcount += 1
			memory[PUNCHCOUNT] = punchcount
			memory[CURRENT_STRATEGY] = "MOVE"
		}
	}
}
func LeaveFlameTower(ecs *ecs.ECS, Entry *donburi.Entry) {
	//ent := ecs.World.Entry(*entity)
	gridPos1 := component.GridPos.Get(Entry)
	flame := archetype.NewHazard(ecs.World, assets.FlametowerRaw)
	flameEnt := ecs.World.Entry(flame)
	flameEnt.RemoveComponent(component.EnemyRoutine)
	flameEnt.AddComponent(component.Transient)
	flameEnt.AddComponent(component.Burner)

	gridPos := component.GridPos.Get(flameEnt)
	gridPos.Col = gridPos1.Col + 1
	gridPos.Row = gridPos1.Row
	component.Burner.Set(flameEnt, &component.BurnerData{
		Damage:  20,
		Element: component.FIRE,
	})
	component.Transient.Set(flameEnt, &component.TransientData{
		Start:    time.Now(),
		Duration: 1500 * time.Millisecond,
	})

}
