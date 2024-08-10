package component

import "github.com/yohamta/donburi"

type ScreenPosComponentData struct {
	X, Y float64
}

var ScreenPos = donburi.NewComponentType[ScreenPosComponentData]()
