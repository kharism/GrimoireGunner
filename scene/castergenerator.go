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
		attack.NewBombConstructCaster(),
		attack.NewAtkBonusCaster(),
		attack.NewPushgunCaster(),
		attack.NewChargeshotCaster(),
		attack.NewBambooLanceCaster(),
		attack.NewSporeBombCaster(),
		attack.NewWallCaster(),
		attack.NewIcespikeCaster(),
		attack.NewFlamethrowerCaster(),
		attack.NewFist(),
		&Medkit{},
		&HPUp{},
		&ENUp{},
	}
	rnd := rand.Int() % len(items)
	return items[rnd]
}
func GenerateTrieReward() (ItemInterface, ItemInterface, ItemInterface) {
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
		attack.NewBombConstructCaster(),
		attack.NewAtkBonusCaster(),
		attack.NewPushgunCaster(),
		attack.NewPullgunCaster(),
		attack.NewChargeshotCaster(),
		attack.NewBambooLanceCaster(),
		attack.NewSporeBombCaster(),
		attack.NewWallCaster(),
		attack.NewIcespikeCaster(),
		attack.NewFlamethrowerCaster(),
		attack.NewFist(),
		&Medkit{},
		&HPUp{},
		&ENUp{},
	}
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
	return items[0], items[1], items[3]
}
func DecorateCaster(caster loadout.Caster) loadout.Caster {
	rnd := rand.Int() % len(CasterDecorList)
	return CasterDecorList[rnd](caster)
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
		attack.NewChargeshotCaster(),
		attack.NewGatlingCastor(),
		attack.NewLightningBolCaster(),
		attack.NewBombConstructCaster(),
		attack.NewAtkBonusCaster(),
		attack.NewPushgunCaster(),
		attack.NewBambooLanceCaster(),
		attack.NewSporeBombCaster(),
		attack.NewWallCaster(),
		attack.NewFist(),
		attack.NewIcespikeCaster(),
		attack.NewFlamethrowerCaster(),
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
