package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
)

type FxData struct {
	Animation Animation
}
type Animation interface {
	Draw(screen *ebiten.Image)
	Update()
}

func (f *FxData) GetPriority() int {
	if jj, ok := f.Animation.(*core.MovableImage); ok {
		x, y := jj.GetPos()
		return int(10*y + x)
	} else if jj, ok := f.Animation.(*core.AnimatedImage); ok {
		x, y := jj.GetPos()
		return int(10*y + x)
	}
	return 99
}

var Fx = donburi.NewComponentType[FxData]()
