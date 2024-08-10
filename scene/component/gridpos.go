package component

import "github.com/yohamta/donburi"

type GridPosComponentData struct {
	Row, Col int
}

var GridPos = donburi.NewComponentType[GridPosComponentData]()
