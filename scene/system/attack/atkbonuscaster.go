package attack

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi/ecs"
)

type AtkBonusCaster struct {
	nextCooldown time.Time
}

func NewAtkBonusCaster() *AtkBonusCaster {
	return &AtkBonusCaster{nextCooldown: time.Now()}
}
func (l *AtkBonusCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (a *AtkBonusCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= a.GetCost() {
		ensource.SetEn(en - a.GetCost())
		a.nextCooldown = time.Now().Add(a.GetCooldownDuration())
		var caster loadout.Caster
		if loadout.CurLoadOut[0] == a {
			caster = loadout.CurLoadOut[1]
		} else {
			caster = loadout.CurLoadOut[0]
		}
		if l, ok := caster.(ModifierGetSetter); ok {
			mod := l.GetModifierEntry()
			if mod == nil {
				// entity := ecs.World.Create(component.CasterModifier, component.PostAtkModifier)
				mod = &component.CasterModifierData{}
				l.SetModifier(mod)
			}

			mod.DamageModifier += 10
			oriVal := mod.PostAtk
			mod.PostAtk = RemoveAtk(l, oriVal)
			l.SetModifier(mod)
			// component.PostAtkModifier.SetValue( RemoveAtk(l, oriVal))

		}
	}

}
func RemoveAtk(origin ModifierGetSetter, originalFunc func(*ecs.ECS, loadout.ENSetGetter)) func(*ecs.ECS, loadout.ENSetGetter) {
	return func(ecs *ecs.ECS, ensource loadout.ENSetGetter) {
		if originalFunc != nil {
			originalFunc(ecs, ensource)
		}

		mod := origin.GetModifierEntry()
		// ecs.World.Remove(mod.Entity())
		// origin.SetModifier(nil)
		// component.CasterModifier.Get(mod).DamageModifier -= 10
		mod.DamageModifier -= 10
		// component.PostAtkModifier.SetValue(mod, originalFunc)
		mod.PostAtk = originalFunc
		// mod.RemoveComponent(component.PostAtkModifier)
	}
}
func (a *AtkBonusCaster) GetCost() int {
	return 300
}
func (a *AtkBonusCaster) GetIcon() *ebiten.Image {
	return assets.AtkUp
}
func (a *AtkBonusCaster) GetCooldown() time.Time {
	return a.nextCooldown
}
func (a *AtkBonusCaster) GetCooldownDuration() time.Duration {
	return 20 * time.Second
}
func (a *AtkBonusCaster) GetDamage() int {
	return 0
}

func (a *AtkBonusCaster) GetDescription() string {
	return "Give +10 damage to caster on other slot"
}
func (a *AtkBonusCaster) GetName() string {
	return "AtkBonus"
}
