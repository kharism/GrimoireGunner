package system

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type gridPosStruct struct{}

func UpdateBurnSystem(ecs *ecs.ECS) {
	// burner := []*donburi.Entry{}
	burnerMap := map[component.GridPosComponentData]*donburi.Entry{}
	component.Burner.Each(ecs.World, func(e *donburi.Entry) {
		if e.HasComponent(component.GridPos) {
			gridPos := component.GridPos.GetValue(e)
			burnerMap[gridPos] = e
		}
	})
	damageble := []*donburi.Entry{}
	component.Health.Each(ecs.World, func(e *donburi.Entry) {
		damageble = append(damageble, e)
	})
	for _, val := range damageble {
		gridPos := component.GridPos.GetValue(val)
		if burnerMap[gridPos] != nil {
			burnerData := component.Burner.Get(burnerMap[gridPos])
			if !val.HasComponent(component.Burned) {
				val.AddComponent(component.Burned)

				burnedData := &component.BurnedData{
					NextBurn:   time.Now(),
					BurnDamage: burnerData.Damage,
					BurnCount:  0,
				}
				component.Burned.Set(val, burnedData)
			}
			burnedData := component.Burned.Get(val)
			if time.Now().After(burnedData.NextBurn) {
				burnedData.NextBurn = time.Now().Add(400 * time.Millisecond)
				// inflict burn damage here
				burnDamageTile := ecs.World.Create(component.GridPos, component.Elements, component.Damage, component.OnHit, component.Transient)
				entry := ecs.World.Entry(burnDamageTile)
				component.Elements.SetValue(entry, burnerData.Element)
				component.Transient.Set(entry, &component.TransientData{Start: time.Now(), Duration: 300 * time.Millisecond})
				component.GridPos.Set(entry, &component.GridPosComponentData{Col: gridPos.Col, Row: gridPos.Row})
				component.Damage.Set(entry, &component.DamageData{Damage: burnedData.BurnDamage})
				component.OnHit.SetValue(entry, OnTakeBurnDamage(ecs, entry, val))
			}

		}
	}
}

func OnTakeBurnDamage(ecs2 *ecs.ECS, projectile, receiver *donburi.Entry) component.OnAtkHit {
	burnedData := component.Burned.Get(receiver)
	return func(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
		burnedData.BurnCount += 1

		attack.SingleHitProjectile(ecs, projectile, receiver)
		if burnedData.BurnCount >= 10 {
			receiver.RemoveComponent(component.Burned)
		}
	}

}
