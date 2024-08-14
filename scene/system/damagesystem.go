package system

import (
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type damageSystem struct {
	DamagableQuery *donburi.Query

	DamagingQuery *donburi.Query
}

var DamageSystem = &damageSystem{
	DamagableQuery: donburi.NewQuery(
		filter.Contains(
			mycomponent.Health,
			mycomponent.GridPos,
		),
	),
	DamagingQuery: donburi.NewQuery(
		filter.Contains(
			mycomponent.Damage,
			mycomponent.GridPos,
			mycomponent.OnHit,
		),
	),
}

func (s *damageSystem) Update(ecs *ecs.ECS) {
	gridMap := [4][8]*donburi.Entry{}
	s.DamagableQuery.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := mycomponent.GridPos.Get(e)
		gridMap[gridPos.Row][gridPos.Col] = e
	})

	s.DamagingQuery.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := mycomponent.GridPos.Get(e)
		if gridMap[gridPos.Row][gridPos.Col] != nil {
			damageableEntity := gridMap[gridPos.Row][gridPos.Col]
			// damage := mycomponent.Damage.Get(e).Damage
			onhit := mycomponent.OnHit.GetValue(e)
			onhit(ecs, e, damageableEntity)
			// mycomponent.Health.Get(damageableEntity).HP -= damage
		}
	})

}
