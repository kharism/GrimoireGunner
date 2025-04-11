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

func NewMoltenSlug(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.MoltenSlug)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 800, Name: "MoltenSlug", Element: component.FIRE})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data["CURRENT_PATTERN"] = 0
	data["FLAME_TOWER_COUNT"] = 0
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: MoltenSlugRoutine, Memory: data})
}

func MoltenSlugRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	health := component.Health.Get(entity)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.MoltenSlug
			if health.HP < 500 {
				memory[CURRENT_STRATEGY] = "MOVE2"
				memory["FLAME_TOWER_COUNT"] = 0
			} else {
				memory[CURRENT_STRATEGY] = "MOVE1"
			}

			memory[WARM_UP] = time.Now().Add(1 * time.Second)
		}
	}
	playerGrid, _ := attack.GetPlayerGridPos(ecs)
	if playerGrid == nil {
		return
	}
	demonPos := component.GridPos.Get(entity)
	tempRow := playerGrid.Row
	tempCol := 4
	// move to the front row and use flamethrower
	if memory[CURRENT_STRATEGY] == "MOVE1" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
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

			memory[CURRENT_STRATEGY] = "WARM_UP_FLAME_THROWER"
			component.Sprite.Get(entity).Image = assets.MoltenSlugWarmup
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
		}
	}
	if memory[CURRENT_STRATEGY] == "WARM_UP_FLAME_THROWER" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.MoltenSlugCooldown
			PyroEyesAttack(ecs, entity)
			memory[CURRENT_STRATEGY] = "WAIT"
			memory[WARM_UP] = time.Now().Add(600 * time.Millisecond)
			memory[CUR_DMG] = memory[CUR_DMG].(int) + 10
		}
	}
	// move to the back row and use spit out flame tower
	if memory[CURRENT_STRATEGY] == "MOVE2" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			tempCol := 7
			tempRow := playerGrid.Row
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
			memory[CURRENT_STRATEGY] = "SPAM_FLAME_TOWER"
			component.Sprite.Get(entity).Image = assets.MoltenSlugWarmup
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
		}
	}
	if memory[CURRENT_STRATEGY] == "SPAM_FLAME_TOWER" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.MoltenSlugCooldown
			atkPattern := memory["CURRENT_PATTERN"].(int)
			switch atkPattern {
			case 0:
				// horizontal pillars move from up to down
				now := time.Now()
				for i := range 4 {
					v1 := &createFlameTowerQueue{Row: i, Col: 0, damage: 30, time: now.Add(time.Duration(750*i) * time.Millisecond)}
					v2 := &createFlameTowerQueue{Row: i, Col: 1, damage: 30, time: now.Add(time.Duration(750*i) * time.Millisecond)}
					v3 := &createFlameTowerQueue{Row: i, Col: 2, damage: 30, time: now.Add(time.Duration(750*i) * time.Millisecond)}
					v4 := &createFlameTowerQueue{Row: i, Col: 3, damage: 30, time: now.Add(time.Duration(750*i) * time.Millisecond)}
					component.EventQueue.AddEvent(v1)
					component.EventQueue.AddEvent(v2)
					component.EventQueue.AddEvent(v3)
					component.EventQueue.AddEvent(v4)
				}
				memory[WARM_UP] = time.Now().Add(2500 * time.Millisecond)
				memory["CURRENT_PATTERN"] = 1
				memory["FLAME_TOWER_COUNT"] = memory["FLAME_TOWER_COUNT"].(int) + 1
			case 1:
				now := time.Now()
				for i := range 4 {
					v1 := &createFlameTowerQueue{Col: i, Row: 0, damage: 30, time: now.Add(time.Duration(750*i) * time.Millisecond)}
					v2 := &createFlameTowerQueue{Col: i, Row: 1, damage: 30, time: now.Add(time.Duration(750*i) * time.Millisecond)}
					v3 := &createFlameTowerQueue{Col: i, Row: 2, damage: 30, time: now.Add(time.Duration(750*i) * time.Millisecond)}
					v4 := &createFlameTowerQueue{Col: i, Row: 3, damage: 30, time: now.Add(time.Duration(750*i) * time.Millisecond)}
					component.EventQueue.AddEvent(v1)
					component.EventQueue.AddEvent(v2)
					component.EventQueue.AddEvent(v3)
					component.EventQueue.AddEvent(v4)
				}
				memory["FLAME_TOWER_COUNT"] = memory["FLAME_TOWER_COUNT"].(int) + 1
				if memory["FLAME_TOWER_COUNT"].(int) == 2 {
					memory[CURRENT_STRATEGY] = "MOVE1"
				} else {
					memory[WARM_UP] = time.Now().Add(2500 * time.Millisecond)
					memory["CURRENT_PATTERN"] = 0
				}
			}
		}
	}

}

type createFlameTowerQueue struct {
	Row    int
	Col    int
	time   time.Time
	damage int
}

func (c *createFlameTowerQueue) Execute(ecs *ecs.ECS) {
	f := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Fx, component.Transient)
	entry := ecs.World.Entry(f)
	component.Damage.Set(entry, &component.DamageData{
		Damage: c.damage,
	})
	targetScrX, targetScrY := assets.GridCoord2Screen(c.Row, c.Col)
	targetScrX -= 50
	targetScrY -= 100
	flameTower := core.NewMovableImage(assets.FlametowerRaw, core.NewMovableImageParams().
		WithMoveParam(core.MoveParam{Sx: targetScrX, Sy: targetScrY, Speed: 3}))
	component.Fx.Set(entry, &component.FxData{Animation: flameTower})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: c.Row, Col: c.Col})
	component.OnHit.SetValue(entry, attack.SingleHitProjectile)
	now := time.Now()
	component.Transient.Set(entry, &component.TransientData{
		Start:    now,
		Duration: 400 * time.Millisecond,
	})

}
func (c *createFlameTowerQueue) GetTime() time.Time {
	return c.time
}
func horizontalFlamePillars(ecs *ecs.ECS, entitySource *donburi.Entry) {

}
