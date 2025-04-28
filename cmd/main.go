package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/kharism/grimoiregunner/scene"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

// return width and height of the scene
func (g *Game) GetLayout() (width, height int) {
	return 1024, 600
}

// return the starting text position where the box containing name of the character appear on the scene
// return negative number if no such box needed
func (g *Game) GetNamePosition() (x, y int) {
	return 128, 600 - 150
}

// get the starting position of the text
func (g *Game) GetTextPosition() (x, y int) {
	return 128, 600 - 120
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GrimoireGunner")
	Level := scene.GenerateLayout1()

	state := &scene.SceneData{
		Bg:            assets.BgMountain,
		PlayerHP:      1000,
		PlayerMaxHP:   1000,
		PlayerCurrEn:  300,
		PlayerMaxEn:   300,
		PlayerEnRegen: 20,
		MainLoadout: []loadout.Caster{
			attack.NewShotgunCaster(),
			attack.NewCannonCaster(),
		},
		PlayerRow:    1,
		PlayerCol:    1,
		Level:        1,
		World:        nil,
		LevelLayout:  Level,
		CurrentLevel: Level.Root,
		// SceneDecor:   scene.,
		SubLoadout1: []loadout.Caster{nil, nil},
		SubLoadout2: []loadout.Caster{nil, nil},
		Inventory:   []scene.ItemInterface{
			// attack.NewCannonCaster(),
			// attack.NewHealCaster(),
		},
	}
	openingScene := scene.NewHanashiScene(scene.Scene1(&Game{}))
	endLevel1Scene := scene.NewHanashiScene(scene.Scene2(&Game{}))
	core.DetectKeyboardNext = func() bool {
		return inpututil.IsKeyJustReleased(ebiten.KeyQ)
	}

	combatScene := &scene.CombatScene{}
	// rewardScene := &scene.RewardScene{}
	ruleSet := map[stagehand.Scene[*scene.SceneData]][]stagehand.Directive[*scene.SceneData]{
		combatScene: {
			stagehand.Directive[*scene.SceneData]{Dest: scene.RewardSceneInstance, Trigger: scene.TriggerToReward},
			stagehand.Directive[*scene.SceneData]{Dest: scene.InventorySceneInstance, Trigger: scene.TriggerToInventory},
			stagehand.Directive[*scene.SceneData]{Dest: scene.StageSelectInstance, Trigger: scene.TriggerToStageSelect},
			stagehand.Directive[*scene.SceneData]{Dest: scene.MainMenuInstance, Trigger: scene.TriggerToMain},
			stagehand.Directive[*scene.SceneData]{Dest: scene.GameClearInstance, Trigger: scene.TriggerToClear},
			stagehand.Directive[*scene.SceneData]{Dest: endLevel1Scene, Trigger: scene.TriggerToPostLv1Story},
		},
		openingScene: {
			stagehand.Directive[*scene.SceneData]{Dest: scene.MainMenuInstance, Trigger: scene.TriggerToMain},
		},
		endLevel1Scene: {
			stagehand.Directive[*scene.SceneData]{Dest: scene.StageSelectInstance, Trigger: scene.TriggerToStageSelect},
		},
		scene.MainMenuInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: combatScene, Trigger: scene.TriggerToCombat},
			stagehand.Directive[*scene.SceneData]{Dest: openingScene, Trigger: scene.TriggerToOpening},
		},
		scene.GameClearInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: scene.MainMenuInstance, Trigger: scene.TriggerToMain},
		},
		scene.RewardSceneInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: combatScene, Trigger: scene.TriggerToCombat},
		},
		scene.InventorySceneInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: combatScene, Trigger: scene.TriggerToCombat},
		},
		scene.StageSelectInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: combatScene, Trigger: scene.TriggerToCombat},
			stagehand.Directive[*scene.SceneData]{Dest: scene.RestSceneInstance, Trigger: scene.TriggerToRest},
			stagehand.Directive[*scene.SceneData]{Dest: scene.WorkshopSceneInstance, Trigger: scene.TriggerToShop},
		},
		scene.RestSceneInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: combatScene, Trigger: scene.TriggerToCombat},
		},
		scene.WorkshopSceneInstance: {
			stagehand.Directive[*scene.SceneData]{Dest: combatScene, Trigger: scene.TriggerToCombat},
		},
	}
	manager := stagehand.NewSceneDirector[*scene.SceneData](scene.MainMenuInstance, state, ruleSet)
	openingScene.SetDoneFunc(func() {
		manager.ProcessTrigger(scene.TriggerToMain)
	})
	endLevel1Scene.SetDoneFunc(func() {
		manager.ProcessTrigger(scene.TriggerToStageSelect)
	})
	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
