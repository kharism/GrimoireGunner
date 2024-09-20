package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi/ecs"
)

// cost 1 EN and cast hitscan piercing bullet
type ShotgunCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
}

func NewShotgunCaster() *ShotgunCaster {
	return &ShotgunCaster{Cost: 100, nextCooldown: time.Now(), Damage: 50, CoolDown: 1 * time.Second}
}
func (l *ShotgunCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Damage 1 target on front and its behind immediately.\nCooldown %.1fs", l.Cost/100, l.Damage, l.CoolDown.Seconds())
}
func (l *ShotgunCaster) GetName() string {
	return "Shotgun"
}
func (l *ShotgunCaster) GetDamage() int {
	return l.Damage
}
func (l *ShotgunCaster) GetCost() int {
	return l.Cost
}
func (l *ShotgunCaster) GetIcon() *ebiten.Image {
	return assets.ShotgunIcon
}
func (l *ShotgunCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *ShotgunCaster) GetCooldownDuration() time.Duration {
	return l.CoolDown
}

func (l *ShotgunCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.Cost {
		ensource.SetEn(en - l.Cost)
		l.nextCooldown = time.Now().Add(l.CoolDown)

		closestTarget := HitScanGetNearestTarget(ecs)
		if closestTarget != nil {
			grid1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
			grid1Entry := ecs.World.Entry(grid1)
			targetGridPos := component.GridPos.Get(closestTarget)
			component.GridPos.Set(grid1Entry, &component.GridPosComponentData{Col: targetGridPos.Col, Row: targetGridPos.Row})
			component.Damage.Set(grid1Entry, &component.DamageData{Damage: l.Damage})
			component.OnHit.SetValue(grid1Entry, SingleHitProjectile)
			if targetGridPos.Col < 7 {
				grid1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Transient)
				grid1Entry := ecs.World.Entry(grid1)
				targetGridPos := component.GridPos.Get(closestTarget)
				component.GridPos.Set(grid1Entry, &component.GridPosComponentData{Col: targetGridPos.Col + 1, Row: targetGridPos.Row})
				component.Damage.Set(grid1Entry, &component.DamageData{Damage: l.Damage})
				component.OnHit.SetValue(grid1Entry, SingleHitProjectile)
				component.Transient.Set(grid1Entry, &component.TransientData{Duration: 1 * time.Second, Start: time.Now()})
			}
		}
	}
}
