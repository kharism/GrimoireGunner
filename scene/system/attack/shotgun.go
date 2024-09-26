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
	ModEntry     *loadout.CasterModifierData
}

func (l *ShotgunCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *ShotgunCaster) SetModifier(e *loadout.CasterModifierData) {
	l.ModEntry = e
}
func NewShotgunCaster() *ShotgunCaster {
	return &ShotgunCaster{Cost: 100, nextCooldown: time.Now(), Damage: 50, CoolDown: 1 * time.Second}
}
func (l *ShotgunCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Damage 1 target on front and its behind immediately.\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.GetCooldownDuration().Seconds())
}
func (l *ShotgunCaster) GetName() string {
	return "Shotgun"
}
func (l *ShotgunCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *ShotgunCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *ShotgunCaster) GetIcon() *ebiten.Image {
	return assets.ShotgunIcon
}
func (l *ShotgunCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *ShotgunCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
func (l *ShotgunCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
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
			component.Damage.Set(grid1Entry, &component.DamageData{Damage: l.GetDamage()})
			component.OnHit.SetValue(grid1Entry, SingleHitProjectile)
			if targetGridPos.Col < 7 {
				grid1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Transient)
				grid1Entry := ecs.World.Entry(grid1)
				targetGridPos := component.GridPos.Get(closestTarget)
				component.GridPos.Set(grid1Entry, &component.GridPosComponentData{Col: targetGridPos.Col + 1, Row: targetGridPos.Row})
				component.Damage.Set(grid1Entry, &component.DamageData{Damage: l.GetDamage()})
				component.OnHit.SetValue(grid1Entry, SingleHitProjectile)
				component.Transient.Set(grid1Entry, &component.TransientData{Duration: 1 * time.Second, Start: time.Now()})
			}
		}
		if l.ModEntry != nil {
			// l := component.PostAtkModifier.GetValue(l.ModEntry)
			if l.ModEntry.PostAtk != nil {
				l.ModEntry.PostAtk(ecs, ensource)
			}
		}
	}
}
