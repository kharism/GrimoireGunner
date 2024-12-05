package attack

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type PushbackDecorator struct {
	caster loadout.Caster
}

func DecorateWithPushbackDecorator(caster loadout.Caster) loadout.Caster {
	if cc, ok := caster.(loadout.ModifierGetSetter); ok {
		// newModifier := ecs.World.Create(component.CasterModifier, component.PostAtkModifier)
		ll := cc.GetModifierEntry()
		if ll == nil {
			ll = &loadout.CasterModifierData{}

		}
		ll.OnHit = PushBackNoDmg
		cc.SetModifier(ll)
		// entry := ecs.World.Entry(newModifier)
		// component.PostAtkModifier.SetValue(entry, DoubleCast(caster, entry))
		// cc.SetModifier(entry)
		return &PushbackDecorator{caster: caster}
	} else {
		return nil
	}
}
func PushBackNoDmg(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	if receiver.HasComponent(component.GridPos) {
		receiverPos := component.GridPos.Get(receiver)
		if receiverPos.Col < 7 {
			if validMove(ecs, receiverPos.Row, receiverPos.Col+1) {
				receiverPos.Col += 1
				scrPos := component.ScreenPos.Get(receiver)
				scrPos.X = 0
				scrPos.Y = 0
			}

		}
	}
}
func (l *PushbackDecorator) GetModifierEntry() *loadout.CasterModifierData {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		return cc.GetModifierEntry()
	}
	return nil
}
func (l *PushbackDecorator) GetElement() component.Elemental {
	if vv, ok := l.caster.(loadout.ElementalCaster); ok {
		return vv.GetElement()
	}
	return component.NEUTRAL
}
func (l *PushbackDecorator) SetModifier(e *loadout.CasterModifierData) {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		cc.SetModifier(e)
	}
}
func (h *PushbackDecorator) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	// h.ensource = ensource
	h.caster.Cast(ensource, ecs)
}

func (h *PushbackDecorator) GetCost() int {
	return h.caster.GetCost() * 2
}
func (h *PushbackDecorator) GetIcon() *ebiten.Image {
	return h.caster.GetIcon()
}
func (h *PushbackDecorator) GetCooldown() time.Time {
	return h.caster.GetCooldown()
}
func (h *PushbackDecorator) GetCooldownDuration() time.Duration {
	return h.caster.GetCooldownDuration()
}
func (h *PushbackDecorator) GetDamage() int {
	return h.caster.GetDamage()
}

func (h *PushbackDecorator) ResetCooldown() {
	h.caster.ResetCooldown()
}

func (h *PushbackDecorator) GetDescription() string {
	return h.caster.GetDescription() + "\nOn hit push target back"
}
func (h *PushbackDecorator) GetName() string {
	return h.caster.GetName() + " +Push"
}
