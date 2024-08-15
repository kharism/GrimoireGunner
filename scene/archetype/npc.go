package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
)

var NpcTag = donburi.NewTag("NPC")

// create an NPC and assign sprite to it
func NewNPC(world donburi.World, sprite *ebiten.Image) *donburi.Entity {
	entity := world.Create(
		mycomponent.Health,
		mycomponent.GridPos,
		mycomponent.ScreenPos,
		mycomponent.Speed,
		mycomponent.EnemyRoutine,
		mycomponent.TargetLocation,
		mycomponent.Sprite,
	)
	mycomponent.Sprite.Set(world.Entry(entity), &mycomponent.SpriteData{Image: sprite})
	return &entity
}
