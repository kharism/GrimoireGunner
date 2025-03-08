package scene

import (
	"math/rand"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/enemies"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
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

type CombatSceneDecorator func(*ecs.ECS, *CombatScene)

var Decorators1 = []CombatSceneDecorator{}
var Decorators2 = []CombatSceneDecorator{}

var OptBoss1 = []CombatSceneDecorator{}
var OptBoss1Name = []string{}

func init() {
	Decorators1 = []CombatSceneDecorator{
		level1Decorator1,
		level1Decorator2,
		level1Decorator3,
		level1Decorator4,
		level1Decorator5,
		level1Decorator6,
		level1Decorator7,
		level1Decorator8,
		level1Decorator10,
		level1Decorator11,
		level1Decorator12,
		level1Decorator13,
	}
	Decorators2 = []CombatSceneDecorator{
		level1Decorator10,
		level1Decorator11,
		level1Decorator12,
		level2Decorator1,
		level2Decorator4,
		level2Decorator3,
		level2Decorator5,
		level2Decorator6,
		level2Decorator7,
		level2Decorator8,
	}
	OptBoss1 = []CombatSceneDecorator{
		level1OptBoss1,
		level1OptBoss2,
		level1OptBoss3,
	}
	OptBoss1Name = []string{
		"Ellone",
		"Yanman",
		"Joji",
	}
}
func RandCombatDecorator1() CombatSceneDecorator {
	i := rand.Int() % len(Decorators1)
	return Decorators1[i]
}

func RandCombatDecorator2() CombatSceneDecorator {
	i := rand.Int() % len(Decorators2)
	return Decorators2[i]
}

func RandBossDecorator1() (CombatSceneDecorator, string) {
	i := rand.Int() % len(OptBoss1)
	return OptBoss1[i], OptBoss1Name[i]
}

// put in 1 rock and 1 cannoneer and 1 rock
func level1Decorator1(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	LoadBoulder(ecs.World, BoulderParam{
		Col: 5,
		Row: 0,
	})
	enemies.NewCannoneer(ecs, 6, 1)
}

// put 1 bloombomber
func level1Decorator2(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	LoadBoulder(ecs.World, BoulderParam{
		Col: 5,
		Row: 0,
	})
	enemies.NewBloombomber(ecs, 6, 0)
}

// put 1 gatlinghoul
func level1Decorator3(ecs *ecs.ECS, combatscene *CombatScene) {
	// enemies.NewGatlingGhoul(ecs, 4, 0)
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewGatlingGhoul(ecs, 4, 3)
}

// put 1 gatlinghoul and 1 reaper
func level1Decorator4(ecs *ecs.ECS, combatscene *CombatScene) {
	// enemies.NewGatlingGhoul(ecs, 4, 0)
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewGatlingGhoul(ecs, 4, 3)
	enemies.NewReaper(ecs, 4, 2)
}

func level1Decorator5(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	LoadBoulder(ecs.World, BoulderParam{
		Col: 4,
		Row: 2,
	})
	enemies.NewHammerghoul(ecs, 5, 2)
}
func level1Decorator6(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewHealslime(ecs, 6, 2)
}

func level1Decorator7(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewDemon(ecs, 4, 2)
}
func level1Decorator8(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = []ItemInterface{}
	for i := 0; i < 3; i++ {
		var temp loadout.Caster
		for {
			temp = DecorateCaster(GenerateCaster())
			if temp != nil {
				break
			}
		}
		combatscene.rewards = append(combatscene.rewards, temp)
	}
	enemies.NewGatlingGhoul(ecs, 6, 1)
	enemies.NewDemon(ecs, 4, 2)
}
func level1OptBoss1(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewSwordwomen(ecs, 5, 1)
}
func level1OptBoss2(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewYanman(ecs, 6, 1)
}
func level1OptBoss3(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewFrostYeti(ecs, 6, 1)
}
func level1Decorator10(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewPoacher(ecs, 6, 1)
}
func level1Decorator11(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewPyroEyes(ecs, 4, 1)
}
func level1Decorator12(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewLightningImp(ecs, 4, 1)
}
func level1Decorator13(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewGatlingGhoulOmega(ecs, 6, 1)
}
func level2Decorator1(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = []ItemInterface{}
	for i := 0; i < 3; i++ {
		var temp loadout.Caster
		for {
			temp = DecorateCaster(GenerateCaster())
			if temp != nil {
				break
			}
		}
		combatscene.rewards = append(combatscene.rewards, temp)
	}

	enemies.NewGatlingGhoul(ecs, 6, 1)
	enemies.NewDemon(ecs, 4, 2)
}
func level2Decorator2(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewPoacher(ecs, 6, 1)
}
func level2Decorator3(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewYeti(ecs, 6, 1)
}
func level2Decorator4(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgCave
	combatscene.rewards = nil
	enemies.NewIceslime(ecs, 6, 1)
}
func level2Decorator5(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgCave
	combatscene.rewards = nil
	enemies.NewIceslime(ecs, 6, 1)
	enemies.NewHealslime(ecs, 4, 2)
}
func level2Decorator6(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgCave
	combatscene.rewards = nil
	enemies.NewBuzzer(ecs, 6, 1)
	enemies.NewHealslime(ecs, 5, 1)
}
func level2Decorator7(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewLightningImp(ecs, 4, 1)
	enemies.NewGatlingGhoul(ecs, 6, 2)
	LoadBoulder(ecs.World, BoulderParam{
		Col: 5,
		Row: 2,
	})
}
func level2Decorator8(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgCave
	combatscene.rewards = nil
	enemies.NewLightningMage(ecs, 4, 1)
}
func finalBoss(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgCave
	combatscene.rewards = nil
	enemies.NewWhiteSnake(ecs, 6, 2)

}
