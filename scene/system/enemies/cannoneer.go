package enemies

import (
	"fmt"
	"math"
	"time"

	csg "github.com/kharism/golang-csg/core"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func NewCannoneer(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Cannoneer)
	entry := ecs.World.Entry(*entity)
	component.Health.Set(entry, &component.HealthData{HP: 200, Name: "Cannoneer"})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[IS_MOVING] = false
	data[LAST_FIRED] = nil
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: CannoneerRoutine, Memory: data})
}

var CURRENT_STRATEGY = "CurrMove"
var WARM_UP = "WarmUp"

// check whether there are obstacle on row-col grid
// for now it checks another character
func validMove(ecs *ecs.ECS, row, col int) bool {
	if row < 0 || row > 3 || col < 0 || col > 7 {
		return false
	}
	ObstacleExist := false
	var QueryHP = donburi.NewQuery(
		filter.Contains(
			component.Health,
			component.GridPos,
		),
	)
	QueryHP.Each(ecs.World, func(e *donburi.Entry) {
		pos := component.GridPos.Get(e)
		if pos.Col == col && pos.Row == row {
			ObstacleExist = true
		}
	})
	return !ObstacleExist
}

// this enemy will move up-down and if it's on the same row as player
// will attack
func CannoneerRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	player, _ := archetype.PlayerTag.First(ecs.World)
	playerGridPos := component.GridPos.Get(player)
	// playerScreenPos := component.ScreenPos.Get(player)
	entityGridPos := component.GridPos.Get(entity)
	entityScreenPos := component.ScreenPos.Get(entity)
	v := component.Speed.Get(entity)
	memory := component.EnemyRoutine.Get(entity).Memory
	sprite := component.Sprite.Get(entity)
	if entityScreenPos == nil {
		return
	}
	if v.Vx == 0 && v.Vy == 0 {
		memory[IS_MOVING] = false
		memory[CURRENT_STRATEGY] = ""
	}
	if memory[CURRENT_STRATEGY] == "" {
		if playerGridPos.Row != entityGridPos.Row {
			if moving, ok := memory[IS_MOVING].(bool); !ok || moving {
				return
			}
			memory[CURRENT_STRATEGY] = "MOVE"
			targetGridRowDirection := playerGridPos.Row - entityGridPos.Row
			distance := float64(targetGridRowDirection) / math.Abs(float64(targetGridRowDirection))
			targetRow := entityGridPos.Row + int(distance)
			if validMove(ecs, targetRow, entityGridPos.Col) {
				targetData := component.MoveTargetData{}
				targetData.Tx, targetData.Ty = assets.GridCoord2Screen(targetRow, entityGridPos.Col)
				pp, yy := assets.Coord2Grid(targetData.Tx, targetData.Ty)
				fmt.Println(pp, yy)
				component.TargetLocation.Set(entity, &targetData)
				Vy := (targetData.Ty - entityScreenPos.Y)
				Vx := 0.0
				speedVector := csg.NewVector(Vx, Vy, 0)
				speedVector = speedVector.Normalize().MultiplyScalar(1)
				component.Speed.Set(entity, &component.SpeedData{Vx: 0, Vy: speedVector.Y})
				memory[IS_MOVING] = true
			}

			// memory[CURRENT_STRATEGY] = ""
		} else if fired, ok := memory[ALREADY_FIRED]; !(!ok || fired.(bool)) {

			memory[CURRENT_STRATEGY] = "SHOOT"
			memory[IS_MOVING] = false
		}
	}
	if memory[CURRENT_STRATEGY] == "SHOOT" {
		memory[WARM_UP] = time.Now().Add(750 * time.Millisecond)
		memory[ALREADY_FIRED] = true
		sprite.Image = assets.CannoneerAtk

		memory[CURRENT_STRATEGY] = ""
	}
	if done, ok := memory[WARM_UP].(time.Time); ok && done.Before(time.Now()) {
		fmt.Println("Enemy attack")
		attack.GenerateMagibullet(ecs, entityGridPos.Row, entityGridPos.Col-1, -15)
		memory[ALREADY_FIRED] = false
		sprite.Image = assets.Cannoneer
		memory[CURRENT_STRATEGY] = ""
	}
	/**
	if playerGridPos.Row != entityGridPos.Row {
		if moving, ok := memory[IS_MOVING].(bool); !ok || moving {
			return
		}
		targetGridRowDirection := playerGridPos.Row - entityGridPos.Row
		distance := float64(targetGridRowDirection) / math.Abs(float64(targetGridRowDirection))
		targetRow := entityGridPos.Row + int(distance)
		targetData := component.MoveTargetData{}
		targetData.Tx, targetData.Ty = assets.GridCoord2Screen(targetRow, entityGridPos.Col)
		pp, yy := assets.Coord2Grid(targetData.Tx, targetData.Ty)
		fmt.Println(pp, yy)
		component.TargetLocation.Set(entity, &targetData)
		Vy := (targetData.Ty - entityScreenPos.Y)
		Vx := 0.0
		speedVector := csg.NewVector(Vx, Vy, 0)
		speedVector = speedVector.Normalize().MultiplyScalar(1)
		component.Speed.Set(entity, &component.SpeedData{Vx: 0, Vy: speedVector.Y})
		memory[IS_MOVING] = true
		// component.Speed.Set(entity, &component.SpeedData{Vx: 0, Vy: 1})
	} else {

		if fired, ok := memory[ALREADY_FIRED]; !ok || fired.(bool) {
			return
		}
		memory[ALREADY_FIRED] = true
		sprite.Image = assets.CannoneerAtk
		// timer := time.NewTimer(1500 * time.Millisecond)
		go func() {
			// <-timer.C
			time.Sleep(1 * time.Second)
			attack.GenerateMagibullet(ecs, entityGridPos.Row, entityGridPos.Col-1, -15)
			// timer.Stop()
			// timer = nil
			// time.Sleep(1 * time.Second)
			memory[ALREADY_FIRED] = false
			memory[IS_MOVING] = false
			sprite.Image = assets.Cannoneer
		}()

	}*/
}
