package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/yohamta/donburi/ecs"
)

func DrawBg(ecs *ecs.ECS, screen *ebiten.Image) {
	drawOpt := ebiten.DrawImageOptions{}
	screen.DrawImage(assets.Bg, &drawOpt)
}
