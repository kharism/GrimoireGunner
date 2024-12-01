package component

import "github.com/yohamta/donburi"

type Elemental int

const (
	NEUTRAL Elemental = iota
	FIRE
	WATER
	ELEC
	WOOD
)

func IsDoubleDamage(attacker, receiver Elemental) bool {
	if attacker == FIRE && receiver == WOOD {
		return true
	}
	if attacker == WATER && receiver == FIRE {
		return true
	}
	if attacker == ELEC && receiver == WATER {
		return true
	}
	if attacker == WOOD && receiver == ELEC {
		return true
	}
	return false
}

var Elements = donburi.NewComponentType[Elemental]()
