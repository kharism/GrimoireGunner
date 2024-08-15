package enemies

import (
	"time"

	csg "github.com/kharism/golang-csg/core"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func NewPyroEyes(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.PyroEyes)
	entry := ecs.World.Entry(*entity)
	component.Health.Set(entry, &component.HealthData{HP: 200, Name: "Pyro-Eyes"})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: PyroEyesRoutine})
}

// this enemy will move up-down and if it's on the same row as player
// will attack
func PyroEyesRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	player, _ := archetype.PlayerTag.First(ecs.World)
	playerGridPos := component.GridPos.Get(player)
	playerScreenPos := component.ScreenPos.Get(player)
	entityGridPos := component.GridPos.Get(entity)
	entityScreenPos := component.ScreenPos.Get(entity)
	if entityScreenPos == nil {
		return
	}
	if playerGridPos.Row != entityGridPos.Row {
		targetData := component.MoveTargetData{}
		targetData.Tx = entityScreenPos.X
		targetData.Ty = playerScreenPos.Y
		component.TargetLocation.Set(entity, &targetData)
		Vy := (targetData.Ty - entityScreenPos.Y)
		Vx := 0.0
		speedVector := csg.NewVector(Vx, Vy, 0)
		speedVector = speedVector.Normalize().MultiplyScalar(1)
		component.Speed.Set(entity, &component.SpeedData{Vx: 0, Vy: speedVector.Y})
		// component.Speed.Set(entity, &component.SpeedData{Vx: 0, Vy: 1})
	} else {
		timer := time.NewTimer(1500 * time.Millisecond)
		go func() {
			<-timer.C
			attack.GenerateMagibullet(ecs, entityGridPos.Row, entityGridPos.Col-1, -15)
		}()

	}
}
