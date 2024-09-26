package attack

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi/ecs"
)

type DoubleCastDecor struct {
	// ensource loadout.ENSetGetter
	caster loadout.Caster
}

func DecorateWithDoubleCast(caster loadout.Caster) loadout.Caster {
	if cc, ok := caster.(loadout.ModifierGetSetter); ok {
		// newModifier := ecs.World.Create(component.CasterModifier, component.PostAtkModifier)
		ll := cc.GetModifierEntry()
		if ll == nil {
			ll = &loadout.CasterModifierData{}

		}
		ll.PostAtk = DoubleCast(caster)
		cc.SetModifier(ll)
		// entry := ecs.World.Entry(newModifier)
		// component.PostAtkModifier.SetValue(entry, DoubleCast(caster, entry))
		// cc.SetModifier(entry)
		return &DoubleCastDecor{caster: caster}
	} else {
		return nil
	}
}

type Add2ndShotEvent struct {
	time     time.Time
	ShotFunc func()
}

func (a *Add2ndShotEvent) Execute(ecs *ecs.ECS) {
	a.ShotFunc()
}
func (a *Add2ndShotEvent) GetTime() time.Time {
	return a.time
}
func DoubleCast(caster loadout.Caster) func(*ecs.ECS, loadout.ENSetGetter) {
	return func(ecs *ecs.ECS, ensource loadout.ENSetGetter) {
		// nextShot := Add2ndShotEvent{time: time.Now().Add(50 * time.Millisecond), ShotFunc: func() {
		// 	// ff := component.PostAtkModifier.GetValue(entry)
		// 	var ff component.PostAtkBehaviour
		// 	if jj, ok := caster.(ModifierGetSetter); ok {
		// 		ff = jj.GetModifierEntry().PostAtk
		// 		jj.GetModifierEntry().PostAtk = nil
		// 		caster.Cast(ensource, ecs)
		// 		jj.GetModifierEntry().PostAtk = ff
		// 	}

		// 	// component.PostAtkModifier.SetValue(entry, ff)
		// }}
		var ff loadout.PostAtkBehaviour
		if jj, ok := caster.(loadout.ModifierGetSetter); ok {
			ff = jj.GetModifierEntry().PostAtk
			jj.GetModifierEntry().PostAtk = nil
			caster.Cast(ensource, ecs)
			jj.GetModifierEntry().PostAtk = ff
		}
		// component.EventQueue.AddEvent(&nextShot)

	}
}
func (l *DoubleCastDecor) GetModifierEntry() *loadout.CasterModifierData {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		return cc.GetModifierEntry()
	}
	return nil
}
func (l *DoubleCastDecor) SetModifier(e *loadout.CasterModifierData) {
	if cc, ok := l.caster.(loadout.ModifierGetSetter); ok {
		cc.SetModifier(e)
	}
}
func (h *DoubleCastDecor) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	// h.ensource = ensource
	h.caster.Cast(ensource, ecs)
}

func (h *DoubleCastDecor) GetCost() int {
	return h.caster.GetCost() * 2
}
func (h *DoubleCastDecor) GetIcon() *ebiten.Image {
	return h.caster.GetIcon()
}
func (h *DoubleCastDecor) GetCooldown() time.Time {
	return h.caster.GetCooldown()
}
func (h *DoubleCastDecor) GetCooldownDuration() time.Duration {
	return h.caster.GetCooldownDuration()
}
func (h *DoubleCastDecor) GetDamage() int {
	return h.caster.GetDamage()
}

func (h *DoubleCastDecor) ResetCooldown() {
	h.caster.ResetCooldown()
}

func (h *DoubleCastDecor) GetDescription() string {
	return h.caster.GetDescription() + "\nDouble cost/cast"
}
func (h *DoubleCastDecor) GetName() string {
	return h.caster.GetName() + " +Double"
}
