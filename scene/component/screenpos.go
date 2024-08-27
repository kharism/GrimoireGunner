package component

import "github.com/yohamta/donburi"

type ScreenPosComponentData struct {
	X, Y float64
}

// ScreenPos define position of object on the screen
var ScreenPos = donburi.NewComponentType[ScreenPosComponentData]()
