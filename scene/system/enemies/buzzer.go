package enemies

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system"
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
func BuzzerRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)

	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	gridPos := component.GridPos.Get(entity)
	playerPos, _ := attack.GetPlayerGridPos(ecs)
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
	shader := system.Element2Shader(component.ELEC)
	if shader != nil {
		component.Shader.Set(receiver, shader)
	}
	attack.SingleHitProjectile(ecs, projectile, receiver)
}
