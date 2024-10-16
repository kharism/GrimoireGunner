package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi/ecs"
)

type HealCaster struct {
	Cost             int
	Heal             int
	Charge           int
	nextCooldown     time.Time
	CooldownDuration time.Duration
	ModEntry         *loadout.CasterModifierData
	OnHit            component.OnAtkHit
}

func NewHealCaster() *HealCaster {
	return &HealCaster{Cost: 100, Heal: 100, Charge: 2, nextCooldown: time.Now(), CooldownDuration: 5 * time.Second}
}
func (l *HealCaster) GetDescription() string {
	return fmt.Sprintf("Heal %d. Has %d uses left", l.Heal, l.Charge)
}
func (l *HealCaster) GetName() string {
	return "Heal"
}
func (l *HealCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *HealCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}

func (l *HealCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *HealCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.CooldownDuration)

		gridPos, playerEnt := GetPlayerGridPos(ecs)
		healthComp := component.Health.Get(playerEnt)
		if healthComp.HP+l.Heal < healthComp.MaxHP {
			healthComp.HP += l.Heal
		} else {
			healthComp.HP = healthComp.MaxHP
		}

		fxEntity := ecs.World.Create(component.Fx, component.Transient)
		l.Charge -= 1
		fx := ecs.World.Entry(fxEntity)

		x, y := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
		x -= 50
		y -= 100
		anim := core.NewMovableImage(assets.HealFx, core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: x, Sy: y}))
		component.Fx.Set(fx, &component.FxData{Animation: anim})
		component.Transient.Set(fx, &component.TransientData{Start: time.Now(), Duration: 500 * time.Millisecond})
		if l.ModEntry != nil {
			if l.ModEntry.PostAtk != nil {
				l.ModEntry.PostAtk(ecs, ensource)
			}
		}
	}
}
func (l *HealCaster) GetCharge() int {
	return l.Charge
}
func (l *HealCaster) SetCharge(i int) {
	l.Charge = i
}
func (l *HealCaster) GetDamage() int {
	return 0
}
func (l *HealCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *HealCaster) GetIcon() *ebiten.Image {
	newImage := ebiten.NewImageFromImage(assets.HealIcon)
	bounds := newImage.Bounds()
	geom := ebiten.GeoM{}
	geom.Translate(float64(bounds.Dx()), 0)
	colorScale := ebiten.ColorScale{}
	colorScale.Scale(1, 1, 1, 1)
	text.Draw(newImage, fmt.Sprintf("%d", l.Charge), assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM:       geom,
			ColorScale: colorScale,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignEnd,
		},
	})
	return newImage
}
func (l *HealCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *HealCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CooldownDuration + l.ModEntry.CooldownModifer
	}
	return l.CooldownDuration
}
