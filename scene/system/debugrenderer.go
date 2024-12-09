package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type debugRenderer struct {
	query *donburi.Query
	// orderedQuery *donburi.OrderedQuery[component.PositionData]
}

var DebugRenderer = &debugRenderer{
	query: donburi.NewQuery(
		filter.Contains(mycomponent.Health),
		// filter.Contains(component.Fx),
	),
}

func (r *debugRenderer) DrawDebug(ecs *ecs.ECS, screen *ebiten.Image) {
	fps := ebiten.ActualFPS()
	// op := &text.DrawOptions{
	// 	LayoutOptions: text.LayoutOptions{
	// 		PrimaryAlign: text.AlignStart,
	// 	},
	// }

	// text.Draw(screen, fmt.Sprintf("fps=%f", fps), assets.FontFace, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%f", fps))
	r.query.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := mycomponent.GridPos.Get(e)
		screenPos := mycomponent.ScreenPos.Get(e)
		if screenPos.X == 0 && screenPos.Y == 0 {
			screenPos.X = assets.TileStartX + float64(gridPos.Col)*float64(assets.TileWidth)
			screenPos.Y = assets.TileStartY + float64(gridPos.Row)*float64(assets.TileHeight)
		}
		translate := ebiten.GeoM{}
		translate.Translate(screenPos.X, screenPos.Y+24)
		op := &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: text.AlignCenter,
			},
			DrawImageOptions: ebiten.DrawImageOptions{
				GeoM: translate,
			},
		}
		// hp := mycomponent.Health.Get(e).HP
		// if animImage, ok := component.Fx.Get(e).Animation.(*core.AnimatedImage); ok {
		// animImage.CurrMove.(*core.MoveAnimation)
		// }
		if e.HasComponent(component.EnemyRoutine) {
			memory := component.EnemyRoutine.Get(e).Memory
			if cc, ok := memory["CurrMove"].(string); ok {
				text.Draw(screen, cc, assets.FontFace, op)
			}

		}

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
