package main

import (
	"log"

	"github.com/kharism/grimoiregunner/scene"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/system"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/yohamta/donburi"

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
	state := &scene.SceneData{
		Bg:            assets.BgForrest,
		PlayerHP:      1000,
		PlayerMaxHP:   1000,
		PlayerCurrEn:  300,
		PlayerMaxEn:   300,
		PlayerEnRegen: 20,
		MainLoadout: []system.Caster{
			attack.NewShotgunCaster(),
			attack.NewWideSwordCaster(),
		},
		PlayerRow:   1,
		PlayerCol:   1,
		Level:       1,
		World:       donburi.NewWorld(),
		SceneDecor:  scene.Level1Decorator1, //scene.RandDecorator(),
		SubLoadout1: []system.Caster{nil, nil},
		SubLoadout2: []system.Caster{nil, nil},
		Inventory:   []scene.ItemInterface{},
	}
	combatScene := &scene.CombatScene{}
	// rewardScene := &scene.RewardScene{}
	ruleSet := map[stagehand.Scene[*scene.SceneData]][]stagehand.Directive[*scene.SceneData]{
		combatScene: {
			stagehand.Directive[*scene.SceneData]{Dest: scene.RewardSceneInstance, Trigger: scene.TriggerToReward},
			stagehand.Directive[*scene.SceneData]{Dest: scene.InventorySceneInstance, Trigger: scene.TriggerToInventory},
		},
		scene.RewardSceneInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: combatScene, Trigger: scene.TriggerToCombat},
		},
		scene.InventorySceneInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: combatScene, Trigger: scene.TriggerToCombat},
		},
	}
	manager := stagehand.NewSceneDirector[*scene.SceneData](combatScene, state, ruleSet)
	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
