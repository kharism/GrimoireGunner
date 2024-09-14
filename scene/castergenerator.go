package scene

import (
	"math/rand"

	"github.com/kharism/grimoiregunner/scene/system"
	"github.com/kharism/grimoiregunner/scene/system/attack"
)

func GenerateCaster() system.Caster {
	rng := rand.Int() % 5
	switch rng {
	case 0:
		return attack.NewBuckshotCaster()
	case 1:
		return attack.NewFirewallCaster()
	case 2:
		return attack.NewLongSwordCaster()
	case 3:
		return attack.NewShockwaveCaster()
	case 4:
		return attack.NewLightningBolCaster()

	}
	return attack.NewGatlingCastor()
}
