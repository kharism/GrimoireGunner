package scene

import (
	"math/rand"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/grimoiregunner/scene/system/enemies"
	"github.com/kharism/grimoiregunner/scene/system/hazard"
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
func LoadBomb(world donburi.World, param BoulderParam) *donburi.Entity {
	bombEntity := archetype.NewConstruct(world, assets.Bomb2)
	bombEntry := world.Entry(*bombEntity)
	bombEntry.AddComponent(component.OnDestroy)
	bombEntry.AddComponent(component.OnHit)

	bombGridPos := component.GridPos.Get(bombEntry)
	bombGridPos.Col = param.Col
	bombGridPos.Row = param.Row
	component.OnDestroy.SetValue(bombEntry, attack.OnBombDestroyed)
	component.OnHit.SetValue(bombEntry, attack.SingleHitProjectile)
	return bombEntity
}

type CombatSceneDecorator func(*ecs.ECS, *CombatScene)

var Decorators1 = []CombatSceneDecorator{}
var Decorators2 = []CombatSceneDecorator{}

var OptBoss1 = []CombatSceneDecorator{}
var OptBoss1Name = []string{}

var OptBoss2 = []CombatSceneDecorator{}
var OptBoss2Name = []string{}

func init() {
	Decorators1 = []CombatSceneDecorator{
		// level1Decorator1,
		// level1Decorator2,
		level1Decorator3,
		level1Decorator4,
		level1Decorator8,
		// level1Decorator10,
		// level1Decorator11,
		level1Decorator12,
		// level1Decorator14,
		level1Decorator16,
		level1Decorator17,
		level1Decorator18,
		level2Decorator4,
		level2Decorator11,
		// level1WavesDecor2,
	}
	Decorators2 = []CombatSceneDecorator{
		level1Decorator10,
		level1Decorator11,
		level1Decorator12,
		level2Decorator1,
		// level2Decorator2,
		level2Decorator4,
		level2Decorator3,
		level2Decorator5,
		level2Decorator6,
		level2Decorator7,
		level2Decorator8,
		level2Decorator9,
		level2Decorator10,
		level2Decorator11,
		level2Decorator12,
	}
	OptBoss1 = []CombatSceneDecorator{
		level1OptBoss1,
		level1OptBoss2,
		level1OptBoss3,
		level1OptBoss4,
	}
	OptBoss1Name = []string{
		"Ellone",
		"Yanman",
		"Joji",
		"Morty",
	}
	OptBoss2 = []CombatSceneDecorator{
		level2OptBoss1,
		level2OptBoss2,
	}
	OptBoss2Name = []string{
		"Ellone",
		"Joji",
	}
}

var lv1Idx = -1
var lv2Idx = -1

func init() {
	rand.Shuffle(len(Decorators1), func(i, j int) {
		Decorators1[i], Decorators1[j] = Decorators1[j], Decorators1[i]
	})
	rand.Shuffle(len(Decorators2), func(i, j int) {
		Decorators2[i], Decorators2[j] = Decorators2[j], Decorators2[i]
	})
}
func ShuffleLevellayout() {
	rand.Shuffle(len(Decorators1), func(i, j int) {
		Decorators1[i], Decorators1[j] = Decorators1[j], Decorators1[i]
	})
	rand.Shuffle(len(Decorators2), func(i, j int) {
		Decorators2[i], Decorators2[j] = Decorators2[j], Decorators2[i]
	})
	lv1Idx = -1
	lv2Idx = -1
}
func RandCombatDecorator1() CombatSceneDecorator {
	lv1Idx += 1
	return Decorators1[lv1Idx]
}

func RandCombatDecorator2() CombatSceneDecorator {
	lv2Idx += 1
	return Decorators2[lv2Idx]
}

func RandBossDecorator1() (CombatSceneDecorator, string) {
	i := rand.Int() % len(OptBoss1)
	return OptBoss1[i], OptBoss1Name[i]
}

func RandBossDecorator2() (CombatSceneDecorator, string) {
	i := rand.Int() % len(OptBoss2)
	return OptBoss2[i], OptBoss2Name[i]
}

func level1WavesDecor1(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil

	enemies.NewGatlingGhoul(ecs, 4, 3)
	combatscene.waves = append(combatscene.waves,
		level1Decorator5,
		level1Decorator15,
		level1Decorator6,
		level1Decorator7,
		level1Decorator13,
	)
}

func level1WavesDecor2(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil

	enemies.NewYeti(ecs, 4, 3)
	combatscene.waves = append(combatscene.waves,
		level1Decorator1,
		level1Decorator2,
		level2Decorator4,
		level2Decorator6,
		level1Decorator10,
	)
}
func level1WavesDecor3(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil

	enemies.NewLightningImp(ecs, 4, 1)
	combatscene.waves = append(combatscene.waves,
		level1Decorator1,
		level1Decorator2,
		level1Decorator19,
		level2Decorator2,
		level1Decorator20,
	)
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
		Col: 4,
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
	enemies.NewStunSpider(ecs, 5, 2)
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
func level1OptBoss4(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewMoltenSlug(ecs, 6, 1)
}

func level1Decorator10(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewPoacher(ecs, 6, 1)
	enemies.NewIceslime(ecs, 4, 2)
}
func level1Decorator11(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewPyroEyes(ecs, 4, 1)
	enemies.NewPoacher(ecs, 5, 3)
}
func level1Decorator12(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewLightningImp(ecs, 4, 1)
	enemies.NewBuzzer(ecs, 4, 1)
}
func level1Decorator13(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewGatlingGhoulOmega(ecs, 6, 1)
}
func level1Decorator14(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewStunSpider(ecs, 4, 1)
}
func level1Decorator15(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewBuzzer(ecs, 6, 1)
}
func level1Decorator16(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewPoacher(ecs, 6, 2)
	enemies.NewInfernoReaper(ecs, 6, 1)
}
func level1Decorator17(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	enemies.NewStunSpider(ecs, 6, 2)
	enemies.NewHealslime(ecs, 6, 1)
}
func level1Decorator18(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	LoadBomb(ecs.World, BoulderParam{Col: 5, Row: 0})
	LoadBomb(ecs.World, BoulderParam{Col: 2, Row: 2})
	enemies.NewPyroEyes(ecs, 6, 1)
}
func level1Decorator19(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	// LoadBomb(ecs.World, BoulderParam{Col: 5, Row: 0})
	// LoadBomb(ecs.World, BoulderParam{Col: 2, Row: 2})
	enemies.NewIceslime(ecs, 6, 1)
}
func level1Decorator20(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgMountain
	combatscene.rewards = nil
	LoadBomb(ecs.World, BoulderParam{Col: 2, Row: 2})
	enemies.NewLightningMage(ecs, 5, 3)
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
	enemies.NewInfernoReaper(ecs, 6, 1)
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
	enemies.NewBlazeBuzzer(ecs, 4, 2)
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
	enemies.NewGatlingGhoulOmega(ecs, 6, 2)
}
func level2Decorator9(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewInfernoReaper(ecs, 6, 1)
	enemies.NewPoacher(ecs, 4, 1)
}
func level2Decorator10(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewHealslime(ecs, 6, 1)
	enemies.NewStunSpider(ecs, 4, 1)
}
func level2WaveDecor1(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = []ItemInterface{
		&Medkit{},
	}
	for i := 0; i < 2; i++ {
		var temp loadout.Caster
		for {
			temp = DecorateCaster(GenerateCaster())
			if temp != nil {
				break
			}
		}
		combatscene.rewards = append(combatscene.rewards, temp)
	}

	hazard.NewRotatingFlame(ecs, 0, 0)
	hazard.NewRotatingFlame(ecs, 7, 3)
	enemies.NewBlizzBuzzer(ecs, 5, 1)
	enemies.NewReaper(ecs, 6, 1)
	combatscene.waves = append(combatscene.waves,
		level1Decorator1,
		level1Decorator2,
		level1Decorator14,
		level2Decorator2,
		level1Decorator10,
	)
}
func level2Decorator11(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewBlower(ecs, 6, 1)
	enemies.NewPoacher(ecs, 4, 1)
}
func level2Decorator12(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	hazard.NewRotatingFlame(ecs, 7, 1)
	hazard.NewRotatingFlame(ecs, 0, 0)
	enemies.NewInfernoReaper(ecs, 6, 1)
}
func level2OptBoss1(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewNyaaito(ecs, 6, 1)
}
func level2OptBoss2(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewFlameYeti(ecs, 6, 1)
}
func level3BossRush1(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil

	enemies.NewFrostYeti(ecs, 4, 1)
	combatscene.waves = append(combatscene.waves,
		level1OptBoss2,
	)
}
func Landon(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgForrest
	combatscene.rewards = nil
	enemies.NewLandon(ecs, 6, 1)
}
func finalBoss(ecs *ecs.ECS, combatscene *CombatScene) {
	combatscene.data.Bg = assets.BgCave
	combatscene.rewards = nil
	enemies.NewWhiteSnake(ecs, 6, 2)

}
