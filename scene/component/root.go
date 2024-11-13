package component

import (
	"time"

	"github.com/yohamta/donburi"
)

type RootData struct {
	ReleaseTime time.Time
}

var Root = donburi.NewComponentType[RootData]()
