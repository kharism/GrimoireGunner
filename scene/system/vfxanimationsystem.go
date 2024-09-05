package system

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateFx(ecs *ecs.ECS) {
	component.Fx.Each(ecs.World, func(e *donburi.Entry) {
		if e.HasComponent(component.Fx) {
			fx := component.Fx.GetValue(e)
			fx.Animation.Update()
		}

	})
}

func RenderFx(ecs *ecs.ECS, screen *ebiten.Image) {
	sorted := []*component.FxData{}
	component.Fx.Each(ecs.World, func(e *donburi.Entry) {
		fx := component.Fx.Get(e)
		sorted = append(sorted, fx)
		// fx.Animation.Draw(screen)
	})
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].GetPriority() < sorted[j].GetPriority() })
	for _, fx := range sorted {
		fx.Animation.Draw(screen)
	}
}
