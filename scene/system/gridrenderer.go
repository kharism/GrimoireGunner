package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	"github.com/kharism/grimoiregunner/scene/assets"
	asset "github.com/kharism/grimoiregunner/scene/assets"
	myComponent "github.com/kharism/grimoiregunner/scene/component"
)

type gridRenderer struct {
	query *donburi.Query
	// orderedQuery *donburi.OrderedQuery[component.PositionData]
	queryDamage *donburi.Query
}

// This renderer will render floor tile
var GridRenderer = &gridRenderer{
	query: donburi.NewQuery(
		filter.Contains(
			myComponent.TileTag,
			myComponent.GridPos,
			myComponent.ScreenPos,
		),
	),
	queryDamage: donburi.NewQuery(
		filter.Contains(
			myComponent.Damage,
			myComponent.GridPos,
		),
	),
}

var TileStartX = float64(165.0)
var TileStartY = float64(360.0)

// return col,row
func GridCoord2Screen(Row, Col int) (float64, float64) {
	return TileStartX + float64(Col)*float64(assets.TileWidth), TileStartY + float64(Row)*float64(assets.TileHeight)
}

// param screen X,Y coords
// return col,row
func Coord2Grid(X, Y float64) (int, int) {
	col := int(X-TileStartX) / assets.TileWidth
	row := int(Y-TileStartY) / assets.TileHeight
	return col, row
}
func (r *gridRenderer) DrawGrid(ecs *ecs.ECS, screen *ebiten.Image) {
	// screen.Fill(color.RGBA{R: 41, G: 44, B: 45, A: 255})
	r.query.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := myComponent.GridPos.Get(e)
		// screenPos := myComponent.GridPos.Get(e)
		var sprite *ebiten.Image
		if gridPos.Col < 4 {
			sprite = asset.BlueTile
		} else {
			sprite = asset.RedTile
		}

		translate := ebiten.GeoM{}
		translate.Translate(-float64(assets.TileWidth)/2, -float64(assets.TileHeight))
		translate.Translate(TileStartX+float64(gridPos.Col)*float64(assets.TileWidth), TileStartY+float64(gridPos.Row)*float64(assets.TileHeight))
		drawOption := &ebiten.DrawImageOptions{
			GeoM: translate,
		}
		screen.DrawImage(sprite, drawOption)
	})
	r.queryDamage.Each(ecs.World, func(e *donburi.Entry) {
		gridPos := myComponent.GridPos.Get(e)
		translate := ebiten.GeoM{}
		translate.Translate(-float64(assets.TileWidth)/2, -float64(assets.TileHeight))
		translate.Translate(TileStartX+float64(gridPos.Col)*float64(assets.TileWidth), TileStartY+float64(gridPos.Row)*float64(assets.TileHeight))
		drawOption := &ebiten.DrawImageOptions{
			GeoM: translate,
		}
		screen.DrawImage(asset.DamageGrid, drawOption)
	})
}
