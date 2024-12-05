package loadout

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// this component holds simple stat change for caster
type CasterModifierData struct {
	DamageModifier  int
	CooldownModifer time.Duration
	CostModifier    int
	SpecialModifier int
	Element         component.Elemental
	PostAtk         PostAtkBehaviour
	OnHit           component.OnAtkHit
}
type PostAtkBehaviour func(*ecs.ECS, ENSetGetter)

var CasterModifier = donburi.NewComponentType[CasterModifierData]()

var PostAtkModifier = donburi.NewComponentType[PostAtkBehaviour]()
