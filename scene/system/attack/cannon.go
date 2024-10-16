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

// cost 1 EN and cast hitscan bullet
type CannonCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func NewCannonCaster() *CannonCaster {
	return &CannonCaster{Cost: 100, nextCooldown: time.Now(), Damage: 80, CoolDown: 1 * time.Second, OnHit: SingleHitProjectile}
}
func (l *CannonCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Damage 1 target on front immediately.\nCooldown %.1fs", l.Cost/100, l.Damage, l.CoolDown.Seconds())
}
func (l *CannonCaster) GetName() string {
	return "Cannon"
}
func (l *CannonCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *CannonCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *CannonCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}

func (l *CannonCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		closestTarget := HitScanGetNearestTarget(ecs)
		if closestTarget != nil {
			grid1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Transient)
			grid1Entry := ecs.World.Entry(grid1)
			targetGridPos := component.GridPos.Get(closestTarget)
			component.GridPos.Set(grid1Entry, &component.GridPosComponentData{Col: targetGridPos.Col, Row: targetGridPos.Row})
			component.Damage.Set(grid1Entry, &component.DamageData{Damage: l.GetDamage()})
			component.Transient.Set(grid1Entry, &component.TransientData{Start: time.Now(), Duration: 100 * time.Millisecond})
			component.OnHit.SetValue(grid1Entry, l.OnHit)
		}
		if l.ModEntry != nil {
			// l := component.PostAtkModifier.GetValue(l.ModEntry)
			if l.ModEntry.PostAtk != nil {
				l.ModEntry.PostAtk(ecs, ensource)
			}
		}
	}
}

func (l *CannonCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *CannonCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *CannonCaster) GetIcon() *ebiten.Image {
	return assets.CannonIcon
}
func (l *CannonCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *CannonCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
