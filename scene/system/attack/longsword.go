package attack

import (
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func NewLongSwordAttack(ecs *ecs.ECS, playerScrLoc component.ScreenPosComponentData, playerGridLoc component.GridPosComponentData) {
	param := DamageGridParam{}
	loc1 := &component.GridPosComponentData{Row: playerGridLoc.Row, Col: playerGridLoc.Col + 1}
	loc2 := &component.GridPosComponentData{Row: playerGridLoc.Row, Col: playerGridLoc.Col + 2}
	param.Locations = []*component.GridPosComponentData{loc1, loc2}
	param.Damage = []int{20, 20}
	param.OnHit = SingleHitProjectile

	param.Fx = assets.NewSwordAtkAnim(assets.SpriteParam{
		ScreenX: playerScrLoc.X + float64(assets.TileWidth)/2,
		ScreenY: playerScrLoc.Y - float64(assets.TileHeight),
		Modulo:  2,
	})
	NonProjectileAtk(ecs, param)
}

type DamageGridParam struct {
	Locations []*component.GridPosComponentData
	Damage    []int
	Fx        *core.AnimatedImage
	OnHit     component.OnAtkHit
}

func NonProjectileAtk(ecs *ecs.ECS, param DamageGridParam) {
	removeList := []donburi.Entity{}
	for idx, l := range param.Locations {
		entity := ecs.World.Create(
			component.Damage,
			component.GridPos,
			component.OnHit,
		)
		entry := ecs.World.Entry(entity)
		removeList = append(removeList, entity)
		component.Damage.Set(entry, &component.DamageData{param.Damage[idx]})
		component.GridPos.Set(entry, l)
		component.OnHit.Set(entry, &param.OnHit)
	}

	entity := ecs.World.Create(component.Fx)
	entry := ecs.World.Entry(entity)
	param.Fx.Done = func() {
		ecs.World.Remove(removeList[0])
		ecs.World.Remove(removeList[1])
		ecs.World.Remove(entity)
	}
	component.Fx.Set(entry, &component.FxData{param.Fx})
}
