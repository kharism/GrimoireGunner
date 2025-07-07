package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi/ecs"
)

type SeedshotCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func NewSeedshotCaster() *SeedshotCaster {
	return &SeedshotCaster{
		Cost: 300, nextCooldown: time.Now(), Damage: 130, CoolDown: 6 * time.Second, OnHit: SingleHitProjectile,
	}
}
func (l *SeedshotCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *SeedshotCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	if l.GetElement() != component.NEUTRAL && e.Element == component.NEUTRAL {
		e.Element = l.GetElement()
	}
	l.ModEntry = e
}
func (l *SeedshotCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nSweeping attack 3 row wide for %d Damage.\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.GetCooldownDuration().Seconds())
}
func (l *SeedshotCaster) GetName() string {
	return "Seedshot"
}
func (l *SeedshotCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *SeedshotCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *SeedshotCaster) GetIcon() *ebiten.Image {
	return assets.SeedshotIcon
}
func (l *SeedshotCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *SeedshotCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
func (l *SeedshotCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (f *SeedshotCaster) GetElement() component.Elemental {
	if f.ModEntry != nil {
		return f.ModEntry.Element
	}
	return component.WOOD
}
func (l *SeedshotCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		gridPosPlayer, _ := GetPlayerGridPos(ecs)
		for row := -1; row <= 1; row++ {
			if gridPosPlayer.Row+row < 0 || gridPosPlayer.Row+row > 3 {
				continue
			}
			proj := archetype.NewProjectile(ecs.World, archetype.ProjectileParam{
				Vx:     25,
				Vy:     0,
				Col:    gridPosPlayer.Col + 1,
				Row:    gridPosPlayer.Row + row,
				Sprite: assets.Projectile1,
				Damage: l.GetDamage(),
				OnHit:  l.OnHit,
			})
			entry := ecs.World.Entry(*proj)
			entry.AddComponent(component.Elements)
			component.Elements.SetValue(entry, l.GetElement())
			entry.AddComponent(component.Shader)
			component.Shader.Set(entry, assets.WoodyShader)
		}
	}
}
