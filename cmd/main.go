package main

import (
	"log"

	"github.com/kharism/grimoiregunner/scene"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/system"
	"github.com/kharism/grimoiregunner/scene/system/attack"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

const (
	screenWidth  = 1024
	screenHeight = 600
)

type Game struct {
	count int
}

func (g *Game) Update() error {
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {

}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GrimoireGunner")
	state := scene.SceneData{
		Bg:            assets.BgForrest,
		PlayerHP:      1000,
		PlayerMaxHP:   1000,
		PlayerCurrEn:  300,
		PlayerMaxEn:   300,
		PlayerEnRegen: 20,
		MainLoadout: []system.Caster{
			attack.NewGatlingCastor(),
			attack.NewWideSwordCaster(),
		},
		PlayerRow:   1,
		PlayerCol:   1,
		Level:       1,
		SceneDecor:  scene.RandDecorator(),
		SubLoadout1: []system.Caster{nil, nil},
		SubLoadout2: []system.Caster{nil, nil},
	}
	combatScene := &scene.CombatScene{}
	ruleSet := map[stagehand.Scene[scene.SceneData]][]stagehand.Directive[scene.SceneData]{}
	manager := stagehand.NewSceneDirector[scene.SceneData](combatScene, state, ruleSet)
	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
