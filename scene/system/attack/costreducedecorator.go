package attack

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi/ecs"
)

type CostReducerDecor struct {
	caster loadout.Caster
}

func DecorateWithCostReducer(caster loadout.Caster) loadout.Caster {
	if cc, ok := caster.(loadout.ModifierGetSetter); ok {
		// newModifier := ecs.World.Create(component.CasterModifier, component.PostAtkModifier)
		ll := cc.GetModifierEntry()
		if ll == nil {
			ll = &loadout.CasterModifierData{}

		}
		ll.CooldownModifer -= 100
		cc.SetModifier(ll)
		// entry := ecs.World.Entry(newModifier)
		// component.PostAtkModifier.SetValue(entry, DoubleCast(caster, entry))
		// cc.SetModifier(entry)
		return &CostReducerDecor{caster: caster}
	} else {
		return nil
	}
}
func (l *CostReducerDecor) GetModifierEntry() *loadout.CasterModifierData {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		return cc.GetModifierEntry()
	}
	return nil
}
func (l *CostReducerDecor) SetModifier(e *loadout.CasterModifierData) {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		cc.SetModifier(e)
	}
}
func (h *CostReducerDecor) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	// h.ensource = ensource
	h.caster.Cast(ensource, ecs)
}
func (h *CostReducerDecor) GetCost() int {
	return h.caster.GetCost()
}
func (h *CostReducerDecor) GetIcon() *ebiten.Image {
	return h.caster.GetIcon()
}
func (h *CostReducerDecor) GetCooldown() time.Time {
	return h.caster.GetCooldown()
}
func (h *CostReducerDecor) GetCooldownDuration() time.Duration {
	return h.caster.GetCooldownDuration()
}
func (h *CostReducerDecor) GetDamage() int {
	return h.caster.GetDamage()
}

func (h *CostReducerDecor) ResetCooldown() {
	h.caster.ResetCooldown()
}

func (h *CostReducerDecor) GetDescription() string {
	return h.caster.GetDescription()
}
func (h *CostReducerDecor) GetName() string {
	return h.caster.GetName() + " -1EN"
}
