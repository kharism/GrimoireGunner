package system

import (
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type enemyAI struct {
	NPCQuery *donburi.Query
}

var EnemyAI = &enemyAI{
	NPCQuery: donburi.NewQuery(
		filter.Contains(
			component.EnemyRoutine,
		),
	),
}

func (e *enemyAI) Update(ecs *ecs.ECS) {
	enemyEntries := []*donburi.Entry{}
	e.NPCQuery.Each(ecs.World, func(entry *donburi.Entry) {
		// routine := component.EnemyRoutine.Get(entry)
		// routine.Routine(ecs, entry)
		enemyEntries = append(enemyEntries, entry)
	})
	for _, e := range enemyEntries {
		if !e.HasComponent(component.Root) {
			routine := component.EnemyRoutine.Get(e)
			routine.Routine(ecs, e)
		}

	}
}
