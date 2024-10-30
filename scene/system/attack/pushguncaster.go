package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// cost 1 EN and cast hitscan piercing bullet
type PushgunCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func (l *PushgunCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *PushgunCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}
func NewPushgunCaster() *PushgunCaster {
	return &PushgunCaster{Cost: 100, nextCooldown: time.Now(), Damage: 50, CoolDown: 3 * time.Second, OnHit: PushbackOnHit}
}
func PushbackOnHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	damage := component.Damage.Get(projectile).Damage
	component.Health.Get(receiver).HP -= damage
	if receiver.HasComponent(component.GridPos) {
		receiverPos := component.GridPos.Get(receiver)
		if receiverPos.Col < 7 && validMove(ecs, receiverPos.Row, receiverPos.Col+1) {
			receiverPos.Col += 1
			scrPos := component.ScreenPos.Get(receiver)
			scrPos.X = 0
			scrPos.Y = 0
		}
	}
	ecs.World.Remove(projectile.Entity())
}
func (l *PushgunCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Damage 1 target on front and Push it behind.\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.GetCooldownDuration().Seconds())
}
func (l *PushgunCaster) GetName() string {
	return "Pushgun"
}
func (l *PushgunCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *PushgunCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *PushgunCaster) GetIcon() *ebiten.Image {
	return assets.PushgunIcon
}
func (l *PushgunCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *PushgunCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
func (l *PushgunCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *PushgunCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.Cost {
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		closestTarget := HitScanGetNearestTarget(ecs)
		AtkSfxQueue.QueueSFX(assets.HitscanFx)
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
