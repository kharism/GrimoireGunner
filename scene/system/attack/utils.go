package attack

import (
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func GetPlayerGridPos(ecs *ecs.ECS) (*component.GridPosComponentData, *donburi.Entry) {
	query := donburi.NewQuery(
		filter.Contains(
			archetype.PlayerTag,
		),
	)

	playerEntry, ok := query.First(ecs.World)
	if !ok {
		return nil, nil
	}
	return component.GridPos.Get(playerEntry), playerEntry
}

// get nearest damageable target
func GetNearestTarget(ecs *ecs.ECS) *donburi.Entry {
	_, playerId := GetPlayerGridPos(ecs)
	var closestTarget *donburi.Entry
	closestCol := 99
	query := donburi.NewQuery(filter.Contains(
		component.Health,
	))
	query.Each(ecs.World, func(e *donburi.Entry) {
		if e == playerId {
			return
		}
		gridPosE := component.GridPos.Get(e)
		if gridPosE.Col < closestCol {
			closestCol = gridPosE.Col
			closestTarget = e
		}
	})
	return closestTarget
}

// get the closest damageable target in the same column as the player
func HitScanGetNearestTarget(ecs *ecs.ECS) *donburi.Entry {
	gridPos, playerId := GetPlayerGridPos(ecs)
	var closestTarget *donburi.Entry
	closestCol := 99
	query := donburi.NewQuery(filter.Contains(
		component.Health,
	))
	query.Each(ecs.World, func(e *donburi.Entry) {
		if e == playerId {
			return
		}
		gridPosE := component.GridPos.Get(e)
		if gridPosE.Row == gridPos.Row {
			if gridPosE.Col < closestCol {
				closestCol = gridPosE.Col
				closestTarget = e
			}
		}
	})
	return closestTarget
}
