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

type SporebombCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func NewSporeBombCaster() *SporebombCaster {
	return &SporebombCaster{
		Cost: 300, nextCooldown: time.Now(), Damage: 90, CoolDown: 6 * time.Second, OnHit: SingleHitProjectile,
	}
}

func (l *SporebombCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *SporebombCaster) SetModifier(e *loadout.CasterModifierData) {
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
func (l *SporebombCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Damage in 4 grid in front and 8 grids around it.\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.GetCooldownDuration().Seconds())
}
func (l *SporebombCaster) GetName() string {
	return "Sporebomb"
}
func (l *SporebombCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *SporebombCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *SporebombCaster) GetIcon() *ebiten.Image {
	return assets.SporebombIcon
}
func (l *SporebombCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *SporebombCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
func (l *SporebombCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (f *SporebombCaster) GetElement() component.Elemental {
	if f.ModEntry != nil {
		return f.ModEntry.Element
	}
	return component.WOOD
}
func (l *SporebombCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		gridPosPlayer, _ := GetPlayerGridPos(ecs)
		gridPos := &component.GridPosComponentData{Row: gridPosPlayer.Row, Col: gridPosPlayer.Col + 4}
		ents := ecs.World.CreateMany(9, component.GridPos, component.Damage, component.Elements, component.OnHit, component.Transient)
		fxEnts := ecs.World.CreateMany(9, component.Fx, component.Transient)
		entsIdx := 0
		now := time.Now()
		for col := gridPos.Col - 1; col <= gridPos.Col+1; col++ {
			for row := gridPos.Row - 1; row <= gridPos.Row+1; row++ {
				if col < 0 || col > 7 || row < 0 || row > 3 {
					continue
				}
				entry := ecs.World.Entry(ents[entsIdx])
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
				sporeAnim := assets.NewSporeAnim(assets.SpriteParam{
					ScreenX: scrX,
					ScreenY: scrY,
					Modulo:  5,
					Done:    func() {},
				})
				fxEntry := ecs.World.Entry(fxEnts[entsIdx])
				component.Fx.Set(fxEntry, &component.FxData{Animation: sporeAnim})
				component.Transient.Set(fxEntry, &component.TransientData{
					Start:    now,
					Duration: 300 * time.Millisecond,
				})
				entsIdx += 1
			}
		}
	}
}
