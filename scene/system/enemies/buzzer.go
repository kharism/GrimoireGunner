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

func NewBuzzer(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Buzzer1)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 500, MaxHP: 500, Name: "Buzzer", Element: component.ELEC})

	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	data[OPTION_LIST] = []donburi.Entity{}
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: BuzzerRoutine, Memory: data})
}
func NewMiniBuzzer(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Buzzer1)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 200, MaxHP: 200, Name: "Buzzer", Element: component.ELEC})

	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	data[OPTION_LIST] = []donburi.Entity{}
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: BuzzerRoutine, Memory: data})
}
func BuzzerRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)

	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	gridPos := component.GridPos.Get(entity)
	playerPos, _ := attack.GetPlayerGridPos(ecs)
	if playerPos == nil {
		return
	}
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			scrPos := component.ScreenPos.Get(entity)
			newCol := playerPos.Col + 1
			newRow := playerPos.Row
			if validMove(ecs, newRow, newCol) {
				gridPos.Row = newRow
				gridPos.Col = newCol
				scrPos.Y = 0
				scrPos.X = 0
				memory[CURRENT_STRATEGY] = "ATTACK_MELEE"
				component.Sprite.Get(entity).Image = assets.Buzzer2
				memory[WARM_UP] = time.Now().Add(800 * time.Millisecond)

			} else if gridPos.Col == newCol && gridPos.Row == newRow {
				// buzzer is already in front of player
				memory[CURRENT_STRATEGY] = "ATTACK_MELEE"
				component.Sprite.Get(entity).Image = assets.Buzzer2
				memory[WARM_UP] = time.Now().Add(800 * time.Millisecond)
			}
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK_MELEE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			jj := ecs.World.Create(component.Damage, component.Elements, component.Fx, component.GridPos, component.Transient, component.OnHit)
			dmgTile := ecs.World.Entry(jj)
			component.GridPos.Set(dmgTile, &component.GridPosComponentData{Row: gridPos.Row, Col: gridPos.Col - 1})
			component.Transient.Set(dmgTile, &component.TransientData{
				Start:    time.Now(),
				Duration: 300 * time.Millisecond,
			})
			component.Damage.Set(dmgTile, &component.DamageData{Damage: dmg})
			component.Elements.SetValue(dmgTile, component.ELEC)
			component.OnHit.SetValue(dmgTile, TurnToElec)
			sX, sY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col-1)
			sX -= 50
			sY -= 50
			anim := core.NewMovableImage(assets.ElecSphere, core.NewMovableImageParams().WithMoveParam(core.MoveParam{
				Sx: sX,
				Sy: sY,
			}))
			component.Fx.Set(dmgTile, &component.FxData{
				Animation: anim,
			})
			memory[CURRENT_STRATEGY] = "WAIT"
			component.Sprite.Get(entity).Image = assets.Buzzer1
			memory[WARM_UP] = time.Now().Add(400 * time.Millisecond)
		}
	}
}

func TurnToElec(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	component.Health.Get(receiver).Element = component.ELEC
	if !receiver.HasComponent(component.Shader) {
		receiver.AddComponent(component.Shader)
	}
	shader := assets.Element2Shader(component.ELEC)
	if shader != nil {
		component.Shader.Set(receiver, shader)
	}
	attack.SingleHitProjectile(ecs, projectile, receiver)
}
func TurnToFire(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	component.Health.Get(receiver).Element = component.FIRE
	if !receiver.HasComponent(component.Shader) {
		receiver.AddComponent(component.Shader)
	}
	shader := assets.Element2Shader(component.FIRE)
	if shader != nil {
		component.Shader.Set(receiver, shader)
	}
	attack.SingleHitProjectile(ecs, projectile, receiver)
}
func TurnToWater(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	component.Health.Get(receiver).Element = component.WATER
	if !receiver.HasComponent(component.Shader) {
		receiver.AddComponent(component.Shader)
	}
	shader := assets.Element2Shader(component.WATER)
	if shader != nil {
		component.Shader.Set(receiver, shader)
	}
	attack.SingleHitProjectile(ecs, projectile, receiver)
}
func NewBlazeBuzzer(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Buzzer1)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	entry.AddComponent(component.Shader)
	component.Health.Set(entry, &component.HealthData{HP: 600, MaxHP: 600, Name: "BlazeBuzzer", Element: component.FIRE})
	component.Shader.Set(entry, assets.DakkaShader)
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50

	data[OPTION_LIST] = []donburi.Entity{}
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: BuzzerRoutine2, Memory: data})
}
func NewBlizzBuzzer(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Buzzer1)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	entry.AddComponent(component.Shader)
	component.Health.Set(entry, &component.HealthData{HP: 600, MaxHP: 600, Name: "BlazeBuzzer", Element: component.WATER})
	component.Shader.Set(entry, assets.IcyShader)
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50

	data[OPTION_LIST] = []donburi.Entity{}
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: BuzzerRoutine2, Memory: data})
}
func BuzzerRoutine2(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)

	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	gridPos := component.GridPos.Get(entity)
	playerPos, _ := attack.GetPlayerGridPos(ecs)
	if playerPos == nil {
		return
	}
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			scrPos := component.ScreenPos.Get(entity)
			newCol := playerPos.Col + 1
			newRow := playerPos.Row
			if validMove(ecs, newRow, newCol) {
				gridPos.Row = newRow
				gridPos.Col = newCol
				scrPos.Y = 0
				scrPos.X = 0
				memory[CURRENT_STRATEGY] = "ATTACK_MELEE"
				component.Sprite.Get(entity).Image = assets.Buzzer2
				memory[WARM_UP] = time.Now().Add(800 * time.Millisecond)

			} else if gridPos.Col == newCol && gridPos.Row == newRow {
				// buzzer is already in front of player
				memory[CURRENT_STRATEGY] = "ATTACK_MELEE"
				component.Sprite.Get(entity).Image = assets.Buzzer2
				memory[WARM_UP] = time.Now().Add(800 * time.Millisecond)
			}
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK_MELEE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			for i := -1; i < 2; i++ {
				if gridPos.Row+i < 0 || gridPos.Row+1 >= 4 {
					continue
				}
				jj := ecs.World.Create(component.Damage, component.Elements, component.Fx, component.GridPos, component.Transient, component.OnHit)
				dmgTile := ecs.World.Entry(jj)
				component.GridPos.Set(dmgTile, &component.GridPosComponentData{Row: gridPos.Row + i, Col: gridPos.Col - 1})
				component.Transient.Set(dmgTile, &component.TransientData{
					Start:    time.Now(),
					Duration: 300 * time.Millisecond,
				})
				component.Damage.Set(dmgTile, &component.DamageData{Damage: dmg})
				sX, sY := assets.GridCoord2Screen(gridPos.Row+i, gridPos.Col-1)
				switch component.Health.Get(entity).Element {
				case component.FIRE:
					component.OnHit.SetValue(dmgTile, TurnToFire)
					component.Elements.SetValue(dmgTile, component.FIRE)
					sX -= 50
					sY -= 100
					anim := core.NewMovableImage(assets.Flamehtrower, core.NewMovableImageParams().WithMoveParam(core.MoveParam{
						Sx: sX,
						Sy: sY,
					}))
					component.Fx.Set(dmgTile, &component.FxData{
						Animation: anim,
					})
				case component.ELEC:
					component.OnHit.SetValue(dmgTile, TurnToElec)
					component.Elements.SetValue(dmgTile, component.ELEC)
				case component.WATER:
					component.OnHit.SetValue(dmgTile, TurnToWater)
					component.Elements.SetValue(dmgTile, component.WATER)
					sX -= 50
					sY -= 100
					anim := core.NewMovableImage(assets.Icicle, core.NewMovableImageParams().WithMoveParam(core.MoveParam{
						Sx: sX,
						Sy: sY,
					}))
					component.Fx.Set(dmgTile, &component.FxData{
						Animation: anim,
					})
				}

			}

			memory[CURRENT_STRATEGY] = "WAIT"
			component.Sprite.Get(entity).Image = assets.Buzzer1
			memory[WARM_UP] = time.Now().Add(400 * time.Millisecond)
		}
	}
}
