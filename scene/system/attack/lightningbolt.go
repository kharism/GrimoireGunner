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
	"github.com/yohamta/donburi/filter"
)

type LightnigAtkParam struct {
	StartRow, StartCol int
	//+1 to create lightning on right of starting point, and -1 to create on left
	Direction int
	Damage    int
	Element   component.Elemental
	Actor     *donburi.Entry
	OnHit     component.OnAtkHit
}

func LightningBoltOnHitfunc(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	damage := component.Damage.Get(projectile).Damage
	health := component.Health.Get(receiver)
	health.HP -= damage
	ecs.World.Remove(projectile.Entity())
}
func NewLigtningAttack(ecs *ecs.ECS, param LightnigAtkParam) {
	startCol := param.StartCol
	now := time.Now()
	for i := startCol; i >= 0 && i <= 7; i += param.Direction {
		entity := ecs.World.Create(
			component.Damage,
			component.GridPos,
			component.OnHit,
			component.Transient,
			component.Elements,
		)
		lightningFx := ecs.World.Create(component.Fx, component.Transient)
		lightningFxEntry := ecs.World.Entry(lightningFx)
		entry := ecs.World.Entry(entity)
		component.Damage.Set(entry, &component.DamageData{Damage: param.Damage})
		component.Elements.SetValue(entry, param.Element)
		component.GridPos.Set(entry, &component.GridPosComponentData{Row: param.StartRow, Col: i})
		scrX, scrY := assets.GridCoord2Screen(param.StartRow, i)
		fxHeight := assets.LightningBolt.Bounds().Dy()
		anim1 := core.NewMovableImage(assets.LightningBolt, core.
			NewMovableImageParams().
			WithMoveParam(core.MoveParam{Sx: scrX - (float64(assets.TileWidth) / 2), Sy: scrY - float64(fxHeight)}))
		anim1.Done = func() {}
		component.Fx.Set(lightningFxEntry, &component.FxData{Animation: anim1})
		component.Transient.Set(entry, &component.TransientData{
			Start:    now,
			Duration: 200 * time.Millisecond,
		})
		component.Transient.Set(lightningFxEntry, &component.TransientData{
			Start:    now,
			Duration: 200 * time.Millisecond,
		})
		component.OnHit.SetValue(entry, param.OnHit)
	}
}

// cost 3 EN and cast lightning bolt
type LightingBoltCaster struct {
	Cost         int
	nextCooldown time.Time
	CoolDown     time.Duration
	Damage       int
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func NewLightningBolCaster() *LightingBoltCaster {
	return &LightingBoltCaster{Cost: 300, Damage: 60, nextCooldown: time.Now(), CoolDown: 8 * time.Second, OnHit: LightningBoltOnHitfunc}
}
func (l *LightingBoltCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (f *LightingBoltCaster) GetElement() component.Elemental {
	if f.ModEntry != nil {
		return f.ModEntry.Element
	}
	return component.ELEC
}
func (l *LightingBoltCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}
func (l *LightingBoltCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		query := donburi.NewQuery(
			filter.Contains(
				archetype.PlayerTag,
			),
		)
		AtkSfxQueue.QueueSFX(assets.LightningFx)

		playerId, ok := query.First(ecs.World)
		if !ok {
			return
		}
		gridPos := component.GridPos.Get(playerId)
		param := LightnigAtkParam{
			StartRow:  gridPos.Row,
			StartCol:  gridPos.Col + 1,
			Direction: 1,
			Actor:     playerId,
			Damage:    l.GetDamage(),
			OnHit:     l.OnHit,
			Element:   l.GetElement(),
		}
		NewLigtningAttack(ecs, param)
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		if l.ModEntry != nil {
			// l := component.PostAtkModifier.GetValue(l.ModEntry)
			if l.ModEntry.PostAtk != nil {
				l.ModEntry.PostAtk(ecs, ensource)
			}
		}
	}
}
func (l *LightingBoltCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *LightingBoltCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *LightingBoltCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nShoots %d damage Lightning bolt instantly.\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.GetCooldownDuration().Seconds())
}
func (l *LightingBoltCaster) GetName() string {
	return "LightningBolt"
}
func (l *LightingBoltCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *LightingBoltCaster) GetIcon() *ebiten.Image {
	return assets.LightningIcon
}
func (l *LightingBoltCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *LightingBoltCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
