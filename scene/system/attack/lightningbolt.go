package attack

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type LightnigAtkParam struct {
	StartRow, StartCol int
	//+1 to create lightning on right of starting point, and -1 to create on left
	Direction int
	Actor     *donburi.Entry
}

func LightningBoltOnHitfunc(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	damage := component.Damage.Get(projectile).Damage
	health := component.Health.Get(receiver)
	health.HP -= damage
	ecs.World.Remove(projectile.Entity())
}
func NewLigtningAttack(ecs *ecs.ECS, param LightnigAtkParam) {
	startCol := param.StartCol
	damage := 60
	for i := startCol; i >= 0 && i <= 7; i += param.Direction {
		entity := ecs.World.Create(
			component.Damage,
			component.GridPos,
			component.OnHit,
			component.Transient,
			component.Fx,
		)
		entry := ecs.World.Entry(entity)
		component.Damage.Set(entry, &component.DamageData{Damage: damage})
		component.GridPos.Set(entry, &component.GridPosComponentData{Row: param.StartRow, Col: i})
		scrX, scrY := assets.GridCoord2Screen(param.StartRow, i)
		fxHeight := assets.LightningBolt.Bounds().Dy()
		anim1 := core.NewMovableImage(assets.LightningBolt, core.
			NewMovableImageParams().
			WithMoveParam(core.MoveParam{Sx: scrX - (float64(assets.TileWidth) / 2), Sy: scrY - float64(fxHeight)}))
		anim1.Done = func() {}
		component.Fx.Set(entry, &component.FxData{Animation: anim1})
		component.Transient.Set(entry, &component.TransientData{
			Start:    time.Now(),
			Duration: 200 * time.Millisecond,
		})
		component.OnHit.SetValue(entry, LightningBoltOnHitfunc)
	}
}
