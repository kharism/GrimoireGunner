package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/system"
	"github.com/yohamta/donburi"
)

type SceneData struct {
	PlayerHP      int
	PlayerMaxHP   int
	PlayerCurrEn  int
	PlayerMaxEn   int
	PlayerEnRegen int
	PlayerRow     int
	PlayerCol     int
	World         donburi.World

	MainLoadout []system.Caster
	SubLoadout1 []system.Caster
	SubLoadout2 []system.Caster

	Inventory []ItemInterface

	Level int //the difficulties

	LevelLayout  *Level
	SceneDecor   CombatSceneDecorator
	CurrentLevel *LevelNode

	Bg *ebiten.Image
}
