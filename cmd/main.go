package main

import (
	"log"

	"github.com/kharism/grimoiregunner/scene"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
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
	Level := scene.GenerateLayout1()
	state := &scene.SceneData{
		Bg:            assets.BgForrest,
		PlayerHP:      1000,
		PlayerMaxHP:   1000,
		PlayerCurrEn:  300,
		PlayerMaxEn:   300,
		PlayerEnRegen: 20,
		MainLoadout: []loadout.Caster{
			attack.NewHealCaster(),
			attack.NewWideSwordCaster(),
		},
		PlayerRow:    1,
		PlayerCol:    1,
		Level:        1,
		World:        donburi.NewWorld(),
		LevelLayout:  Level,
		CurrentLevel: Level.Root,
		SceneDecor:   Level.Root.Decorator, //scene.RandDecorator(),
		SubLoadout1:  []loadout.Caster{nil, nil},
		SubLoadout2:  []loadout.Caster{nil, nil},
		Inventory:    []scene.ItemInterface{},
	}
	combatScene := &scene.CombatScene{}
	// rewardScene := &scene.RewardScene{}
	ruleSet := map[stagehand.Scene[*scene.SceneData]][]stagehand.Directive[*scene.SceneData]{
		combatScene: {
			stagehand.Directive[*scene.SceneData]{Dest: scene.RewardSceneInstance, Trigger: scene.TriggerToReward},
			stagehand.Directive[*scene.SceneData]{Dest: scene.InventorySceneInstance, Trigger: scene.TriggerToInventory},
			stagehand.Directive[*scene.SceneData]{Dest: scene.StageSelectInstance, Trigger: scene.TriggerToStageSelect},
		},
		scene.RewardSceneInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: combatScene, Trigger: scene.TriggerToCombat},
		},
		scene.InventorySceneInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: combatScene, Trigger: scene.TriggerToCombat},
		},
		scene.StageSelectInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: combatScene, Trigger: scene.TriggerToCombat},
		},
	}
	manager := stagehand.NewSceneDirector[*scene.SceneData](combatScene, state, ruleSet)
	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
