package scene

import (
	"math/rand"

	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
)

func GenerateReward() ItemInterface {
	items := []ItemInterface{
		attack.NewBuckshotCaster(),
		attack.NewFirewallCaster(),
		attack.NewLongSwordCaster(),
		attack.NewShockwaveCaster(),
		attack.NewCannonCaster(),
		attack.NewWideSwordCaster(),
		attack.NewShotgunCaster(),
		attack.NewHealCaster(),
		attack.NewGatlingCastor(),
		attack.NewLightningBolCaster(),
		&Medkit{},
		&HPUp{},
		&ENUp{},
	}
	rnd := rand.Int() % len(items)
	return items[rnd]
}
func GenerateCaster() loadout.Caster {
	casters := []loadout.Caster{
		attack.NewBuckshotCaster(),
		attack.NewFirewallCaster(),
		attack.NewLongSwordCaster(),
		attack.NewShockwaveCaster(),
		attack.NewCannonCaster(),
		attack.NewWideSwordCaster(),
		attack.NewShotgunCaster(),
		attack.NewHealCaster(),
		attack.NewGatlingCastor(),
		attack.NewLightningBolCaster(),
	}

	rnd := rand.Int() % len(casters)
	return casters[rnd]

	// switch rng {
	// case 0:
	// 	return attack.NewBuckshotCaster()
	// case 1:
	// 	return attack.NewFirewallCaster()
	// case 2:
	// 	return attack.NewLongSwordCaster()
	// case 3:
	// 	return attack.NewShockwaveCaster()
	// case 4:
	// 	return attack.NewLightningBolCaster()

	// }
	// return attack.NewGatlingCastor()
}
