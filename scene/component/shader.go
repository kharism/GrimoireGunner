package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

var Shader = donburi.NewComponentType[ebiten.Shader]()
