package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
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

	MainLoadout []loadout.Caster
	SubLoadout1 []loadout.Caster
	SubLoadout2 []loadout.Caster

	Inventory []ItemInterface

	Level int //the difficulties

	LevelLayout  *Level
	SceneDecor   CombatSceneDecorator
	CurrentLevel *LevelNode

	//predetermined rewards
	rewards []ItemInterface

	Bg *ebiten.Image
}
