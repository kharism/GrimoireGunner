package component

import "github.com/yohamta/donburi"

type HealthData struct {
	HP   int
	Name string
}

var Health = donburi.NewComponentType[HealthData]()
