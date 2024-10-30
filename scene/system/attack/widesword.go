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

type WideSwordCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func NewWideSwordCaster() *WideSwordCaster {
	return &WideSwordCaster{Cost: 200, Damage: 100, nextCooldown: time.Now(), OnHit: OnWideswordHit}
}
func (l *WideSwordCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *WideSwordCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}
func (l *WideSwordCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *WideSwordCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nHit 1 column in front, up, down for %d damage.\nNo cooldown", l.Cost/100, l.GetDamage())
}
func (l *WideSwordCaster) GetName() string {
	return "WideSword"
}
func (l *WideSwordCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())

		playerGridLoc, _ := GetPlayerGridPos(ecs)
		var entry1 *donburi.Entry
		var entry2 *donburi.Entry
		var entry3 *donburi.Entry
		now := time.Now()

		if playerGridLoc.Row > 0 {
			hitbox1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Transient)
			entry1 = ecs.World.Entry(hitbox1)
			component.Damage.Set(entry1, &component.DamageData{Damage: l.GetDamage()})
			component.GridPos.Set(entry1, &component.GridPosComponentData{Row: playerGridLoc.Row - 1, Col: playerGridLoc.Col + 1})
			component.OnHit.SetValue(entry1, l.OnHit)
			component.Transient.Set(entry1, &component.TransientData{Start: now, Duration: 100 * time.Millisecond})
		}
		hitbox2 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Transient)
		entry2 = ecs.World.Entry(hitbox2)
		component.Damage.Set(entry2, &component.DamageData{Damage: l.GetDamage()})
		component.GridPos.Set(entry2, &component.GridPosComponentData{Row: playerGridLoc.Row, Col: playerGridLoc.Col + 1})
		component.OnHit.SetValue(entry2, l.OnHit)
		component.Transient.Set(entry2, &component.TransientData{Start: now, Duration: 100 * time.Millisecond})
		if playerGridLoc.Row < 3 {
			hitbox3 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Transient)
			entry3 = ecs.World.Entry(hitbox3)
			component.Damage.Set(entry3, &component.DamageData{Damage: l.GetDamage()})
			component.GridPos.Set(entry3, &component.GridPosComponentData{Row: playerGridLoc.Row + 1, Col: playerGridLoc.Col + 1})
			component.OnHit.SetValue(entry3, l.OnHit)
			component.Transient.Set(entry3, &component.TransientData{Start: now, Duration: 100 * time.Millisecond})
		}
		fxEntity := ecs.World.Create(component.Fx)
		AtkSfxQueue.QueueSFX(assets.SlashFx)
		fx := ecs.World.Entry(fxEntity)
		scrX, scrY := assets.GridCoord2Screen(playerGridLoc.Row-1, playerGridLoc.Col+1)
		scrX -= 50
		scrY -= 50
		wideSword := assets.NewWideSlashAtkAnim(assets.SpriteParam{
			ScreenX: scrX,
			ScreenY: scrY,
			Modulo:  5,
			Done: func() {

				ecs.World.Remove(fxEntity)
			},
		})
		wideSword.FlipHorizontal = true
		component.Fx.Set(fx, &component.FxData{Animation: wideSword})
		if l.ModEntry != nil {
			// l := component.PostAtkModifier.GetValue(l.ModEntry)
			if l.ModEntry.PostAtk != nil {
				l.ModEntry.PostAtk(ecs, ensource)
			}
		}
	}
}
func (l *WideSwordCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func OnWideswordHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	DmgComponent := component.Damage.Get(projectile)
	Health := component.Health.Get(receiver)
	Health.HP -= DmgComponent.Damage
	ecs.World.Remove(projectile.Entity())
}
func (l *WideSwordCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *WideSwordCaster) GetIcon() *ebiten.Image {
	return assets.WideSwordIcon
}
func (l *WideSwordCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *WideSwordCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.ModEntry.CooldownModifer
	}
	return 0
}
