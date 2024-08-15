package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateFx(ecs *ecs.ECS) {
	component.Fx.Each(ecs.World, func(e *donburi.Entry) {
		fx := component.Fx.GetValue(e)
		fx.Update()
	})
}

func RenderFx(ecs *ecs.ECS, screen *ebiten.Image) {
	component.Fx.Each(ecs.World, func(e *donburi.Entry) {
		fx := component.Fx.GetValue(e)
		fx.Draw(screen)
	})
}
