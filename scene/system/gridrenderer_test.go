package system

import (
	"testing"

	"github.com/kharism/grimoiregunner/scene/assets"
)

func TestCoordTrans(t *testing.T) {
	X, Y := 265.0, 320.0
	assets.TileHeight = 50.0
	assets.TileWidth = 100
	col, row := assets.Coord2Grid(X, Y)
	if col != 1 {
		t.Fail()
	}
	if row != 0 {
		t.Fail()
	}
}
