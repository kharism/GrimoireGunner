package system

import (
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
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
		// health := mycomponent.Health.Get(e)
		// fmt.Println(e.Entity(), gridPos, gridPos.Row, gridPos.Col, health.Name)
		gridMap[gridPos.Row][gridPos.Col] = e
	})

	s.DamagingQuery.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := mycomponent.GridPos.Get(e)
		if gridMap[gridPos.Row][gridPos.Col] != nil {
			damageableEntity := gridMap[gridPos.Row][gridPos.Col]
			// damage := mycomponent.Damage.Get(e).Damage
			onhit := mycomponent.OnHit.GetValue(e)
			onhit(ecs, e, damageableEntity)
			if component.Health.Get(damageableEntity).HP <= 0 {
				scrPos := mycomponent.ScreenPos.GetValue(damageableEntity)
				ecs.World.Remove(damageableEntity.Entity())
				explosionAnim := assets.NewExplosionAnim(assets.SpriteParam{
					ScreenX: scrPos.X - float64(assets.TileWidth)/2,
					ScreenY: scrPos.Y - 75,
					Modulo:  5,
				})
				entity := ecs.World.Create(component.Fx)
				entry := ecs.World.Entry(entity)
				explosionAnim.Done = func() {
					ecs.World.Remove(entity)
				}
				component.Fx.Set(entry, &explosionAnim)
			}
			// mycomponent.Health.Get(damageableEntity).HP -= damage
		}
	})

}
