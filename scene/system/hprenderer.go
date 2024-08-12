package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type hpRenderer struct {
	query *donburi.Query
	// orderedQuery *donburi.OrderedQuery[component.PositionData]
}

var NPCRenderer = &hpRenderer{
	query: donburi.NewQuery(
		filter.Contains(mycomponent.Health),
	),
}

func (r *hpRenderer) DrawHP(ecs *ecs.ECS, screen *ebiten.Image) {
	r.query.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := mycomponent.GridPos.Get(e)
		screenPos := mycomponent.ScreenPos.Get(e)
		if screenPos.X == 0 && screenPos.Y == 0 {
			screenPos.X = TileStartX + float64(gridPos.Col)*float64(tileWidth)
			screenPos.Y = TileStartY + float64(gridPos.Row)*float64(tileHeight)
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
	})
}
