package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type FxData struct {
	Animation Animation
}
type Animation interface {
	Draw(screen *ebiten.Image)
	Update()
}

var Fx = donburi.NewComponentType[FxData]()
