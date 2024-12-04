package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi/ecs"
)

type BambooLanceCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func NewBambooLanceCaster() *BambooLanceCaster {
	return &BambooLanceCaster{Cost: 100, nextCooldown: time.Now(), Damage: 150, CoolDown: 6 * time.Second, OnHit: SingleHitProjectile}
}
func (l *BambooLanceCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nHit target in the last column for %d damage.\nCooldown %.1fs", l.Cost/100, l.Damage, l.CoolDown.Seconds())
}
func (l *BambooLanceCaster) GetName() string {
	return "BambooLance"
}
func (l *BambooLanceCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *BambooLanceCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *BambooLanceCaster) SetModifier(e *loadout.CasterModifierData) {
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
func (f *BambooLanceCaster) GetElement() component.Elemental {
	if f.ModEntry != nil {
		return f.ModEntry.Element
	}
	return component.WOOD
}
func (l *BambooLanceCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		now := time.Now()
		for i := 0; i < 4; i++ {
			jj := ecs.World.Create(component.GridPos, component.Damage, component.OnHit, component.Transient, component.Elements)
			bambooLance := ecs.World.Entry(jj)
			gridPos := &component.GridPosComponentData{Row: i, Col: 7}
			component.GridPos.Set(bambooLance, gridPos)
			component.Damage.Set(bambooLance, &component.DamageData{Damage: l.GetDamage()})
			component.OnHit.SetValue(bambooLance, l.OnHit)
			component.Transient.Set(bambooLance, &component.TransientData{
				Start:    now,
				Duration: 300 * time.Millisecond,
			})
			component.Elements.SetValue(bambooLance, component.WOOD)
			fxEntity := ecs.World.Create(component.Fx, component.Transient)
			bambooLanceFx := ecs.World.Entry(fxEntity)
			sx, sy := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
			sx -= 50
			sy -= 100
			bounds := assets.BambooLance.Bounds()
			width := bounds.Dx()
			gg := core.NewMovableImageParams().WithMoveParam(
				core.MoveParam{Sx: sx + float64(width), Sy: sy},
			)
			gg.WithScale(&core.ScaleParam{Sx: -1, Sy: 1})
			fx := core.NewMovableImage(assets.BambooLance, gg)
			component.Fx.Set(bambooLanceFx, &component.FxData{Animation: fx})
			component.Transient.Set(bambooLanceFx, &component.TransientData{
				Start:    now,
				Duration: 300 * time.Millisecond,
			})
		}
	}
}

func (l *BambooLanceCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *BambooLanceCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *BambooLanceCaster) GetIcon() *ebiten.Image {
	return assets.BambooLanceIcon
}
func (l *BambooLanceCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *BambooLanceCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
