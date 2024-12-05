package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
)

var PlayerTag = donburi.NewTag("Player")

func NewPlayer(world donburi.World, sprite *ebiten.Image) *donburi.Entity {
	entity := world.Create(
		mycomponent.Health,
		mycomponent.GridPos,
		mycomponent.ScreenPos,
		// mycomponent.Speed,
		// mycomponent.TargetLocation,
		mycomponent.PlayerDataComponent,
		mycomponent.Elements,
		// mycomponent.Shader,
		mycomponent.Sprite, PlayerTag)
	entId := world.Entry(entity)
	mycomponent.Sprite.Set(entId, &mycomponent.SpriteData{Image: sprite})
	mycomponent.Health.Set(entId, &mycomponent.HealthData{HP: 1000, Name: "Player", OnTakeDamage: mycomponent.AddIFrame})
	gridPos := &mycomponent.GridPosComponentData{Row: 0, Col: 0}
	mycomponent.GridPos.Set(entId, gridPos)
	mycomponent.PlayerDataComponent.Set(entId, mycomponent.NewPlayerData())
	mycomponent.ScreenPos.Set(entId, &mycomponent.ScreenPosComponentData{X: 0, Y: 0})
	// mycomponent.Speed.Set(entId, &mycomponent.SpeedData{Vx: 0, Vy: 0})
	// mycomponent.TargetLocation.Set(entId, &mycomponent.MoveTargetData{Tx: 0, Ty: 0})
	return &entity
}
