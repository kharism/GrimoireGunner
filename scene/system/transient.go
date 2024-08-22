package system

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type TransientSystem struct {
	query *donburi.Query
}

func NewTransientSystem() *TransientSystem {
	return &TransientSystem{
		query: donburi.NewQuery(
			filter.Contains(
				component.Transient,
			),
		),
	}
}
func (t *TransientSystem) Update(ecs *ecs.ECS) {
	t.query.Each(ecs.World, func(e *donburi.Entry) {
		transientData := component.Transient.Get(e)
		now := time.Now()
		if now.After(transientData.Start.Add(transientData.Duration)) {
			ecs.World.Remove(e.Entity())
		}
	})
}
