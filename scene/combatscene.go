package scene

import (
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/layers"
	"github.com/kharism/grimoiregunner/scene/system"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/grimoiregunner/scene/system/enemies"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type CombatScene struct {
	data  SceneData
	sm    *stagehand.SceneDirector[SceneData]
	world donburi.World
	ecs   *ecs.ECS
	// grid store entity id or 0 if no entity occupy the cell
	entitygrid [4][8]int64
}

func (c *CombatScene) Update() error {
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
	enemies.NewCannoneer(s.ecs, 6, 1)
	assets.Bg = state.Bg

	attack.GenerateMagibullet(s.ecs, 1, 5, -15)
	s.ecs.
		AddSystem(system.NewPlayerMoveSystem(playerEntity).Update).
		AddSystem(system.DamageSystem.Update).
		AddSystem(system.NPMoveSystem.Update).
		AddSystem(system.NewPlayerAttackSystem(playerEntity).Update).
		AddSystem(system.NewTransientSystem().Update).
		AddSystem(system.UpdateFx).
		AddSystem(system.EnemyAI.Update).
		AddRenderer(layers.LayerBackground, system.DrawBg).
		AddRenderer(layers.LayerGrid, system.GridRenderer.DrawGrid).
		AddRenderer(layers.LayerCharacter, system.CharacterRenderer.DrawCharacter).
		AddRenderer(layers.LayerFx, system.RenderFx).
		AddRenderer(layers.LayerDebug, system.DebugRenderer.DrawDebug).
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
