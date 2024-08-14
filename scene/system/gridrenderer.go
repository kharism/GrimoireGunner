package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	asset "github.com/kharism/grimoiregunner/scene/assets"
	myComponent "github.com/kharism/grimoiregunner/scene/component"
)

type gridRenderer struct {
	query *donburi.Query
	// orderedQuery *donburi.OrderedQuery[component.PositionData]
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
}
var tileWidth int
var tileHeight int

var TileStartX = float64(165.0)
var TileStartY = float64(360.0)

// return col,row
func GridCoord2Screen(Row, Col int) (float64, float64) {
	return TileStartX + float64(Col)*float64(tileWidth), TileStartY + float64(Row)*float64(tileHeight)
}

// param screen X,Y coords
// return col,row
func Coord2Grid(X, Y float64) (int, int) {
	col := int(X-TileStartX) / tileWidth
	row := int(Y-TileStartY) / tileHeight
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
		if tileWidth == 0 {
			rect := sprite.Bounds()
			tileWidth = rect.Dx()
			tileHeight = rect.Dy()
		}

		translate := ebiten.GeoM{}
		translate.Translate(-float64(tileWidth)/2, -float64(tileHeight))
		translate.Translate(TileStartX+float64(gridPos.Col)*float64(tileWidth), TileStartY+float64(gridPos.Row)*float64(tileHeight))
		drawOption := &ebiten.DrawImageOptions{
			GeoM: translate,
		}
		screen.DrawImage(sprite, drawOption)
	})
}
