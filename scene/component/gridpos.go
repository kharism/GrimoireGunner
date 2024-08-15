package component

import "github.com/yohamta/donburi"

type GridPosComponentData struct {
	Row, Col int
}

func (f GridPosComponentData) Order() int {
	return 10*f.Row + f.Col
}

var GridPos = donburi.NewComponentType[GridPosComponentData]()
