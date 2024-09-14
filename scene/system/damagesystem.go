package system

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/hanashi/core"
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
			component.Health,
			component.GridPos,
		),
	),
	DamagingQuery: donburi.NewQuery(
		filter.Contains(
			component.Damage,
			component.GridPos,
			component.OnHit,
		),
	),
}

// add hit animation to damagedEntity. Assuming the hit animation is 128x128
func AddHitAnim(ecs *ecs.ECS, damagedEntity donburi.Entity) {
	entry := ecs.World.Entry(damagedEntity)
	screenPos := component.ScreenPos.Get(entry)
	hitfx := assets.NewHitAnim(assets.SpriteParam{
		Modulo:  2,
		ScreenX: screenPos.X - 64,
		ScreenY: screenPos.Y - 100,
	})
	entityFx := ecs.World.Create(component.Fx)
	entryFx := ecs.World.Entry(entityFx)
	component.Fx.Set(entryFx, &component.FxData{hitfx})
	hitfx.Done = func() {
		ecs.World.Remove(entityFx)
	}
}
func (s *damageSystem) Update(ecs *ecs.ECS) {
	gridMap := [4][8]*donburi.Entry{}
	s.DamagableQuery.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := component.GridPos.Get(e)
		// health := mycomponent.Health.Get(e)
		// fmt.Println(e.Entity(), gridPos, gridPos.Row, gridPos.Col, health.Name)
		gridMap[gridPos.Row][gridPos.Col] = e
	})

	s.DamagingQuery.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := component.GridPos.Get(e)
		if gridMap[gridPos.Row][gridPos.Col] != nil {
			damageableEntity := gridMap[gridPos.Row][gridPos.Col]
			// damage := mycomponent.Damage.Get(e).Damage
			onhit := component.OnHit.GetValue(e)
			invisTime := component.Health.Get(damageableEntity).InvisTime
			isZero := invisTime.IsZero()
			before := (!isZero && invisTime.Before(time.Now()))
			if isZero || before {
				onhit(ecs, e, damageableEntity)
				AddHitAnim(ecs, damageableEntity.Entity())
				if component.Health.Get(damageableEntity).OnTakeDamage != nil {
					damageParam := component.DamageDetail{}
					component.Health.Get(damageableEntity).OnTakeDamage(ecs, damageableEntity, damageParam)
				}
			}

			if component.Health.Get(damageableEntity).HP <= 0 {
				playerEnt, _ := archetype.PlayerTag.First(ecs.World)
				if playerEnt == damageableEntity {
					//TODO: game over screen
					//trigger game over here
				} else {
					// destroy anim
					scrPos := component.ScreenPos.GetValue(damageableEntity)
					gridMap[gridPos.Row][gridPos.Col] = nil
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
					component.Fx.Set(entry, &component.FxData{explosionAnim})
					// check for other enemies
					enemyCount := 0
					component.EnemyTag.Each(ecs.World, func(e *donburi.Entry) {
						enemyCount += 1
					})
					if enemyCount == 0 {
						// TODO: end of battle
						// show stage clear
						stgClrDim := assets.StageClear.Bounds()
						movableImg := core.NewMovableImage(assets.StageClear,
							core.NewMovableImageParams().WithMoveParam(core.MoveParam{
								Sx:    float64(-stgClrDim.Dx()),
								Sy:    float64(300 + stgClrDim.Dy()/2),
								Speed: 10}))
						movableImg.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{
							Tx:    float64(600 - stgClrDim.Dx()/2 - 60),
							Ty:    float64(300 + stgClrDim.Dy()/2),
							Speed: 5,
						}))
						movableImg.Done = func() {

						}
						//turn off attack system
						PlayerAttackSystem.State = CombatClearState
						//attach the stageclear to fx system
						stgDone := ecs.World.Create(component.Fx)
						component.Fx.Set(ecs.World.Entry(stgDone), &component.FxData{
							Animation: movableImg,
						})

					}
				}

			}
			// mycomponent.Health.Get(damageableEntity).HP -= damage
		}
	})

}
