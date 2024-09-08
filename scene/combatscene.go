package scene

import (
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/layers"
	"github.com/kharism/grimoiregunner/scene/system"
	"github.com/kharism/grimoiregunner/scene/system/attack"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type CombatScene struct {
	data       SceneData
	sm         *stagehand.SceneDirector[SceneData]
	world      donburi.World
	ecs        *ecs.ECS
	debugPause bool
	// grid store entity id or 0 if no entity occupy the cell
	entitygrid [4][8]int64
}

func (c *CombatScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		c.debugPause = !c.debugPause
	}

	if c.debugPause {
		return nil
	}
	c.ecs.Update()
	return nil
}

func (c *CombatScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	c.ecs.DrawLayer(layers.LayerBackground, screen)
	c.ecs.DrawLayer(layers.LayerGrid, screen)
	c.ecs.DrawLayer(layers.LayerCharacter, screen)
	c.ecs.DrawLayer(layers.LayerFx, screen)
	c.ecs.DrawLayer(layers.LayerHP, screen)
	c.ecs.DrawLayer(layers.LayerDebug, screen)
	c.ecs.DrawLayer(layers.LayerUI, screen)
}
func LoadGrid(world donburi.World) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			idx := world.Create(mycomponent.ScreenPos, mycomponent.GridPos, mycomponent.TileTag)
			entId := world.Entry(idx)
			mycomponent.GridPos.Set(entId, &mycomponent.GridPosComponentData{Col: j, Row: i})

		}
	}
}
func LoadPlayer(world donburi.World) *donburi.Entity {
	playerEntity := archetype.NewPlayer(world, assets.Player1Stand)
	gridPos := component.GridPos.Get(world.Entry(*playerEntity))
	screenPos := component.ScreenPos.Get(world.Entry(*playerEntity))
	screenPos.X = assets.TileStartX + float64(gridPos.Col)*float64(assets.TileWidth)
	screenPos.Y = assets.TileStartY + float64(gridPos.Row)*float64(assets.TileHeight)
	return playerEntity
}

type BoulderParam struct {
	Col, Row int
}

func LoadBoulder(world donburi.World, param BoulderParam) *donburi.Entity {
	entity := world.Create(
		mycomponent.Health,
		mycomponent.GridPos,
		mycomponent.ScreenPos,
		mycomponent.Sprite,
		archetype.ConstructTag,
	)
	entry := world.Entry(entity)
	mycomponent.Sprite.Set(entry, &mycomponent.SpriteData{Image: assets.Boulder})
	mycomponent.Health.Set(entry, &mycomponent.HealthData{HP: 200, Name: "Boulder"})
	mycomponent.GridPos.Set(entry, &mycomponent.GridPosComponentData{Col: param.Col, Row: param.Row})
	return &entity
}
func (s *CombatScene) Load(state SceneData, manager stagehand.SceneController[SceneData]) {
	// your load code
	s.world = donburi.NewWorld()
	s.entitygrid = [4][8]int64{}
	s.ecs = ecs.NewECS(s.world)
	s.debugPause = false
	//add tiles entity
	LoadGrid(s.world)
	LoadBoulder(s.world, BoulderParam{
		Col: 5,
		Row: 1,
	})
	playerEntity := LoadPlayer(s.world)
	LoadBoulder(s.world, BoulderParam{
		Col: 2,
		Row: 2,
	})
	// enemies.NewCannoneer(s.ecs, 6, 1)
	// enemies.NewGatlingGhoul(s.ecs, 4, 1)
	// enemies.NewReaper(s.ecs, 4, 1)
	assets.Bg = state.Bg
	system.CurLoadOut[0] = attack.NewLightningBolCaster()
	system.CurLoadOut[1] = attack.NewShockwaveCaster() //attack.NewLongSwordCaster()

	system.SubLoadOut1[0] = attack.NewWideSwordCaster() //attack.NewLongSwordCaster()
	system.SubLoadOut1[1] = attack.NewBuckshotCaster()
	system.SubLoadOut2[0] = attack.NewFirewallCaster()
	system.SubLoadOut2[1] = attack.NewGatlingCastor()

	Ensystemrenderer := system.EnergySystem
	eq := system.EventQueueSystem{}
	// attack.GenerateMagibullet(s.ecs, 1, 5, -15)
	s.ecs.
		AddSystem(system.NewPlayerMoveSystem(playerEntity).Update).
		AddSystem(system.DamageSystem.Update).
		AddSystem(system.NPMoveSystem.Update).
		AddSystem(system.NewPlayerAttackSystem(playerEntity).Update).
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

	s.sm = manager.(*stagehand.SceneDirector[SceneData]) // This type assertion is important
}
func (s *CombatScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *CombatScene) Unload() SceneData {
	// your unload code
	return s.data
}
