package system

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type characterRenderer struct {
	query *donburi.OrderedQuery[component.GridPosComponentData]
}

// Render anything that has sprite and grid position
var CharacterRenderer = &characterRenderer{
	query: donburi.NewOrderedQuery[component.GridPosComponentData](
		filter.Contains(
			mycomponent.Sprite,
			mycomponent.GridPos,
		),
	),
}

func (r *characterRenderer) DrawCharacter(ecs *ecs.ECS, screen *ebiten.Image) {
	entries := []*donburi.Entry{}
	r.query.Each(ecs.World, func(e *donburi.Entry) {
		entries = append(entries, e)
	})
	sort.Slice(entries, func(i, j int) bool {
		gridPosI := mycomponent.GridPos.Get(entries[i])
		gridPosJ := mycomponent.GridPos.Get(entries[j])
		return gridPosI.Order() < gridPosJ.Order()
	})
	for _, e := range entries {
		gridPos := mycomponent.GridPos.Get(e)

		// fmt.Println(e.Entity(), gridPos.Col, gridPos.Order(), component.Health.Get(e).Name)
		screenPos := mycomponent.ScreenPos.Get(e)
		if screenPos.X == 0 && screenPos.Y == 0 {
			screenPos.X = assets.TileStartX + float64(gridPos.Col)*float64(assets.TileWidth)
			screenPos.Y = assets.TileStartY + float64(gridPos.Row)*float64(assets.TileHeight)
		}
		sprite := mycomponent.Sprite.Get(e).Image
		bound := sprite.Bounds()
		translate := ebiten.GeoM{}
		translate.Translate(-float64(bound.Dx())/2, -float64(bound.Dy()))
		translate.Translate(screenPos.X, screenPos.Y)
		drawOption := &ebiten.DrawImageOptions{
			GeoM: translate,
		}
		screen.DrawImage(sprite, drawOption)
	}
}
