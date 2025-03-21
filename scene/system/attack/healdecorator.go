package attack

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi/ecs"
)

type HealDecor struct {
	caster loadout.Caster
}

func DecorateWithHeal(caster loadout.Caster) loadout.Caster {
	if cc, ok := caster.(loadout.ModifierGetSetter); ok {
		var mod *loadout.CasterModifierData
		if cc.GetModifierEntry() != nil {
			mod = cc.GetModifierEntry()
		} else {
			mod = &loadout.CasterModifierData{}
		}
		if mod.PostAtk != nil {
			mod.PostAtk = func(ecs *ecs.ECS, ensource loadout.ENSetGetter) {
				mod.PostAtk(ecs, ensource)
				AddHeal(ecs, ensource)
			}
		} else {
			mod.PostAtk = AddHeal
		}
		cc.SetModifier(mod)
		return &HealDecor{caster}
	} else {
		return nil
	}

}

// this is a test
func AddHeal(ecs *ecs.ECS, ensource loadout.ENSetGetter) {
	gridPos, playerEnt := GetPlayerGridPos(ecs)
	healthComp := component.Health.Get(playerEnt)
	AtkSfxQueue.QueueSFX(assets.HealsFx)
	if healthComp.HP+5 < healthComp.MaxHP {
		healthComp.HP += 5
	} else {
		healthComp.HP = healthComp.MaxHP
	}

	fxEntity := ecs.World.Create(component.Fx, component.Transient)
	fx := ecs.World.Entry(fxEntity)

	x, y := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
	x -= 50
	y -= 100
	anim := core.NewMovableImage(assets.HealFx, core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: x, Sy: y}))
	component.Fx.Set(fx, &component.FxData{Animation: anim})
	component.Transient.Set(fx, &component.TransientData{Start: time.Now(), Duration: 500 * time.Millisecond})
}
func (h *HealDecor) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	h.caster.Cast(ensource, ecs)
}
func (l *HealDecor) GetElement() component.Elemental {
	if vv, ok := l.caster.(loadout.ElementalCaster); ok {
		return vv.GetElement()
	}
	return component.NEUTRAL
}
func (l *HealDecor) GetModifierEntry() *loadout.CasterModifierData {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		return cc.GetModifierEntry()
	}
	return nil
}
func (l *HealDecor) SetModifier(e *loadout.CasterModifierData) {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		cc.SetModifier(e)
	}
}
func (h *HealDecor) GetCost() int {
	return h.caster.GetCost()
}
func (h *HealDecor) GetIcon() *ebiten.Image {
	return h.caster.GetIcon()
}
func (h *HealDecor) GetCooldown() time.Time {
	return h.caster.GetCooldown()
}
func (h *HealDecor) GetCooldownDuration() time.Duration {
	return h.caster.GetCooldownDuration()
}
func (h *HealDecor) GetDamage() int {
	return h.caster.GetDamage()
}

func (h *HealDecor) ResetCooldown() {
	h.caster.ResetCooldown()
}

func (h *HealDecor) GetDescription() string {
	return h.caster.GetDescription() + "\nHeal 5HP on cast"
}
func (h *HealDecor) GetName() string {
	return h.caster.GetName() + " +Heal"
}
