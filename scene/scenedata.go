package scene

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi"
)

type SceneData struct {
	PlayerHP       int
	PlayerMaxHP    int
	PlayerCurrEn   int
	PlayerMaxEn    int
	PlayerEnRegen  int
	PlayerRow      int
	PlayerCol      int
	World          donburi.World
	HanashiChoices map[string]any

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

	MusicSeek time.Duration

	Bg *ebiten.Image
}
