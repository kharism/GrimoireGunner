package attack

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi/ecs"
)

type AtkBonusCastDecor struct {
	// ensource loadout.ENSetGetter
	caster loadout.Caster
}

func DecorateWithBonus10(caster loadout.Caster) loadout.Caster {
	if cc, ok := caster.(loadout.ModifierGetSetter); ok {
		// newModifier := ecs.World.Create(component.CasterModifier, component.PostAtkModifier)
		ll := cc.GetModifierEntry()
		if ll == nil {
			ll = &loadout.CasterModifierData{}

		}
		ll.DamageModifier += 10
		cc.SetModifier(ll)
		// entry := ecs.World.Entry(newModifier)
		// component.PostAtkModifier.SetValue(entry, DoubleCast(caster, entry))
		// cc.SetModifier(entry)
		return &AtkBonusCastDecor{caster: caster}
	} else {
		return nil
	}
}
func (l *AtkBonusCastDecor) GetElement() component.Elemental {
	if vv, ok := l.caster.(loadout.ElementalCaster); ok {
		return vv.GetElement()
	}
	return component.NEUTRAL
}
func (l *AtkBonusCastDecor) GetModifierEntry() *loadout.CasterModifierData {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		return cc.GetModifierEntry()
	}
	return nil
}
func (l *AtkBonusCastDecor) SetModifier(e *loadout.CasterModifierData) {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		cc.SetModifier(e)
	}
}
func (h *AtkBonusCastDecor) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	// h.ensource = ensource
	h.caster.Cast(ensource, ecs)
}
func (h *AtkBonusCastDecor) GetCost() int {
	return h.caster.GetCost()
}
func (h *AtkBonusCastDecor) GetIcon() *ebiten.Image {
	return h.caster.GetIcon()
}
func (h *AtkBonusCastDecor) GetCooldown() time.Time {
	return h.caster.GetCooldown()
}
func (h *AtkBonusCastDecor) GetCooldownDuration() time.Duration {
	return h.caster.GetCooldownDuration()
}
func (h *AtkBonusCastDecor) GetDamage() int {
	return h.caster.GetDamage()
}

func (h *AtkBonusCastDecor) ResetCooldown() {
	h.caster.ResetCooldown()
}

func (h *AtkBonusCastDecor) GetDescription() string {
	return h.caster.GetDescription()
}
func (h *AtkBonusCastDecor) GetName() string {
	return h.caster.GetName() + " +10"
}
