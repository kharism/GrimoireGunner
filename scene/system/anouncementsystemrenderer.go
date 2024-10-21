package system

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateAnouncement(ecs *ecs.ECS) {
	component.Anouncement.Each(ecs.World, func(e *donburi.Entry) {
		if e.HasComponent(component.Anouncement) {
			fx := component.Anouncement.GetValue(e)
			fx.Animation.Update()
		}

	})
}

func RenderAnouncement(ecs *ecs.ECS, screen *ebiten.Image) {
	sorted := []*component.FxData{}
	component.Anouncement.Each(ecs.World, func(e *donburi.Entry) {
		fx := component.Anouncement.Get(e)
		sorted = append(sorted, fx)
		// fx.Animation.Draw(screen)
	})
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].GetPriority() < sorted[j].GetPriority() })
	for _, fx := range sorted {
		fx.Animation.Draw(screen)
	}
}
