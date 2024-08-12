package component

import "github.com/yohamta/donburi"

type SpeedData struct {
	Vx, Vy float64
}

var Speed = donburi.NewComponentType[SpeedData]()
