package system

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type damageSystem struct {
	DamagableQuery      *donburi.Query
	isGameOver          bool
	DamageEventConsumer DamageEventConsumer
	DamagingQuery       *donburi.Query
}
type DamageEventConsumer interface {
	OnCombatClear()
	OnGameOver()
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

func AddDoubleDamageFx(ecs *ecs.ECS, damagedEntity donburi.Entity) {
	entry := ecs.World.Entry(damagedEntity)
	screenPos := component.ScreenPos.Get(entry)
	scrX, scrY := screenPos.X, screenPos.Y
	if scrX == 0 && scrY == 0 {
		gridPos := component.GridPos.Get(entry)
		scrX, scrY = assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
		scrY -= 50
	}
	hitfx := core.NewMovableImage(assets.DoubleDamage,
		core.NewMovableImageParams().WithMoveParam(
			core.MoveParam{
				Sx: scrX - 50,
				Sy: scrY - 50,
			},
		),
	)
	entityFx := ecs.World.Create(component.Fx, component.Transient)
	entryFx := ecs.World.Entry(entityFx)
	component.Fx.Set(entryFx, &component.FxData{hitfx})
	component.Transient.Set(entryFx, &component.TransientData{
		Start:    time.Now(),
		Duration: 300 * time.Millisecond,
	})
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
	damageableList := []*donburi.Entry{}
	s.DamagableQuery.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := component.GridPos.Get(e)
		damageableList = append(damageableList, e)
		// health := mycomponent.Health.Get(e)
		// fmt.Println(e.Entity(), gridPos, gridPos.Row, gridPos.Col, health.Name)
		if gridPos.Row >= 0 && gridPos.Row <= 3 && gridPos.Col >= 0 && gridPos.Row <= 7 {
			gridMap[gridPos.Row][gridPos.Col] = e
		}

	})
	damagingEntries := []*donburi.Entry{}
	s.DamagingQuery.Each(ecs.World, func(e *donburi.Entry) {
		damagingEntries = append(damagingEntries, e)
	})
	for _, e := range damagingEntries {
		gridPos := component.GridPos.Get(e)
		if gridMap[gridPos.Row][gridPos.Col] != nil {
			damageableEntity := gridMap[gridPos.Row][gridPos.Col]
			// damage := mycomponent.Damage.Get(e).Damage
			onhit := component.OnHit.GetValue(e)
			invisTime := component.Health.Get(damageableEntity).InvisTime
			isZero := invisTime.IsZero()
			before := (!isZero && invisTime.Before(time.Now()))
			Health := component.Health.Get(damageableEntity)
			if isZero || before {
				f := e.HasComponent(component.Elements)
				if f && Health.Element != component.NEUTRAL {
					jj := component.Elements.GetValue(e)
					if component.IsDoubleDamage(jj, Health.Element) {
						dmg := component.Damage.Get(e).Damage
						Health.HP -= dmg
						AddDoubleDamageFx(ecs, damageableEntity.Entity())
						// onhit(ecs, e, damageableEntity)
					}
				}
				if onhit != nil {
					onhit(ecs, e, damageableEntity)
				}

				attack.AtkSfxQueue.QueueSFX(assets.ImpactFx)
				AddHitAnim(ecs, damageableEntity.Entity())
				if damageableEntity.HasComponent(component.Health) && component.Health.Get(damageableEntity).OnTakeDamage != nil {
					damageParam := component.DamageDetail{}
					component.Health.Get(damageableEntity).OnTakeDamage(ecs, damageableEntity, damageParam)
				}
			}

			// mycomponent.Health.Get(damageableEntity).HP -= damage
		}
	}
	for _, damageableEntity := range damageableList {
		if damageableEntity.HasComponent(component.Health) && component.Health.Get(damageableEntity).HP <= 0 {
			gridPos := component.GridPos.Get(damageableEntity)
			playerEnt, _ := archetype.PlayerTag.First(ecs.World)
			if playerEnt == damageableEntity && !s.isGameOver {
				s.isGameOver = true
				//TODO: game over screen
				//trigger game over here
				stgClrDim := assets.GameOver.Bounds()
				movableImg := core.NewMovableImage(assets.GameOver,
					core.NewMovableImageParams().WithMoveParam(core.MoveParam{
						Sx:    float64(-stgClrDim.Dx()),
						Sy:    float64(300 + stgClrDim.Dy()/2),
						Speed: 10}))
				movableImg.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{
					Tx:    float64(600 - stgClrDim.Dx()/2 - 60),
					Ty:    float64(300 + stgClrDim.Dy()/2),
					Speed: 10,
				}))
				movableImg.Done = func() {
					s.isGameOver = false
					PlayerAttackSystem.State = GameOverState

				}
				//turn off attack system
				PlayerAttackSystem.State = DoNothingState
				//attach the stageclear to fx system
				stgDone := ecs.World.Create(component.Anouncement)
				component.Anouncement.Set(ecs.World.Entry(stgDone), &component.FxData{
					Animation: movableImg,
				})
			} else {
				// trigger on destroy if it has any
				if damageableEntity.HasComponent(component.OnDestroy) {
					component.OnDestroy.GetValue(damageableEntity)(ecs, damageableEntity)
				}
				attack.AtkSfxQueue.QueueSFX(assets.ExplosionFx)
				// destroy anim
				scrPos := component.ScreenPos.GetValue(damageableEntity)
				if gridPos.Row < 0 || gridPos.Col < 0 {
					return
				}
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
					if s.DamageEventConsumer != nil {
						s.DamageEventConsumer.OnCombatClear()
					}
					stgClrDim := assets.StageClear.Bounds()
					movableImg := core.NewMovableImage(assets.StageClear,
						core.NewMovableImageParams().WithMoveParam(core.MoveParam{
							Sx:    float64(-stgClrDim.Dx()),
							Sy:    float64(300 + stgClrDim.Dy()/2),
							Speed: 10}))
					movableImg.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{
						Tx:    float64(600 - stgClrDim.Dx()/2 - 60),
						Ty:    float64(300 + stgClrDim.Dy()/2),
						Speed: 10,
					}))
					movableImg.Done = func() {
						PlayerAttackSystem.State = CombatClearState
					}
					//turn off attack system
					PlayerAttackSystem.State = DoNothingState
					//attach the stageclear to fx system
					stgDone := ecs.World.Create(component.Anouncement)
					component.Anouncement.Set(ecs.World.Entry(stgDone), &component.FxData{
						Animation: movableImg,
					})

				}
			}

		}
	}
}
