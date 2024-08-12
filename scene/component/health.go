package component

import "github.com/yohamta/donburi"

type HealthData struct {
	HP int
}

var Health = donburi.NewComponentType[HealthData]()
