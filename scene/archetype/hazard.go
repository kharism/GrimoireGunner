package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
)

func NewHazard(world donburi.World, sprite *ebiten.Image) donburi.Entity {
	entity := world.Create(
		// mycomponent.Health,
		mycomponent.GridPos,
		mycomponent.ScreenPos,
		// mycomponent.Speed,
		// mycomponent.TargetLocation,
		mycomponent.EnemyRoutine,
		mycomponent.Sprite, ConstructTag)
	entId := world.Entry(entity)
	mycomponent.Sprite.Set(entId, &mycomponent.SpriteData{Image: sprite})
	// mycomponent.Health.Set(entId, &mycomponent.HealthData{HP: 200})
	mycomponent.GridPos.Set(entId, &mycomponent.GridPosComponentData{Row: 1, Col: 1})
	mycomponent.ScreenPos.Set(entId, &mycomponent.ScreenPosComponentData{X: 0, Y: 0})
	return entity
}
