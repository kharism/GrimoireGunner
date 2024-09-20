package component

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// this component holds simple stat change for caster
type CasterModifierData struct {
	DamageModifier  int
	CooldownModifer time.Duration
	CostModifier    int
	SpecialModifier int
}
type PostAtkBehaviour func(*ecs.ECS)

var CasterModifier = donburi.NewComponentType[CasterModifierData]()

var PostAtkModifier = donburi.NewComponentType[PostAtkBehaviour]()
