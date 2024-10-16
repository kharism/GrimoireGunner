package attack

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi/ecs"
)

// add 30atk and add 5s cooldown
type TankingCastDecor struct {
	// ensource loadout.ENSetGetter
	caster loadout.Caster
}

func DecorateWithTank(caster loadout.Caster) loadout.Caster {
	if cc, ok := caster.(loadout.ModifierGetSetter); ok {
		// newModifier := ecs.World.Create(component.CasterModifier, component.PostAtkModifier)
		ll := cc.GetModifierEntry()
		if ll == nil {
			ll = &loadout.CasterModifierData{}

		}
		// ll.PostAtk = DoubleCast(caster)
		ll.CooldownModifer += 5 * time.Second
		ll.DamageModifier += 30
		cc.SetModifier(ll)
		// entry := ecs.World.Entry(newModifier)
		// component.PostAtkModifier.SetValue(entry, DoubleCast(caster, entry))
		// cc.SetModifier(entry)
		return &TankingCastDecor{caster: caster}
	} else {
		return nil
	}
}

func (l *TankingCastDecor) GetModifierEntry() *loadout.CasterModifierData {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		return cc.GetModifierEntry()
	}
	return nil
}
func (l *TankingCastDecor) SetModifier(e *loadout.CasterModifierData) {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		cc.SetModifier(e)
	}
}
func (h *TankingCastDecor) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	// h.ensource = ensource
	h.caster.Cast(ensource, ecs)
}
func (h *TankingCastDecor) GetCost() int {
	return h.caster.GetCost()
}
func (h *TankingCastDecor) GetIcon() *ebiten.Image {
	return h.caster.GetIcon()
}
func (h *TankingCastDecor) GetCooldown() time.Time {
	return h.caster.GetCooldown()
}
func (h *TankingCastDecor) GetCooldownDuration() time.Duration {
	return h.caster.GetCooldownDuration()
}
func (h *TankingCastDecor) GetDamage() int {
	return h.caster.GetDamage()
}

func (h *TankingCastDecor) ResetCooldown() {
	h.caster.ResetCooldown()
}

func (h *TankingCastDecor) GetDescription() string {
	return h.caster.GetDescription()
}
func (h *TankingCastDecor) GetName() string {
	return h.caster.GetName() + " +T"
}
