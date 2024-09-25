package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type GatlingCaster struct {
	Cost             int
	Damage           int
	ShotAmount       int
	nextCooldown     time.Time
	CooldownDuration time.Duration
	ModEntry         *component.CasterModifierData
}

func NewGatlingCastor() *GatlingCaster {
	return &GatlingCaster{Cost: 200, Damage: 5, ShotAmount: 10, nextCooldown: time.Now(), CooldownDuration: 3 * time.Second}
}
func (l *GatlingCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *GatlingCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nShoots %d damage bullet %d times.\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.ShotAmount, l.CooldownDuration.Seconds())
}
func (l *GatlingCaster) GetName() string {
	return "Gatling"
}
func (l *GatlingCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		now := time.Now()
		for i := 0; i < l.ShotAmount; i++ {
			ev := &addGatlingShot{Damage: l.GetDamage(), Time: now.Add(time.Duration(100*i) * time.Millisecond)}
			component.EventQueue.Queue = append(component.EventQueue.Queue, ev)
		}
		if l.ModEntry != nil {
			// if l.ModEntry.HasComponent(component.PostAtkModifier) {
			// l := component.PostAtkModifier.GetValue(l.ModEntry)
			if l.ModEntry.PostAtk != nil {
				l.ModEntry.PostAtk(ecs, ensource)
			}

			// }
		}
	}

}
func (l *GatlingCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *GatlingCaster) GetModifierEntry() *component.CasterModifierData {
	return l.ModEntry
}
func (l *GatlingCaster) SetModifier(e *component.CasterModifierData) {
	l.ModEntry = e
}

type addGatlingShot struct {
	Damage int
	Time   time.Time
}

func (a *addGatlingShot) Execute(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.Contains(
			archetype.PlayerTag,
		),
	)

	playerEntry, ok := query.First(ecs.World)
	if !ok {
		return
	}
	// playerScrLoc := component.ScreenPos.GetValue(playerEntry)
	playerGridLoc := component.GridPos.GetValue(playerEntry)
	archetype.NewProjectile(ecs.World, archetype.ProjectileParam{
		Vx:     15,
		Vy:     0,
		Col:    playerGridLoc.Col + 1,
		Row:    playerGridLoc.Row,
		Sprite: assets.Projectile1,
		Damage: a.Damage,
		OnHit:  SingleHitProjectile,
	})
}
func (a *addGatlingShot) GetTime() time.Time {
	return a.Time
}
func (l *GatlingCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *GatlingCaster) GetIcon() *ebiten.Image {
	return assets.GatlingIcon
}
func (l *GatlingCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *GatlingCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CooldownDuration + l.ModEntry.CooldownModifer
	}
	return l.CooldownDuration
}
