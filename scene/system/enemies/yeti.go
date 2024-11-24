package enemies

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func NewYeti(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Yeti)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 500, MaxHP: 500, Name: "Yeti"})

	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: YetiRoutine, Memory: data})
}

func YetiRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
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
				memory[CURRENT_STRATEGY] = "SUMMON_CONSTRUCT"
				// component.Sprite.Get(entity).Image = assets.YetiWarmup
				memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
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
			memory[CURRENT_STRATEGY] = "WAIT"
			memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
		}
	}
}

func YetiOnPunchHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	if receiver.HasComponent(archetype.ConstructTag) {
		newDamage := component.Health.Get(receiver).HP
		receiver.RemoveComponent(component.Health)
		receiver.AddComponent(component.Damage)
		receiver.AddComponent(archetype.ProjectileTag)
		receiver.AddComponent(component.Speed)
		receiver.AddComponent(component.OnHit)
		receiver.AddComponent(component.TargetLocation)
		screenPos := component.ScreenPos.Get(receiver)
		screenPos.X = 0
		screenPos.Y = 0

		component.Damage.Set(receiver, &component.DamageData{Damage: newDamage})
		component.OnHit.SetValue(receiver, attack.SingleHitProjectile)
		component.Speed.Set(receiver, &component.SpeedData{
			Vx: -5,
			Vy: 0,
		})
	} else {
		damage := component.Damage.Get(projectile).Damage
		component.Health.Get(receiver).HP -= damage
		ecs.World.Remove(projectile.Entity())
	}
}
