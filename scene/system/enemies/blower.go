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

func NewBlower(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Blower1)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 500, MaxHP: 500, Name: "Buzzer", Element: component.WOOD})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	data[OPTION_LIST] = []donburi.Entity{}
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: BlowerRoutine, Memory: data})
}
func BlowerRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	// dmg := memory[CUR_DMG].(int)
	// gridPos := component.GridPos.Get(entity)

	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "BLOW_FORWARD"
		component.Sprite.Set(entity, &component.SpriteData{Image: assets.Blower2})
		memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
		blowScrPos := component.ScreenPos.Get(entity)
		scrFxXpos := blowScrPos.X - 50
		scrFxYpos := blowScrPos.Y - 100
		blowAnim := assets.NewBlowAnim(assets.SpriteParam{
			ScreenX: scrFxXpos,
			ScreenY: scrFxYpos,
			Modulo:  15,
			Done:    nil,
		})
		jj := ecs.World.Create(component.Fx)
		pp := ecs.World.Entry(jj)
		component.Fx.Set(pp, &component.FxData{Animation: blowAnim})
	}
	playerPos, _ := attack.GetPlayerGridPos(ecs)
	if playerPos == nil {
		return
	}
	if memory[CURRENT_STRATEGY] == "BLOW_FORWARD" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			if validMove(ecs, playerPos.Row, playerPos.Col-1) {
				playerPos.Col -= 1
			}
			memory[WARM_UP] = time.Now().Add(800 * time.Millisecond)
		}
	}
}
