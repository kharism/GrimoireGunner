package scene

import (
	mycomponent "github.com/kharism/mmbn_clone/scene/component"
	"github.com/kharism/mmbn_clone/scene/layers"
	"github.com/kharism/mmbn_clone/scene/system"

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
}
func (s *CombatScene) Load(state SceneData, manager stagehand.SceneController[SceneData]) {
	// your load code
	s.world = donburi.NewWorld()
	s.entitygrid = [4][8]int64{}
	s.ecs = ecs.NewECS(s.world)
	//add tiles entity
	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			s.world.Create(mycomponent.ScreenPos, mycomponent.GridPos)
		}
	}
	s.ecs.AddRenderer(layers.LayerBackground, system.DrawGrid)

	s.sm = manager.(*stagehand.SceneDirector[SceneData]) // This type assertion is important
}
func (s *CombatScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *CombatScene) Unload() SceneData {
	// your unload code
	return s.data
}
