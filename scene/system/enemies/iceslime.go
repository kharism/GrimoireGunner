package enemies

import (
	"math/rand"
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func NewIceslime(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Slime)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	entry.AddComponent(component.Shader)
	component.Health.Set(entry, &component.HealthData{HP: 300, MaxHP: 300, Name: "IceSlime", Element: component.WATER})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	component.Shader.Set(entry, assets.IcyShader)

	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: IceSlimeRoutine, Memory: data})
}
func IceSlimeRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	gridPos := component.GridPos.Get(entity)
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "MOVE"
		}
	}
	if memory[CURRENT_STRATEGY] == "MOVE" {
		component.Sprite.Set(entity, &component.SpriteData{Image: assets.Slime})
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			moveCount := memory[MOVE_COUNT].(int)
			if moveCount < 4 {

				scrPos := component.ScreenPos.Get(entity)
			checkMove:
				for {
					rndMove := rand.Int() % 4
					switch rndMove {
					case 0:
						if validMove(ecs, gridPos.Row-1, gridPos.Col) {
							gridPos.Row -= 1
							scrPos.Y -= float64(assets.TileHeight)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					case 1:
						if validMove(ecs, gridPos.Row, gridPos.Col-1) && gridPos.Col-1 >= 4 {
							gridPos.Col -= 1
							scrPos.X -= float64(assets.TileWidth)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					case 2:
						if validMove(ecs, gridPos.Row+1, gridPos.Col) {
							gridPos.Row += 1
							scrPos.Y += float64(assets.TileHeight)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					case 3:
						if validMove(ecs, gridPos.Row, gridPos.Col+1) {
							gridPos.Col += 1
							scrPos.X += float64(assets.TileWidth)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					}
				}
			} else {
				memory[MOVE_COUNT] = 0
				memory[CURRENT_STRATEGY] = "ATTACK"
				component.Sprite.Set(entity, &component.SpriteData{Image: assets.Slime2})
				memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
			}
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			pair := map[int]map[int]bool{}
			for i := 0; i < 4; i++ {
				var col int
				var row int
				for true {
					col = rand.Int() % 4
					row = rand.Int() % 4
					if _, ok := pair[row][col]; ok {
						continue
					} else {
						if pair[row] == nil {
							pair[row] = map[int]bool{}
						}
						pair[row][col] = true
						break
					}
				}
				targetZone := ecs.World.Create(component.GridPos, component.GridTarget, component.Damage, component.Transient)
				targetZoneEntry := ecs.World.Entry(targetZone)
				component.GridPos.Set(targetZoneEntry, &component.GridPosComponentData{Col: col, Row: row})
				component.Damage.Set(targetZoneEntry, &component.DamageData{
					Damage: dmg,
				})
				component.Transient.Set(targetZoneEntry, &component.TransientData{
					Start:            time.Now(),
					Duration:         700 * time.Millisecond,
					OnRemoveCallback: CreateIcicle,
				})
				memory[CURRENT_STRATEGY] = "WAIT"
				memory[WARM_UP] = time.Now().Add(1100 * time.Millisecond)
			}
		}
	}
}
func CreateIcicle(ecs *ecs.ECS, entity *donburi.Entry) {
	gridPos := component.GridPos.Get(entity)
	scrX, scrY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
	scrX -= 50
	scrY -= 100
	dmg := component.Damage.Get(entity).Damage
	fxImg := core.NewMovableImage(assets.Icicle, core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: scrX, Sy: scrY}))
	icicleFx := ecs.World.Create(component.Fx, component.GridPos, component.Elements, component.Damage, component.OnHit, component.Transient)
	iciclefxEnt := ecs.World.Entry(icicleFx)
	component.Fx.Set(iciclefxEnt, &component.FxData{
		Animation: fxImg,
	})
	component.Elements.SetValue(iciclefxEnt, component.WATER)
	component.GridPos.Set(iciclefxEnt, gridPos)
	component.Transient.Set(iciclefxEnt, &component.TransientData{
		Start:    time.Now(),
		Duration: 400 * time.Millisecond,
	})
	component.Damage.Set(iciclefxEnt, &component.DamageData{Damage: dmg})
	component.OnHit.SetValue(iciclefxEnt, attack.SingleHitProjectile)
}
