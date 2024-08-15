package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type hpRenderer struct {
	query *donburi.Query
	// orderedQuery *donburi.OrderedQuery[component.PositionData]
}

var HPRenderer = &hpRenderer{
	query: donburi.NewQuery(
		filter.Contains(mycomponent.Health),
	),
}

func (r *hpRenderer) DrawHP(ecs *ecs.ECS, screen *ebiten.Image) {
	r.query.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := mycomponent.GridPos.Get(e)
		screenPos := mycomponent.ScreenPos.Get(e)
		if screenPos.X == 0 && screenPos.Y == 0 {
			screenPos.X = TileStartX + float64(gridPos.Col)*float64(assets.TileWidth)
			screenPos.Y = TileStartY + float64(gridPos.Row)*float64(assets.TileHeight)
		}
		translate := ebiten.GeoM{}
		translate.Translate(screenPos.X, screenPos.Y)
		op := &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: text.AlignCenter,
			},
			DrawImageOptions: ebiten.DrawImageOptions{
				GeoM: translate,
			},
		}
		hp := mycomponent.Health.Get(e).HP
		text.Draw(screen, fmt.Sprintf("%d", hp), assets.FontFace, op)
		// sprite := mycomponent.Sprite.Get(e).Image
		// bound := sprite.Bounds()
		// translate := ebiten.GeoM{}
		// translate.Translate(-float64(bound.Dx())/2, -float64(bound.Dy()))
		// translate.Translate(screenPos.X, screenPos.Y)
		// drawOption := &ebiten.DrawImageOptions{
		// 	GeoM: translate,
		// }
		// screen.DrawImage(sprite, drawOption)
	})
}
