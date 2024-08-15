package component

import (
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
)

var Fx = donburi.NewComponentType[*core.AnimatedImage]()
