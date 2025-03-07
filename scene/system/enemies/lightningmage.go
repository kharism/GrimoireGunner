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

func NewLightningMage(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.LightningMage)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 500, Name: "LightningMage"})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})

	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 0 //it actualy damage bonus
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: LightningMageRoutine, Memory: data})
}

func LightningMageRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
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
			if playerGrid == nil {
				return
			}
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
			component.Sprite.Get(entity).Image = assets.LightningMageWarmup
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
		}
	}
	if memory[CURRENT_STRATEGY] == "WARM_UP" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.LightningMage
			LightningAttack(ecs, entity)
			memory[CURRENT_STRATEGY] = "WAIT"
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
			memory[CUR_DMG] = memory[CUR_DMG].(int) + 10
		}
	}
}
func LightningAttack(ecs *ecs.ECS, entry *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entry).Memory
	bonusDmg := memory[CUR_DMG].(int)
	demonPos := component.GridPos.Get(entry)
	// demonScreenPosX, demonScreenPosY := assets.GridCoord2Screen(demonPos.Row, demonPos.Col)
	playerPos, _ := attack.GetPlayerGridPos(ecs)
	if playerPos == nil {
		return
	}
	now := time.Now()
	for i := demonPos.Col - 1; i >= 0; i-- {
		DAMAGE := 50 + bonusDmg
		hitbox := ecs.World.Create(
			component.Damage,
			component.GridPos,
			component.OnHit,
			component.Transient,
			component.Elements,
		)
		lightningFx := ecs.World.Create(component.Fx, component.Transient)
		lightningFxEntry := ecs.World.Entry(lightningFx)

		entry = ecs.World.Entry(hitbox)
		component.Damage.Set(entry, &component.DamageData{Damage: DAMAGE})
		component.GridPos.Set(entry, &component.GridPosComponentData{Row: demonPos.Row, Col: i})
		component.OnHit.SetValue(entry, onReaperHit)
		component.Elements.SetValue(entry, component.ELEC)
		component.Transient.Set(entry, &component.TransientData{
			Start:    now,
			Duration: 200 * time.Millisecond,
		})
		scrX, scrY := assets.GridCoord2Screen(demonPos.Row, i)
		fxHeight := assets.LightningBolt.Bounds().Dy()
		fxWidth := assets.LightningBolt.Bounds().Dx()
		anim1 := core.NewMovableImage(assets.LightningBolt,
			core.NewMovableImageParams().
				WithMoveParam(core.MoveParam{Sx: scrX - (float64(assets.TileWidth) / 2) + float64(fxWidth), Sy: scrY - float64(fxHeight)}).
				WithScale(&core.ScaleParam{Sx: -1, Sy: 1}),
		)
		anim1.Done = func() {}
		component.Fx.Set(lightningFxEntry, &component.FxData{Animation: anim1})
		component.Transient.Set(lightningFxEntry, &component.TransientData{
			Start:    now,
			Duration: 200 * time.Millisecond,
		})

	}
}
