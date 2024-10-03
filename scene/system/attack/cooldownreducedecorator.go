package attack

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi/ecs"
)

type CooldownReduceCastDecor struct {
	// ensource loadout.ENSetGetter
	caster loadout.Caster
}

func DecorateWithCooldownReduce(caster loadout.Caster) loadout.Caster {
	if cc, ok := caster.(loadout.ModifierGetSetter); ok {
		// newModifier := ecs.World.Create(component.CasterModifier, component.PostAtkModifier)
		ll := cc.GetModifierEntry()
		if ll == nil {
			ll = &loadout.CasterModifierData{}

		}
		ll.CooldownModifer -= time.Second
		cc.SetModifier(ll)
		// entry := ecs.World.Entry(newModifier)
		// component.PostAtkModifier.SetValue(entry, DoubleCast(caster, entry))
		// cc.SetModifier(entry)
		return &CooldownReduceCastDecor{caster: caster}
	} else {
		return nil
	}
}
func (l *CooldownReduceCastDecor) GetModifierEntry() *loadout.CasterModifierData {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		return cc.GetModifierEntry()
	}
	return nil
}
func (l *CooldownReduceCastDecor) SetModifier(e *loadout.CasterModifierData) {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		cc.SetModifier(e)
	}
}
func (h *CooldownReduceCastDecor) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	// h.ensource = ensource
	h.caster.Cast(ensource, ecs)
}
func (h *CooldownReduceCastDecor) GetCost() int {
	return h.caster.GetCost()
}
func (h *CooldownReduceCastDecor) GetIcon() *ebiten.Image {
	return h.caster.GetIcon()
}
func (h *CooldownReduceCastDecor) GetCooldown() time.Time {
	return h.caster.GetCooldown()
}
func (h *CooldownReduceCastDecor) GetCooldownDuration() time.Duration {
	return h.caster.GetCooldownDuration()
}
func (h *CooldownReduceCastDecor) GetDamage() int {
	return h.caster.GetDamage()
}

func (h *CooldownReduceCastDecor) ResetCooldown() {
	h.caster.ResetCooldown()
}

func (h *CooldownReduceCastDecor) GetDescription() string {
	return h.caster.GetDescription()
}
func (h *CooldownReduceCastDecor) GetName() string {
	return h.caster.GetName() + " -1s"
}
