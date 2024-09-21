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
	ModEntry     *donburi.Entry
}

func NewWideSwordCaster() *WideSwordCaster {
	return &WideSwordCaster{Cost: 200, Damage: 100, nextCooldown: time.Now()}
}
func (l *WideSwordCaster) GetModifierEntry() *donburi.Entry {
	return l.ModEntry
}
func (l *WideSwordCaster) SetModifier(e *donburi.Entry) {
	l.ModEntry = e
}
func (l *WideSwordCaster) GetDamage() int {
	if l.ModEntry != nil {
		mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + mod.DamageModifier
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
		if playerGridLoc.Row > 0 {
			hitbox1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
			entry1 = ecs.World.Entry(hitbox1)
			component.Damage.Set(entry1, &component.DamageData{Damage: l.GetDamage()})
			component.GridPos.Set(entry1, &component.GridPosComponentData{Row: playerGridLoc.Row - 1, Col: playerGridLoc.Col + 1})
			component.OnHit.SetValue(entry1, OnWideswordHit)
		}
		hitbox2 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
		entry2 = ecs.World.Entry(hitbox2)
		component.Damage.Set(entry2, &component.DamageData{Damage: l.GetDamage()})
		component.GridPos.Set(entry2, &component.GridPosComponentData{Row: playerGridLoc.Row, Col: playerGridLoc.Col + 1})
		component.OnHit.SetValue(entry2, OnWideswordHit)
		if playerGridLoc.Row < 3 {
			hitbox3 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
			entry3 = ecs.World.Entry(hitbox3)
			component.Damage.Set(entry3, &component.DamageData{Damage: l.GetDamage()})
			component.GridPos.Set(entry3, &component.GridPosComponentData{Row: playerGridLoc.Row + 1, Col: playerGridLoc.Col + 1})
			component.OnHit.SetValue(entry3, OnWideswordHit)
		}
		fxEntity := ecs.World.Create(component.Fx)
		fx := ecs.World.Entry(fxEntity)
		scrX, scrY := assets.GridCoord2Screen(playerGridLoc.Row-1, playerGridLoc.Col+1)
		scrX -= 50
		scrY -= 50
		wideSword := assets.NewWideSlashAtkAnim(assets.SpriteParam{
			ScreenX: scrX,
			ScreenY: scrY,
			Modulo:  5,
			Done: func() {
				if entry1 != nil {
					ecs.World.Remove(entry1.Entity())
				}
				if entry2 != nil {
					ecs.World.Remove(entry2.Entity())
				}
				if entry3 != nil {
					ecs.World.Remove(entry3.Entity())
				}
				ecs.World.Remove(fxEntity)
			},
		})
		wideSword.FlipHorizontal = true
		component.Fx.Set(fx, &component.FxData{Animation: wideSword})

	}
}

func OnWideswordHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	DmgComponent := component.Damage.Get(projectile)
	Health := component.Health.Get(receiver)
	Health.HP -= DmgComponent.Damage
	ecs.World.Remove(projectile.Entity())
}
func (l *WideSwordCaster) GetCost() int {
	if l.ModEntry != nil {
		mod := component.CasterModifier.Get(l.ModEntry)
		return l.Cost + mod.CostModifier
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
		mod := component.CasterModifier.Get(l.ModEntry)
		return mod.CooldownModifer
	}
	return 0
}
