package component

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type CasterModifierData struct {
	DamageModifier  int
	CooldownModifer time.Duration
	CostModifier    int
	SpecialModifier int

	//execute this after attack is hit
	PostAtkBehaviour func(*ecs.ECS)
}

var CasterModifier = donburi.NewComponentType[CasterModifierData]()
