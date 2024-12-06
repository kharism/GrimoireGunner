package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type FlamethrowerCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	Cooldown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func NewFlamethrowerCaster() *FlamethrowerCaster {
	return &FlamethrowerCaster{Cost: 200, nextCooldown: time.Now(), Cooldown: 6 * time.Second, Damage: 80, OnHit: OnTowerHit}
}
func (l *FlamethrowerCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *FlamethrowerCaster) SetModifier(e *loadout.CasterModifierData) {
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
func (f *FlamethrowerCaster) GetDamage() int {
	if f.ModEntry != nil {
		// mod := component.CasterModifier.Get(f.ModEntry)
		return f.Damage + f.ModEntry.DamageModifier
	}
	return f.Damage
}
func (l *FlamethrowerCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nCreate flame for 3 grid long for %d damage.\nCooldown %.1fs", l.Cost/100, l.Damage, l.Cooldown.Seconds())
}
func (l *FlamethrowerCaster) GetName() string {
	return "Flamethrower"
}
func (f *FlamethrowerCaster) GetElement() component.Elemental {
	if f.ModEntry != nil {
		return f.ModEntry.Element
	}
	return component.FIRE
}
func (f *FlamethrowerCaster) GetCost() int {
	if f.ModEntry != nil {
		// mod := component.CasterModifier.Get(f.ModEntry)
		if f.Cost+f.ModEntry.CostModifier < 0 {
			return 0
		}
		return f.Cost + f.ModEntry.CostModifier
	}
	return f.Cost
}
func (l *FlamethrowerCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (f *FlamethrowerCaster) GetIcon() *ebiten.Image {
	return assets.FlamethrowerIcon
}
func (f *FlamethrowerCaster) GetCooldown() time.Time {
	return f.nextCooldown
}
func (f *FlamethrowerCaster) GetCooldownDuration() time.Duration {
	if f.ModEntry != nil {
		// mod := component.CasterModifier.Get(f.ModEntry)
		return f.Cooldown + f.ModEntry.CooldownModifer
	}
	return f.Cooldown
}
func (f *FlamethrowerCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	curEn := ensource.GetEn()
	if curEn >= f.Cost {
		AtkSfxQueue.QueueSFX(assets.HitscanFx)
		ensource.SetEn(curEn - f.Cost)
		f.nextCooldown = time.Now().Add(f.GetCooldownDuration())
		now := time.Now()
		gridPos, playerEnt := GetPlayerGridPos(ecs)
		playerEnt.AddComponent(component.Root)
		for i := 1; i <= 3; i++ {
			jj := ecs.World.Create(component.GridPos, component.Damage, component.OnHit, component.Elements, component.Transient)
			flameGrid := ecs.World.Entry(jj)
			component.GridPos.Set(flameGrid, &component.GridPosComponentData{Row: gridPos.Row, Col: gridPos.Col + i})
			component.Elements.SetValue(flameGrid, f.GetElement())
			component.OnHit.SetValue(flameGrid, SingleHitProjectile)
			component.Damage.Set(flameGrid, &component.DamageData{Damage: f.GetDamage()})
			component.Transient.Set(flameGrid, &component.TransientData{
				Start:            now,
				Duration:         300 * time.Millisecond,
				OnRemoveCallback: UnRoot,
			})
			jk := ecs.World.Create(component.Fx, component.Transient)
			flameFx := ecs.World.Entry(jk)
			scrX, scrY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col+i)
			fxHeight := assets.Flamehtrower.Bounds().Dy()
			fxWidth := assets.Flamehtrower.Bounds().Dx()
			fx := core.NewMovableImage(assets.Flamehtrower, core.NewMovableImageParams().
				WithScale(&core.ScaleParam{Sx: -1, Sy: 1}).
				WithMoveParam(core.MoveParam{Sx: scrX - (float64(assets.TileWidth) / 2) + float64(fxWidth), Sy: scrY - float64(fxHeight)}))
			component.Fx.Set(flameFx, &component.FxData{
				Animation: fx,
			})
			component.Transient.Set(flameFx, &component.TransientData{
				Start:    now,
				Duration: 300 * time.Millisecond,
			})

		}
	}
}
func UnRoot(ecs *ecs.ECS, entity *donburi.Entry) {
	_, playerEnt := GetPlayerGridPos(ecs)
	playerEnt.RemoveComponent(component.Root)
}
