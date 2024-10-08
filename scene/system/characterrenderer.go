package system

import (
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type characterRenderer struct {
	query   *donburi.OrderedQuery[component.GridPosComponentData]
	counter uint // this to help with blinking function
}

// Render anything that has sprite and grid position
var CharacterRenderer = &characterRenderer{
	query: donburi.NewOrderedQuery[component.GridPosComponentData](
		filter.Contains(
			component.Sprite,
			component.GridPos,
		),
	),
}

func (r *characterRenderer) DrawCharacter(ecs *ecs.ECS, screen *ebiten.Image) {
	entries := []*donburi.Entry{}
	r.counter += 1
	r.query.Each(ecs.World, func(e *donburi.Entry) {
		entries = append(entries, e)
	})
	sort.Slice(entries, func(i, j int) bool {
		gridPosI := component.GridPos.Get(entries[i])
		gridPosJ := component.GridPos.Get(entries[j])
		return gridPosI.Order() < gridPosJ.Order()
	})
	for _, e := range entries {
		gridPos := component.GridPos.Get(e)

		// fmt.Println(e.Entity(), gridPos.Col, gridPos.Order(), component.Health.Get(e).Name)
		screenPos := component.ScreenPos.Get(e)
		if screenPos.X == 0 && screenPos.Y == 0 {
			screenPos.X = assets.TileStartX + float64(gridPos.Col)*float64(assets.TileWidth)
			screenPos.Y = assets.TileStartY + float64(gridPos.Row)*float64(assets.TileHeight)
		}
		blink := false
		if e.HasComponent(component.Health) {
			health := component.Health.Get(e)
			invisTime := health.InvisTime
			if !invisTime.IsZero() && invisTime.After(time.Now()) && r.counter%20 >= 15 && r.counter%20 <= 19 {
				blink = true
			}
		}
		if blink {
			continue
		}
		sprite := component.Sprite.Get(e).Image
		bound := sprite.Bounds()
		translate := ebiten.GeoM{}
		translate.Translate(-float64(bound.Dx())/2, -float64(bound.Dy()))
		translate.Translate(screenPos.X, screenPos.Y)
		drawOption := &ebiten.DrawImageOptions{
			GeoM: translate,
		}
		screen.DrawImage(sprite, drawOption)
	}
}
