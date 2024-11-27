package system

import (
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type kamikazeSystem struct {
	KamikazeQuery *donburi.Query
	HealthQuery   *donburi.Query
}

var Kamikaze = &kamikazeSystem{
	KamikazeQuery: donburi.NewQuery(
		filter.Contains(
			component.KamikazeTag,
			component.GridPos,
			component.Damage,
		),
	),
	HealthQuery: donburi.NewQuery(
		filter.Contains(
			component.Health,
			component.GridPos,
		),
	),
}

func (s *kamikazeSystem) Update(ecs *ecs.ECS) {
	tileMapKamikaze := [4][8]donburi.Entity{
		[8]donburi.Entity{0, 0, 0, 0, 0, 0, 0, 0},
		[8]donburi.Entity{0, 0, 0, 0, 0, 0, 0, 0},
		[8]donburi.Entity{0, 0, 0, 0, 0, 0, 0, 0},
		[8]donburi.Entity{0, 0, 0, 0, 0, 0, 0, 0},
	}
	s.KamikazeQuery.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := component.GridPos.Get(e)
		tileMapKamikaze[gridPos.Row][gridPos.Col] = e.Entity()
	})
	removeList := []*donburi.Entry{}
	s.HealthQuery.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := component.GridPos.Get(e)
		if tileMapKamikaze[gridPos.Row][gridPos.Col] != 0 {
			if tileMapKamikaze[gridPos.Row][gridPos.Col] == e.Entity() {
				return
			} else {
				ll := ecs.World.Entry(tileMapKamikaze[gridPos.Row][gridPos.Col])
				removeList = append(removeList, ll)
			}
		}

	})
	for _, entry := range removeList {
		gridPos := component.GridPos.Get(entry)
		dmg := component.Damage.Get(entry)
		dmgTile := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
		dmgTileEntry := ecs.World.Entry(dmgTile)
		component.GridPos.Set(dmgTileEntry, gridPos)
		component.Damage.Set(dmgTileEntry, dmg)
		component.OnHit.SetValue(dmgTileEntry, attack.SingleHitProjectile)
		ecs.World.Remove(entry.Entity())
	}
}
