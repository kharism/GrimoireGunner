package system

import (
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var SRC_VFX = "VFX"
var SRC_CHR = "CHR"

// impplements Animation interface, but also add the source of the Animation so we know whether
// it was from VFX or Character
type AnimationSource interface {
	component.Animation
	Source() string
	ScrPos() (float64, float64)
}

type AnimFromCharacterEntry struct {
	entry   *donburi.Entry
	counter int
}

// empty implementation
func (a *AnimFromCharacterEntry) Update() {
	// TODO: idle animation implementation
}
func (a *AnimFromCharacterEntry) Source() string {
	return SRC_CHR
}
func (a *AnimFromCharacterEntry) ScrPos() (float64, float64) {
	gridPos := component.GridPos.Get(a.entry)
	scrX, scrY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
	return scrX, scrY
}
func (a *AnimFromCharacterEntry) Draw(screen *ebiten.Image) {
	gridPos := component.GridPos.Get(a.entry)
	a.counter += 1
	screenPos := component.ScreenPos.Get(a.entry)
	if screenPos.X == 0 && screenPos.Y == 0 {
		screenPos.X = assets.TileStartX + float64(gridPos.Col)*float64(assets.TileWidth)
		screenPos.Y = assets.TileStartY + float64(gridPos.Row)*float64(assets.TileHeight)
	}
	blink := false
	if a.entry.HasComponent(component.Health) {
		health := component.Health.Get(a.entry)
		invisTime := health.InvisTime
		if !invisTime.IsZero() && invisTime.After(time.Now()) && a.counter%20 >= 15 && a.counter%20 <= 19 {
			blink = true
		}
	}
	if blink {
		return
	}
	sprite := component.Sprite.Get(a.entry).Image
	bound := sprite.Bounds()
	translate := ebiten.GeoM{}
	translate.Translate(-float64(bound.Dx())/2, -float64(bound.Dy()))
	translate.Translate(screenPos.X, screenPos.Y)
	if a.entry.HasComponent(component.Shader) {
		opts := &ebiten.DrawRectShaderOptions{
			GeoM: translate,
		}
		opts.Images[0] = sprite
		shader := component.Shader.Get(a.entry)
		screen.DrawRectShader(bound.Dx(), bound.Dy(), shader, opts)
	} else {

		drawOption := &ebiten.DrawImageOptions{
			GeoM: translate,
		}
		screen.DrawImage(sprite, drawOption)
	}

}
func AnimationSourceFromHP(entry *donburi.Entry) AnimationSource {
	return &AnimFromCharacterEntry{entry: entry, counter: 0}
}

type AnimFromVfxData struct {
	component.Animation
}

func (a *AnimFromVfxData) Source() string {
	return SRC_VFX
}
func (f *AnimFromVfxData) ScrPos() (float64, float64) {
	if jj, ok := f.Animation.(*core.MovableImage); ok {
		x, y := jj.GetPos()
		_, height := jj.GetSize()
		return x, y + height
	} else if jj, ok := f.Animation.(*core.AnimatedImage); ok {
		x, y := jj.GetPos()
		_, height := jj.GetSize()
		return x, y + height
	}
	return 100, 100
}
func UnifiedRenderer(ecs *ecs.ECS, screen *ebiten.Image) {
	drawables := []AnimationSource{}
	query := donburi.NewOrderedQuery[component.GridPosComponentData](
		filter.Contains(
			component.Sprite,
			component.GridPos,
		),
	)
	query.Each(ecs.World, func(e *donburi.Entry) {
		jj := AnimationSourceFromHP(e)
		drawables = append(drawables, jj)
	})
	component.Fx.Each(ecs.World, func(e *donburi.Entry) {
		fx := component.Fx.Get(e)
		drawables = append(drawables, &AnimFromVfxData{fx.Animation})
		// fx.Animation.Draw(screen)
	})
	sort.Slice(drawables, func(i, j int) bool {
		_, y1 := drawables[i].ScrPos()
		_, y2 := drawables[j].ScrPos()
		if y1 < y2 {
			return true
		}
		if y2 < y1 {
			return false
		}
		if drawables[i].Source() == SRC_CHR {
			return false
		}
		return false
		// return gridPosI.Order() < gridPosJ.Order()
	})
	for _, v := range drawables {
		v.Draw(screen)
	}
}
