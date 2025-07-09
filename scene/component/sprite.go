package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
)

type SpriteData struct {
	Image *ebiten.Image
	Scale *core.ScaleParam
}

var Sprite = donburi.NewComponentType[SpriteData]()
