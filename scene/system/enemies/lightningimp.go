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

func NewLightningImp(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.LightningImp)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 300, MaxHP: 300, Name: "LightningImp", Element: component.ELEC})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})

	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = time.Now()
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 80

	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: LightingImpRoutine, Memory: data})
}

const TARG_GRID = "TARG_GRID"

func LightingImpRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	gridPos := component.GridPos.Get(entity)
	dmg := memory[CUR_DMG].(int)
	if memory[CURRENT_STRATEGY] == "" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "MOVE"
			component.Sprite.Set(entity, &component.SpriteData{
				Image: assets.LightningImp,
			})
			legalPos := []component.GridPosComponentData{
				component.GridPosComponentData{Row: 1, Col: 2},
				component.GridPosComponentData{Row: 2, Col: 2},
				component.GridPosComponentData{Row: 2, Col: 1},
				component.GridPosComponentData{Row: 1, Col: 1},
			}
			var targ_grid *component.GridPosComponentData
			if gridPos.Col == 1 && gridPos.Row == 1 {
				targ_grid = &component.GridPosComponentData{Row: 1, Col: 2}
			} else if gridPos.Col == 2 && gridPos.Row == 1 {
				targ_grid = &component.GridPosComponentData{Row: 2, Col: 2}
			} else if gridPos.Col == 2 && gridPos.Row == 2 {
				targ_grid = &component.GridPosComponentData{Row: 2, Col: 1}
			} else if gridPos.Col == 1 && gridPos.Row == 2 {
				targ_grid = &component.GridPosComponentData{Row: 1, Col: 1}
			}
			if targ_grid == nil {
				for _, ll := range legalPos {
					if validMove(ecs, ll.Row, ll.Col) {
						targ_grid = &component.GridPosComponentData{
							Row: ll.Row,
							Col: ll.Col,
						}
						break
					}
				}
			}
			memory[TARG_GRID] = targ_grid
			tX, tY := assets.GridCoord2Screen(targ_grid.Row, targ_grid.Col)
			// tX -= 50
			// tY -= 50
			component.TargetLocation.Set(entity, &component.MoveTargetData{
				Tx: tX,
				Ty: tY,
			})
			component.Speed.Set(entity, &component.SpeedData{V: 2})
			// memory[WARM_UP] = time.Now().Add(1 * time.Second)

		}
	}

	if memory[CURRENT_STRATEGY] == "MOVE" {
		speed := component.Speed.Get(entity)
		if speed.Vx == 0 && speed.Vy == 0 && speed.V == 0 {
			memory[CURRENT_STRATEGY] = "ATTACK"
			component.Sprite.Set(entity, &component.SpriteData{Image: assets.LightningImpWarmup})
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK" {
		memory[CURRENT_STRATEGY] = ""
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
		now := time.Now()
		for col := gridPos.Col - 1; col <= gridPos.Col+1; col++ {
			if col < 0 {
				continue
			}
			for row := gridPos.Row - 1; row <= gridPos.Row+1; row++ {
				if row < 0 || row > 3 {
					continue
				}
				if row == gridPos.Row && col == gridPos.Col {
					continue
				}
				targetGrid := ecs.World.Create(component.GridPos, component.Damage, component.Transient, component.GridTarget)
				targetGridEntry := ecs.World.Entry(targetGrid)
				component.Damage.Set(targetGridEntry, &component.DamageData{Damage: dmg})
				component.GridPos.Set(targetGridEntry, &component.GridPosComponentData{
					Row: row,
					Col: col,
				})
				component.Transient.Set(targetGridEntry, &component.TransientData{
					Start:            now,
					Duration:         300 * time.Millisecond,
					OnRemoveCallback: generateLightningBolt,
				})
			}
		}
	}
}

// remove target grid and create damage grid around lighting imp
func generateLightningBolt(ecs *ecs.ECS, entity *donburi.Entry) {
	gridPos := component.GridPos.Get(entity)
	dmgGridEntity := ecs.World.Create(
		component.GridPos,
		component.Transient,
		component.Damage,
		component.OnHit)
	dmgGrid := ecs.World.Entry(dmgGridEntity)
	component.Damage.Set(dmgGrid, component.Damage.Get(entity))
	component.GridPos.Set(dmgGrid, &component.GridPosComponentData{
		Col: gridPos.Col,
		Row: gridPos.Row,
	})
	now := time.Now()
	Sx, Sy := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
	Sx -= 50
	Sy -= 100
	component.OnHit.SetValue(dmgGrid, attack.SingleHitProjectile)
	component.Transient.Set(dmgGrid, &component.TransientData{Start: now, Duration: 250 * time.Millisecond})
	fx := core.NewMovableImage(assets.ElecSphere,
		core.NewMovableImageParams().WithMoveParam(core.MoveParam{
			Sx: Sx,
			Sy: Sy,
		}),
	)
	hh := ecs.World.Create(
		component.Transient,
		component.Fx,
	)
	hhE := ecs.World.Entry(hh)
	component.Transient.Set(hhE, &component.TransientData{
		Start: now, Duration: 250 * time.Millisecond,
	})
	component.Fx.Set(hhE, &component.FxData{Animation: fx})
}
