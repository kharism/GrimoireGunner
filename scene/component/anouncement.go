package component

import (
	"github.com/yohamta/donburi"
)

// this is basically same with FX, just so we can use different layer when rendering
var Anouncement = donburi.NewComponentType[FxData]()
