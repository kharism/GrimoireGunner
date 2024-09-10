package enemies

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func NewReaper(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Reaper)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 400, Name: "Reaper"})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: ReaperRoutine, Memory: data})
}

var ORIPOS = "OriPos"

func ReaperRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(3 * time.Second)
	}
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			// gets closer to the player
			playerQuery := donburi.NewQuery(filter.Contains(archetype.PlayerTag))
			playerEntry, _ := playerQuery.First(ecs.World)
			gridPos := component.GridPos.Get(playerEntry)
			var pos1 component.GridPosComponentData
			// counter := byte(0)
			for {
				// rowAdd := counter & 1
				// colAdd := counter & 0b10 >> 1
				pos1 = component.GridPosComponentData{Row: gridPos.Row + 1, Col: gridPos.Col + 1}
				if validMove(ecs, pos1.Row, pos1.Col) {
					break
				}
				pos1 = component.GridPosComponentData{Row: gridPos.Row - 1, Col: gridPos.Col + 1}
				if validMove(ecs, pos1.Row, pos1.Col) {
					break
				} else {
					memory[WARM_UP] = time.Now().Add(1 * time.Second)
					return
				}
			}

			memory[CURRENT_STRATEGY] = "WARMUP"
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
			component.Sprite.Get(entity).Image = assets.ReaperWarmup
			memory[ORIPOS] = component.GridPos.Get(entity)
			component.GridPos.Set(entity, &pos1)
			component.ScreenPos.Get(entity).X = 0
			component.ScreenPos.Get(entity).Y = 0
			now := time.Now()
			if pos1.Row > 0 {
				target1 := ecs.World.Create(component.GridPos, component.GridTarget, component.Transient)
				entry1 := ecs.World.Entry(target1)
				component.Transient.Set(entry1, &component.TransientData{Start: now, Duration: 745 * time.Millisecond})
				component.GridPos.Set(entry1, &component.GridPosComponentData{Col: pos1.Col - 1, Row: pos1.Row - 1})
			}

			target2 := ecs.World.Create(component.GridPos, component.GridTarget, component.Transient)
			entry2 := ecs.World.Entry(target2)
			component.Transient.Set(entry2, &component.TransientData{Start: now, Duration: 745 * time.Millisecond})
			component.GridPos.Set(entry2, &component.GridPosComponentData{Col: pos1.Col - 1, Row: pos1.Row})

			if pos1.Row < 3 {
				target3 := ecs.World.Create(component.GridPos, component.GridTarget, component.Transient)
				entry3 := ecs.World.Entry(target3)
				component.Transient.Set(entry3, &component.TransientData{Start: now, Duration: 745 * time.Millisecond})
				component.GridPos.Set(entry3, &component.GridPosComponentData{Col: pos1.Col - 1, Row: pos1.Row + 1})
			}

			memory[WARM_UP] = time.Now().Add(750 * time.Millisecond)
		}
	}
	if memory[CURRENT_STRATEGY] == "WARMUP" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			reaperGridPos := component.GridPos.Get(entity)
			reaperScreenX, reaperScreenY := assets.GridCoord2Screen(reaperGridPos.Row, reaperGridPos.Col)
			component.Sprite.Get(entity).Image = assets.ReaperCooldown
			var entry1 *donburi.Entry
			var entry2 *donburi.Entry
			var entry3 *donburi.Entry
			if reaperGridPos.Row > 0 {
				hitbox1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
				entry1 = ecs.World.Entry(hitbox1)
				component.Damage.Set(entry1, &component.DamageData{Damage: 20})
				component.GridPos.Set(entry1, &component.GridPosComponentData{Row: reaperGridPos.Row - 1, Col: reaperGridPos.Col - 1})
				component.OnHit.SetValue(entry1, onReaperHit)
			}

			hitbox2 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
			entry2 = ecs.World.Entry(hitbox2)
			component.Damage.Set(entry2, &component.DamageData{Damage: 20})
			component.GridPos.Set(entry2, &component.GridPosComponentData{Row: reaperGridPos.Row, Col: reaperGridPos.Col - 1})
			component.OnHit.SetValue(entry2, onReaperHit)

			if reaperGridPos.Row < 3 {
				hitbox3 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
				entry3 = ecs.World.Entry(hitbox3)
				component.Damage.Set(entry3, &component.DamageData{Damage: 20})
				component.GridPos.Set(entry3, &component.GridPosComponentData{Row: reaperGridPos.Row + 1, Col: reaperGridPos.Col - 1})
				component.OnHit.SetValue(entry3, onReaperHit)
			}

			fx := ecs.World.Create(component.Fx)
			fxEntry := ecs.World.Entry(fx)
			wideSlash := assets.NewWideSlashAtkAnim(assets.SpriteParam{
				ScreenX: reaperScreenX - float64(assets.TileWidth)*1.5,
				ScreenY: reaperScreenY - float64(assets.TileHeight)*2,
				Done: func() {
					if entry1 != nil {
						ecs.World.Remove(entry1.Entity())
					}

					ecs.World.Remove(entry2.Entity())
					if entry3 != nil {
						ecs.World.Remove(entry3.Entity())
					}
					ecs.World.Remove(fx)
				},
				Modulo: 3,
			})
			component.Fx.Set(fxEntry, &component.FxData{Animation: wideSlash})
			memory[CURRENT_STRATEGY] = "WAITTORETURN"
			memory[WARM_UP] = time.Now().Add(time.Second)
		}
	}
	if memory[CURRENT_STRATEGY] == "WAITTORETURN" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = ""
			opiPos := memory[ORIPOS].(*component.GridPosComponentData)
			component.GridPos.Set(entity, opiPos)
			scrPos := component.ScreenPos.Get(entity)
			scrPos.X = 0
			scrPos.Y = 0
			component.Sprite.Get(entity).Image = assets.Reaper
		}
	}
}

func onReaperHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	DmgComponent := component.Damage.Get(projectile)
	Health := component.Health.Get(receiver)
	Health.HP -= DmgComponent.Damage
	ecs.World.Remove(projectile.Entity())
}
