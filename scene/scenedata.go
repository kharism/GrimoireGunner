package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/system"
)

type SceneData struct {
	PlayerHP      int
	PlayerMaxHP   int
	PlayerCurrEn  int
	PlayerMaxEn   int
	PlayerEnRegen int
	PlayerRow     int
	PlayerCol     int

	MainLoadout []system.Caster
	SubLoadout1 []system.Caster
	SubLoadout2 []system.Caster

	Level int //the difficulties

	SceneDecor CombatSceneDecorator

	Bg *ebiten.Image
}
