package component

import "github.com/yohamta/donburi"

type MoveTargetData struct {
	Tx, Ty float64
}

var TargetLocation = donburi.NewComponentType[MoveTargetData]()
