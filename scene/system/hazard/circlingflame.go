package hazard

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// create rotating flame
func NewRotatingFlame(ecs *ecs.ECS, col, row int) {
	flame := archetype.NewHazard(ecs.World, assets.FlametowerRaw)
	flameEnt := ecs.World.Entry(flame)
	flameEnt.AddComponent(component.Burner)
	gridPos := component.GridPos.Get(flameEnt)
	gridPos.Col = col
	gridPos.Row = row
	component.Burner.Set(flameEnt, &component.BurnerData{
		Damage:  20,
		Element: component.FIRE,
	})
	data := map[string]any{}

	data["WARM_UP"] = time.Now().Add(500 * time.Millisecond)

	component.EnemyRoutine.Set(flameEnt, &component.EnemyRoutineData{
		Routine: circlingAround,
		Memory:  data,
	})
}

var scanOrder1 = [20][2]int{
	{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, //move right
	{1, 7}, {2, 7}, {3, 7}, // down
	{3, 6}, {3, 5}, {3, 4}, {3, 3}, {3, 2}, {3, 1}, {3, 0}, // move left
	{2, 0}, {1, 0}, // up

}

func getOrder1Index(row, col int) int {
	for idx, val := range scanOrder1 {
		if val[0] == row && val[1] == col {
			return idx
		}
	}
	return -1
}
func circlingAround(ecs *ecs.ECS, self *donburi.Entry) {
	memory := component.EnemyRoutine.Get(self).Memory
	gridPos := component.GridPos.Get(self)
	curIdx := getOrder1Index(gridPos.Row, gridPos.Col)
	if waitTime, ok := memory["WARM_UP"].(time.Time); ok && waitTime.Before(time.Now()) {
		nextIdx := (curIdx + 1) % 20
		newPos := scanOrder1[nextIdx]
		gridPos.Row = newPos[0]
		gridPos.Col = newPos[1]
		scrPos := component.ScreenPos.Get(self)
		scrPos.X = 0
		scrPos.Y = 0
		memory["WARM_UP"] = time.Now().Add(500 * time.Millisecond)
	}
}
