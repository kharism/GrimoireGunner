package system

import "testing"

func TestCoordTrans(t *testing.T) {
	X, Y := 265.0, 320.0
	tileHeight = 50.0
	tileWidth = 100
	col, row := Coord2Grid(X, Y)
	if col != 1 {
		t.Fail()
	}
	if row != 0 {
		t.Fail()
	}
}
