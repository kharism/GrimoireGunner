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

func NewPyroEyes(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.PyroEyes)
	entry := ecs.World.Entry(*entity)
	component.Health.Set(entry, &component.HealthData{HP: 200, Name: "Pyro-Eyes", Element: component.FIRE})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: PyroEyesRoutine, Memory: data})
}

const ALREADY_FIRED = "already_fired" //
const IS_MOVING = "is_moving"

// this enemy will move up-down and if it's on the same row as player
// will attack
func PyroEyesRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
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
			component.Sprite.Get(entity).Image = assets.PyroEyesWarmup
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
		}
	}
	if memory[CURRENT_STRATEGY] == "WARM_UP" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.PyroEyes
			PyroEyesAttack(ecs, entity)
			memory[CURRENT_STRATEGY] = "WAIT"
			memory[WARM_UP] = time.Now().Add(600 * time.Millisecond)
			memory[CUR_DMG] = memory[CUR_DMG].(int) + 10
		}
	}
}
func PyroEyesAttack(ecs *ecs.ECS, entity *donburi.Entry) {
	rangeMax := 3
	// entry := ecs.World.Entry(entity)
	memory := component.EnemyRoutine.Get(entity).Memory
	damage := memory[CUR_DMG].(int)
	gridPos := component.GridPos.Get(entity)
	now := time.Now()
	for i := 1; i <= rangeMax; i++ {
		width := i - 1
		for j := gridPos.Row - width; j <= gridPos.Row+width; j++ {
			if j < 0 || j >= 4 {
				continue
			}
			damageTileEntity := ecs.World.Create(component.GridPos, component.Damage, component.OnHit, component.Transient, component.Elements)
			damageTile := ecs.World.Entry(damageTileEntity)
			component.GridPos.Set(damageTile, &component.GridPosComponentData{Row: j, Col: gridPos.Col - i})
			component.Damage.Set(damageTile, &component.DamageData{Damage: damage})
			component.Elements.SetValue(damageTile, component.FIRE)
			component.OnHit.SetValue(damageTile, attack.SingleHitProjectile)
			component.Transient.Set(damageTile, &component.TransientData{Start: now.Add(time.Duration(i*100) * time.Millisecond), Duration: 300 * time.Millisecond})
			scrX, scrY := assets.GridCoord2Screen(j, gridPos.Col-i)
			fxHeight := assets.Flamehtrower.Bounds().Dy()
			fx := core.NewMovableImage(assets.Flamehtrower, core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: scrX - (float64(assets.TileWidth) / 2), Sy: scrY - float64(fxHeight)}))
			fxEntity := ecs.World.Create(component.Transient, component.Fx)
			fxEntry := ecs.World.Entry(fxEntity)
			component.Fx.Set(fxEntry, &component.FxData{
				Animation: fx,
			})
			component.Transient.Set(fxEntry, &component.TransientData{Start: now.Add(time.Duration(i*100) * time.Millisecond), Duration: 100 * time.Millisecond})
		}
	}
}
