package scene

import (
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/events"
	"github.com/kharism/grimoiregunner/scene/layers"
	"github.com/kharism/grimoiregunner/scene/system"
	"github.com/kharism/grimoiregunner/scene/system/loadout"

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
	debugPause bool

	sandboxMode bool
	// grid store entity id or 0 if no entity occupy the cell
	entitygrid [4][8]int64
}

func (c *CombatScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		c.debugPause = !c.debugPause
	}
	if c.sandboxMode && inpututil.IsKeyJustPressed(ebiten.KeyI) {
		c.sm.ProcessTrigger(TriggerToInventory)
	}
	if c.sandboxMode && inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		c.sm.ProcessTrigger(TriggerToStageSelect)
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
	c.ecs.DrawLayer(layers.LayerDebug, screen)
	c.ecs.DrawLayer(layers.LayerUI, screen)

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

func (s *CombatScene) Load(state *SceneData, manager stagehand.SceneController[*SceneData]) {
	// your load code
	s.sm = manager.(*stagehand.SceneDirector[*SceneData]) // This type assertion is important
	s.world = donburi.NewWorld()
	s.entitygrid = [4][8]int64{}
	s.ecs = ecs.NewECS(s.world)
	s.data = state
	if !RegisterCombatClear {
		events.CombatClearEvent.Subscribe(s.world, func(w donburi.World, event events.CombatClearData) {
			RegisterCombatClear = false
			s.sm.ProcessTrigger(TriggerToReward)
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
	assets.Bg = state.Bg

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
	} else {
		s.sandboxMode = true
	}

	Ensystemrenderer := system.EnergySystem
	eq := system.EventQueueSystem{}
	system.PlayerAttackSystem.PlayerIndex = playerEntity
	system.PlayerAttackSystem.State = system.CombatState
	// attack.GenerateMagibullet(s.ecs, 1, 5, -15)
	s.ecs.
		AddSystem(system.NewPlayerMoveSystem(playerEntity).Update).
		AddSystem(system.DamageSystem.Update).
		AddSystem(system.NPMoveSystem.Update).
		AddSystem(system.PlayerAttackSystem.Update).
		AddSystem(system.NewTransientSystem().Update).
		AddSystem(Ensystemrenderer.Update).
		AddSystem(eq.Update).
		AddSystem(system.EnemyAI.Update).
		AddSystem(system.UpdateFx).
		AddRenderer(layers.LayerBackground, system.DrawBg).
		AddRenderer(layers.LayerGrid, system.GridRenderer.DrawGrid).
		AddRenderer(layers.LayerCharacter, system.CharacterRenderer.DrawCharacter).
		AddRenderer(layers.LayerFx, system.RenderFx).
		AddRenderer(layers.LayerDebug, system.DebugRenderer.DrawDebug).
		AddRenderer(layers.LayerUI, Ensystemrenderer.DrawEnBar).
		AddRenderer(layers.LayerUI, system.RenderLoadOut).
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
	s.data.SubLoadout1 = loadout.SubLoadOut1[:]
	s.data.SubLoadout2 = loadout.SubLoadOut2[:]
	return s.data
}
