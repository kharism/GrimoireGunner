package attack

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi/ecs"
)

type IcespikeCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func NewIcespikeCaster() *IcespikeCaster {
	return &IcespikeCaster{
		Cost: 200, nextCooldown: time.Now(), Damage: 130, CoolDown: 6 * time.Second, OnHit: SingleHitProjectile,
	}
}
func (l *IcespikeCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *IcespikeCaster) SetModifier(e *loadout.CasterModifierData) {
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
func (l *IcespikeCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Damage in 4 grid in front and 4 grids around it.\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.GetCooldownDuration().Seconds())
}
func (l *IcespikeCaster) GetName() string {
	return "IceSpike"
}
func (l *IcespikeCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *IcespikeCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *IcespikeCaster) GetIcon() *ebiten.Image {
	return assets.IcespikeIcon
}
func (l *IcespikeCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *IcespikeCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
func (l *IcespikeCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (f *IcespikeCaster) GetElement() component.Elemental {
	if f.ModEntry != nil {
		return f.ModEntry.Element
	}
	return component.WATER
}
func (l *IcespikeCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		gridPosPlayer, _ := GetPlayerGridPos(ecs)
		gridPos := &component.GridPosComponentData{Row: gridPosPlayer.Row, Col: gridPosPlayer.Col + 4}
		now := time.Now()
		for col := gridPos.Col - 1; col <= gridPos.Col+1; col++ {
			for row := gridPos.Row - 1; row <= gridPos.Row+1; row++ {
				dist := math.Abs(float64(gridPos.Col-col)) + math.Abs(float64(gridPos.Row-row))
				if col < 0 || col > 7 || row < 0 || row > 3 || dist > 1 {
					continue
				}
				jj := ecs.World.Create(component.GridPos, component.Damage, component.Elements, component.OnHit, component.Transient)
				entry := ecs.World.Entry(jj)
				component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
				component.Damage.Set(entry, &component.DamageData{Damage: l.GetDamage()})
				component.Elements.SetValue(entry, l.GetElement())
				component.OnHit.SetValue(entry, l.OnHit)
				component.Transient.Set(entry, &component.TransientData{
					Start:    now,
					Duration: 300 * time.Millisecond,
				})

				scrX, scrY := assets.GridCoord2Screen(row, col)
				scrX -= 50
				scrY -= 100
				spikeAnim := core.NewMovableImage(assets.Icicle, core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: scrX, Sy: scrY}))
				jk := ecs.World.Create(component.Fx, component.Transient)
				fxEntry := ecs.World.Entry(jk)
				component.Fx.Set(fxEntry, &component.FxData{Animation: spikeAnim})
				component.Transient.Set(fxEntry, &component.TransientData{
					Start:    now,
					Duration: 300 * time.Millisecond,
				})

			}
		}
	}
}
