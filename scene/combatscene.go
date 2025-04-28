package scene

import (
	"fmt"
	"math/rand"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/events"
	"github.com/kharism/grimoiregunner/scene/layers"
	"github.com/kharism/grimoiregunner/scene/system"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/grimoiregunner/scene/system/enemies"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type CombatScene struct {
	data       *SceneData
	sm         *stagehand.SceneDirector[*SceneData]
	world      donburi.World
	ecs        *ecs.ECS
	waves      []func(*ecs.ECS, *CombatScene)
	debugPause bool

	sandboxMode bool
	// grid store entity id or 0 if no entity occupy the cell
	entitygrid  [4][8]int64
	musicPlayer *assets.AudioPlayer
	rewards     []ItemInterface
	loopMusic   bool
}

func (c *CombatScene) Update() error {
	if c.loopMusic && !c.musicPlayer.AudioPlayer().IsPlaying() {
		c.musicPlayer.AudioPlayer().Rewind()
		c.musicPlayer.AudioPlayer().Play()
	}
	if c.musicPlayer != nil {
		c.musicPlayer.Update()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		c.debugPause = !c.debugPause
	}
	if c.sandboxMode && inpututil.IsKeyJustPressed(ebiten.KeyI) {
		c.sm.ProcessTrigger(TriggerToInventory)
	}
	if c.sandboxMode && inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		defTrigger := TriggerToStageSelect
		if len(c.data.CurrentLevel.NextNode) == 0 {
			c.data.Level += 1
			// generate new level
			switch c.data.Level {
			case 2:
				c.data.LevelLayout = GenerateLayout2()
				defTrigger = TriggerToPostLv1Story
			}

			c.data.CurrentLevel = nil //c.data.LevelLayout.Root
		}
		c.sm.ProcessTrigger(defTrigger)
	}
	if c.debugPause {
		return nil
	}
	c.ecs.Update()
	return nil
}

func init() {
	MonogramFace = &text.GoTextFace{
		Source: assets.MonogramFont,
		Size:   25,
	}
}

var MonogramFace *text.GoTextFace

func (c *CombatScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	c.ecs.DrawLayer(layers.LayerBackground, screen)
	c.ecs.DrawLayer(layers.LayerGrid, screen)
	c.ecs.DrawLayer(layers.LayerCharacter, screen)
	c.ecs.DrawLayer(layers.LayerFx, screen)
	c.ecs.DrawLayer(layers.LayerHP, screen)

	c.ecs.DrawLayer(layers.LayerUI, screen)
	c.ecs.DrawLayer(layers.LayerDebug, screen)

	if c.sandboxMode {

		textTranslate := ebiten.GeoM{}
		textTranslate.Translate(950, 50)

		textDrawOpt := text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: text.AlignEnd,
			},
			DrawImageOptions: ebiten.DrawImageOptions{
				GeoM: textTranslate,
			},
		}
		text.Draw(screen, "Press i for inventory", MonogramFace, &textDrawOpt)
		textTranslate.Translate(0, -30)
		textDrawOpt.DrawImageOptions.GeoM = textTranslate
		text.Draw(screen, "Press tab for map", MonogramFace, &textDrawOpt)
	}
}
func LoadGrid(world donburi.World) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			idx := world.Create(component.ScreenPos, component.GridPos, component.TileTag)
			entId := world.Entry(idx)
			component.GridPos.Set(entId, &component.GridPosComponentData{Col: j, Row: i})

		}
	}
}
func LoadPlayer(world donburi.World, state *SceneData) *donburi.Entity {
	playerEntity := archetype.NewPlayer(world, assets.Player1Stand)
	gridPos := component.GridPos.Get(world.Entry(*playerEntity))
	gridPos.Col = state.PlayerCol
	gridPos.Row = state.PlayerRow
	component.Health.Get(world.Entry(*playerEntity)).HP = state.PlayerHP
	component.Health.Get(world.Entry(*playerEntity)).MaxHP = state.PlayerMaxHP
	system.EnergySystem.SetEn(state.PlayerCurrEn)
	system.EnergySystem.MaxEN = state.PlayerMaxEn
	system.EnergySystem.ENRegen = state.PlayerEnRegen
	screenPos := component.ScreenPos.Get(world.Entry(*playerEntity))
	screenPos.X = assets.TileStartX + float64(gridPos.Col)*float64(assets.TileWidth)
	screenPos.Y = assets.TileStartY + float64(gridPos.Row)*float64(assets.TileHeight)
	return playerEntity
}

var RegisterCombatClear bool

// play fanfare
func (s *CombatScene) OnCombatClear() {
	if len(s.waves) == 0 {
		s.musicPlayer.AudioPlayer().Pause()
		s.musicPlayer, _ = assets.NewAudioPlayer(assets.Fanfare, assets.TypeMP3)
		s.loopMusic = false
		s.musicPlayer.AudioPlayer().Play()
		stgClrDim := assets.StageClear.Bounds()
		movableImg := core.NewMovableImage(assets.StageClear,
			core.NewMovableImageParams().WithMoveParam(core.MoveParam{
				Sx:    float64(-stgClrDim.Dx()),
				Sy:    float64(300 + stgClrDim.Dy()/2),
				Speed: 10}))
		movableImg.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{
			Tx:    float64(600 - stgClrDim.Dx()/2 - 60),
			Ty:    float64(300 + stgClrDim.Dy()/2),
			Speed: 10,
		}))
		movableImg.Done = func() {
			system.PlayerAttackSystem.State = system.CombatClearState
		}
		//turn off attack system
		system.PlayerAttackSystem.State = system.DoNothingState
		//attach the stageclear to fx system
		stgDone := s.ecs.World.Create(component.Anouncement)
		component.Anouncement.Set(s.ecs.World.Entry(stgDone), &component.FxData{
			Animation: movableImg,
		})
	} else {
		// a hack to keep the rewards/BG intact
		realReward := s.rewards
		realBg := s.data.Bg
		s.waves[0](s.ecs, s)
		s.waves = s.waves[1:]
		s.rewards = realReward
		s.data.Bg = realBg
	}
}
func (s *CombatScene) OnGameOver() {

}
func (s *CombatScene) Load(state *SceneData, manager stagehand.SceneController[*SceneData]) {
	// your load code
	s.sm = manager.(*stagehand.SceneDirector[*SceneData]) // This type assertion is important
	s.world = donburi.NewWorld()
	s.entitygrid = [4][8]int64{}
	s.ecs = ecs.NewECS(s.world)
	s.data = state
	if !RegisterCombatClear {
		events.CombatClearEvent.Subscribe(s.world, func(w donburi.World, event events.CombatClearData) {
			if event.IsGameOver {
				RegisterCombatClear = false
				s.sm.ProcessTrigger(TriggerToMain)
			} else {
				RegisterCombatClear = false
				if len(s.data.CurrentLevel.NextNode) == 0 && s.data.Level == 2 {
					s.sm.ProcessTrigger(TriggerToClear)
				} else {
					s.sm.ProcessTrigger(TriggerToReward)
				}

			}

		})
		RegisterCombatClear = true
	}

	s.debugPause = false
	//add tiles entity
	LoadGrid(s.world)
	// LoadBoulder(s.world, BoulderParam{
	// 	Col: 5,
	// 	Row: 1,
	// })
	playerEntity := LoadPlayer(s.world, state)
	// LoadBoulder(s.world, BoulderParam{
	// 	Col: 2,
	// 	Row: 2,
	// })
	// enemies.NewCannoneer(s.ecs, 6, 1)
	// enemies.NewGatlingGhoul(s.ecs, 4, 1)
	// enemies.NewReaper(s.ecs, 4, 1)

	loadout.CurLoadOut[0] = state.MainLoadout[0]
	loadout.CurLoadOut[1] = state.MainLoadout[1] //attack.NewLongSwordCaster()

	loadout.SubLoadOut1[0] = state.SubLoadout1[0]
	loadout.SubLoadOut1[1] = state.SubLoadout1[1]
	loadout.SubLoadOut2[0] = state.SubLoadout2[0]
	loadout.SubLoadOut2[1] = state.SubLoadout2[1]
	if state.CurrentLevel != nil && state.CurrentLevel.SelectedStage != nil {
		state.CurrentLevel.SelectedStage.DecorSceneData(s.data)
	}
	if state.SceneDecor != nil {
		state.SceneDecor(s.ecs, s)
		s.sandboxMode = false
		s.musicPlayer = nil
	} else {
		s.sandboxMode = true
		s.musicPlayer = nil
	}
	assets.Bg = state.Bg

	Ensystemrenderer := system.EnergySystem
	eq := system.EventQueueSystem{}
	system.PlayerAttackSystem.PlayerIndex = playerEntity
	system.PlayerAttackSystem.State = system.CombatState
	s.loopMusic = true
	if s.musicPlayer == nil {
		var err error
		if s.sandboxMode {
			s.musicPlayer, err = assets.NewAudioPlayer(assets.IntermissionMusic, assets.TypeMP3)
			s.musicPlayer.AudioPlayer().SetPosition(s.data.MusicSeek)
		} else {
			jj := rand.Int() % 4
			switch jj {
			case 0:
				s.musicPlayer, err = assets.NewAudioPlayer(assets.BattleMusic1, assets.TypeMP3)
			case 1:
				s.musicPlayer, err = assets.NewAudioPlayer(assets.BattleMusic2, assets.TypeMP3)
			case 2:
				s.musicPlayer, err = assets.NewAudioPlayer(assets.BattleMusic3, assets.TypeMP3)
			case 3:
				s.musicPlayer, err = assets.NewAudioPlayer(assets.BattleMusic4, assets.TypeMP3)
			}

		}

		if err != nil {
			fmt.Println(err.Error())
		}
		s.musicPlayer.AudioPlayer().Play()
		// set interfaces for sfx
		attack.AtkSfxQueue = s.musicPlayer
		enemies.EnemySfxQueue = s.musicPlayer
	} else {
		// s.musicPlayer.audioPlayer.Rewind()
		s.musicPlayer.AudioPlayer().Play()
	}
	system.DamageSystem.DamageEventConsumer = s
	// attack.GenerateMagibullet(s.ecs, 1, 5, -15)
	s.ecs.
		AddSystem(system.NewPlayerMoveSystem(playerEntity).Update).
		AddSystem(system.Kamikaze.Update).
		AddSystem(system.UpdateBurnSystem).
		AddSystem(system.DamageSystem.Update).
		AddSystem(system.NPMoveSystem.Update).
		AddSystem(system.PlayerAttackSystem.Update).
		AddSystem(system.NewTransientSystem().Update).
		AddSystem(Ensystemrenderer.Update).
		AddSystem(eq.Update).
		AddSystem(system.EnemyAI.Update).
		AddSystem(system.UpdateFx).
		AddSystem(system.UpdateAnouncement).
		AddRenderer(layers.LayerBackground, system.DrawBg).
		AddRenderer(layers.LayerGrid, system.GridRenderer.DrawGrid).
		// AddRenderer(layers.LayerCharacter, system.CharacterRenderer.DrawCharacter).
		// AddRenderer(layers.LayerFx, system.RenderFx).
		AddRenderer(layers.LayerCharacter, system.UnifiedRenderer).
		// AddRenderer(layers.LayerDebug, system.DebugRenderer.DrawDebug).
		AddRenderer(layers.LayerUI, Ensystemrenderer.DrawEnBar).
		AddRenderer(layers.LayerUI, system.RenderLoadOut).
		AddRenderer(layers.LayerUI, system.RenderAnouncement).
		AddRenderer(layers.LayerHP, system.HPRenderer.DrawHP)

}
func OnCombatClear(w donburi.World, event events.CombatClearData) {

}
func (s *CombatScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *CombatScene) Unload() *SceneData {
	// your unload code
	s.data.MainLoadout = loadout.CurLoadOut[:]
	s.data.rewards = s.rewards
	if s.data.MainLoadout[0] != nil {
		s.data.MainLoadout[0].ResetCooldown()
	}
	if s.data.MainLoadout[1] != nil {
		s.data.MainLoadout[1].ResetCooldown()
	}
	s.data.SubLoadout1 = loadout.SubLoadOut1[:]
	if s.data.SubLoadout1[0] != nil {
		s.data.SubLoadout1[0].ResetCooldown()
	}
	if s.data.SubLoadout1[1] != nil {
		s.data.SubLoadout1[1].ResetCooldown()
	}
	s.data.SubLoadout2 = loadout.SubLoadOut2[:]
	if s.data.SubLoadout2[0] != nil {
		s.data.SubLoadout2[0].ResetCooldown()
	}
	if s.data.SubLoadout2[1] != nil {
		s.data.SubLoadout2[1].ResetCooldown()
	}
	s.loopMusic = false
	player, ok := archetype.PlayerTag.First(s.ecs.World)
	if ok {
		pp := component.Health.Get(player)
		s.data.PlayerHP = pp.HP
		s.data.PlayerMaxHP = pp.MaxHP
	}

	s.data.MusicSeek = s.musicPlayer.AudioPlayer().Position()
	s.musicPlayer.AudioPlayer().Rewind()
	s.musicPlayer.AudioPlayer().Pause()
	return s.data
}
