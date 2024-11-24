package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type FistCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func NewFist() *FistCaster {
	return &FistCaster{Cost: 100, Damage: 100, CoolDown: 0, nextCooldown: time.Now(), OnHit: fistOnHit}
}
func (l *FistCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *FistCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}
func (l *FistCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Damage 1 target on front.\nIf it hits a construct, throw the construct.\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.GetCooldownDuration().Seconds())
}
func (l *FistCaster) GetName() string {
	return "Fist"
}
func (l *FistCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *FistCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *FistCaster) GetIcon() *ebiten.Image {
	return assets.FistIcon
}
func (l *FistCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *FistCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
func (l *FistCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *FistCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		AtkSfxQueue.QueueSFX(assets.ImpactFx)
		gridPos, _ := GetPlayerGridPos(ecs)
		dmgTileEntity := ecs.World.Create(component.Damage, component.GridPos, component.Transient, component.Fx, component.OnHit)
		dmgTile := ecs.World.Entry(dmgTileEntity)
		component.Damage.Set(dmgTile, &component.DamageData{Damage: l.Damage})
		component.GridPos.Set(dmgTile, &component.GridPosComponentData{
			Row: gridPos.Row,
			Col: gridPos.Col + 1,
		})
		now := time.Now()
		component.OnHit.SetValue(dmgTile, l.OnHit)
		component.Transient.Set(dmgTile, &component.TransientData{Start: now, Duration: 300 * time.Millisecond})
		fistBound := assets.Fist.Bounds()
		posX, posY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
		fxSprite := core.NewMovableImage(assets.Fist,
			core.NewMovableImageParams().WithMoveParam(core.MoveParam{
				Sx: posX - 20,
				Sy: posY - float64(fistBound.Dy()),
			}))
		component.Fx.Set(dmgTile, &component.FxData{
			Animation: fxSprite,
		})
		if l.ModEntry != nil {
			// l := component.PostAtkModifier.GetValue(l.ModEntry)
			if l.ModEntry.PostAtk != nil {
				l.ModEntry.PostAtk(ecs, ensource)
			}
		}
	}
}
func fistOnHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	if receiver.HasComponent(archetype.ConstructTag) {
		newDamage := component.Health.Get(receiver).HP
		receiver.RemoveComponent(component.Health)
		receiver.AddComponent(component.Damage)
		receiver.AddComponent(archetype.ProjectileTag)
		receiver.AddComponent(component.Speed)
		receiver.AddComponent(component.OnHit)
		receiver.AddComponent(component.TargetLocation)
		screenPos := component.ScreenPos.Get(receiver)
		screenPos.X = 0
		screenPos.Y = 0

		component.Damage.Set(receiver, &component.DamageData{Damage: newDamage})
		component.OnHit.SetValue(receiver, SingleHitProjectile)
		component.Speed.Set(receiver, &component.SpeedData{
			Vx: 5,
			Vy: 0,
		})
	} else {
		damage := component.Damage.Get(projectile).Damage
		component.Health.Get(receiver).HP -= damage
		ecs.World.Remove(projectile.Entity())
	}
}
