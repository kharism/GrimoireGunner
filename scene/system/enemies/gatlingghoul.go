package enemies

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func NewGatlingGhoul(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Gatlingghoul)
	entry := ecs.World.Entry(*entity)
	component.Health.Set(entry, &component.HealthData{HP: 800, Name: "Gatlinghoul"})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: GatlinghoulRoutine, Memory: data})
}

var scanOrder1 = [16][2]int{
	{0, 0}, {0, 1}, {0, 2}, {0, 3},
	{1, 3}, {1, 2}, {1, 1}, {1, 0},
	{2, 0}, {2, 1}, {2, 2}, {2, 3},
	{3, 3}, {3, 2}, {3, 1}, {3, 0}}

func GatlinghoulRoutine(ecs_ *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(2 * time.Second)
	}
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "WARMUP"
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
			component.Sprite.Get(entity).Image = assets.GatlingghoulAtk
		}
	}
	if memory[CURRENT_STRATEGY] == "WARMUP" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "ATTACK"
			memory["CurTarget"] = 0

			memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
			targetGrid := ecs_.World.Create(component.GridPos, component.GridTarget)
			entry := ecs_.World.Entry(targetGrid)
			memory["TargetEntity"] = targetGrid
			component.GridPos.Set(entry, &component.GridPosComponentData{Row: scanOrder1[0][0], Col: scanOrder1[0][1]})
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK" {
		if curIdx, ok := memory["CurTarget"].(int); ok && curIdx < 16 {
			if memory[WARM_UP].(time.Time).Before(time.Now()) {
				entityDust := ecs_.World.Create(component.GridPos, component.Damage, component.OnHit, component.Transient, component.Fx)
				entry := ecs_.World.Entry(entityDust)
				jj := scanOrder1[curIdx]
				component.GridPos.Set(entry, &component.GridPosComponentData{Row: jj[0], Col: jj[1]})
				component.Damage.Set(entry, &component.DamageData{Damage: 20})
				component.OnHit.SetValue(entry, GatlingGhoulOnAtkHit)
				component.Transient.Set(entry, &component.TransientData{Start: time.Now(), Duration: time.Second})
				screenX, screenY := assets.GridCoord2Screen(jj[0], jj[1])
				screenX -= 50
				screenY -= 50
				dustAnimParam := assets.SpriteParam{
					ScreenX: screenX,
					ScreenY: screenY,
					Modulo:  2,
					Done:    func() {},
				}
				kk := memory["TargetEntity"].(donburi.Entity)
				ecs_.World.Remove(kk)
				component.Fx.Set(entry, &component.FxData{Animation: assets.NewDustAnim(dustAnimParam)})
				memory[WARM_UP] = time.Now().Add(time.Second)
				if curIdx < 15 {
					memory["CurTarget"] = curIdx + 1
					targetGrid2 := ecs_.World.Create(component.GridPos, component.GridTarget)
					entry2 := ecs_.World.Entry(targetGrid2)
					memory["TargetEntity"] = targetGrid2

					component.GridPos.Set(entry2, &component.GridPosComponentData{Row: scanOrder1[curIdx+1][0], Col: scanOrder1[curIdx+1][1]})
				} else {
					memory[CURRENT_STRATEGY] = ""
					memory["CurTarget"] = 0
					component.Transient.Get(entry).OnRemoveCallback = func(ecs *ecs.ECS, L *donburi.Entry) {
						component.Sprite.Get(entity).Image = assets.Gatlingghoul
					}

				}

			}
		} else if curIdx == 16 {

		}
	}
}
func GatlingGhoulOnAtkHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	damageData := component.Damage.Get(projectile)
	component.Health.Get(receiver).HP -= damageData.Damage
	ecs.World.Remove(projectile.Entity())
}
