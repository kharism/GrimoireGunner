package scene

import (
	"math/rand"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/enemies"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type BoulderParam struct {
	Col, Row int
}

func LoadBoulder(world donburi.World, param BoulderParam) *donburi.Entity {
	entity := world.Create(
		mycomponent.Health,
		mycomponent.GridPos,
		mycomponent.ScreenPos,
		mycomponent.Sprite,
		archetype.ConstructTag,
	)
	entry := world.Entry(entity)
	mycomponent.Sprite.Set(entry, &mycomponent.SpriteData{Image: assets.Boulder})
	mycomponent.Health.Set(entry, &mycomponent.HealthData{HP: 200, Name: "Boulder"})
	mycomponent.GridPos.Set(entry, &mycomponent.GridPosComponentData{Col: param.Col, Row: param.Row})
	return &entity
}

type CombatSceneDecorator func(*ecs.ECS)

var Decorators = []CombatSceneDecorator{}

func init() {
	Decorators = []CombatSceneDecorator{
		level1Decorator1,
		level1Decorator2,
		level1Decorator3,
		level1Decorator4,
	}
}
func RandDecorator() CombatSceneDecorator {
	i := rand.Int() % len(Decorators)
	return Decorators[i]
}

// put in 1 rock and 1 cannoneer and 1 rock
func level1Decorator1(ecs *ecs.ECS) {
	LoadBoulder(ecs.World, BoulderParam{
		Col: 5,
		Row: 0,
	})
	enemies.NewCannoneer(ecs, 6, 1)
}

// put 1 bloombomber
func level1Decorator2(ecs *ecs.ECS) {
	LoadBoulder(ecs.World, BoulderParam{
		Col: 5,
		Row: 0,
	})
	enemies.NewBloombomber(ecs, 6, 0)
}

// put 1 gatlinghoul
func level1Decorator3(ecs *ecs.ECS) {
	// enemies.NewGatlingGhoul(ecs, 4, 0)
	enemies.NewGatlingGhoul(ecs, 4, 3)
}

// put 1 gatlinghoul and 1 reaper
func level1Decorator4(ecs *ecs.ECS) {
	// enemies.NewGatlingGhoul(ecs, 4, 0)
	enemies.NewGatlingGhoul(ecs, 4, 3)
	enemies.NewReaper(ecs, 4, 2)
}

func Level1Decorator5(ecs *ecs.ECS) {
	enemies.NewHammerghoul(ecs, 5, 2)
}
