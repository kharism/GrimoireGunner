package component

import (
	"time"

	"github.com/yohamta/donburi"
)

// store who gets to burn
type BurnerData struct {
	Damage int
}

// store data on who gets burned
type BurnedData struct {
	NextBurn   time.Time // the next time the burnable takes damage. Takes damage every 800 milliseconds by default
	BurnCount  int       // how many times the object has taken burn damage. Limit this to 10
	BurnDamage int       // amunt of burn damage
}

// the things that do the burning. Usually tile which is on fire.
// Attach this to grid tile for now, Probably
var Burner = donburi.NewComponentType[BurnerData]()
var Burned = donburi.NewComponentType[BurnedData]()
