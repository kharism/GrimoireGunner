package attack

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type LobAtkParams struct {
	EntrySource          *donburi.Entry
	TargetRow, TargetCol int

	Sprite *ebiten.Image
	Damage int
}

func NewLobAttack(ecs_ *ecs.ECS, param LobAtkParams) {
	targetGrid := ecs_.World.Create(component.GridPos, component.GridTarget)
	ii := ecs_.World.Entry(targetGrid)
	component.GridPos.Set(ii, &component.GridPosComponentData{
		Row: param.TargetRow,
		Col: param.TargetCol,
	})
	// time.Sleep(2000 * time.Millisecond)
	sprite := component.Sprite.Get(param.EntrySource)
	sprite.Image = assets.BloombomberAtk
	topLimit := -110.0
	entityScreenPos := component.ScreenPos.Get(param.EntrySource)
	// ii := ecs.World.Entry(targetGrid)
	targetGridPos := component.GridPos.Get(ii)
	targetScreenPosX, targetScreenPosY := assets.GridCoord2Screen(targetGridPos.Row, targetGridPos.Col)
	anim1 := core.NewMovableImage(assets.Bomb1,
		core.NewMovableImageParams().
			WithMoveParam(
				core.MoveParam{
					Sx:    entityScreenPos.X - (float64(assets.TileWidth) / 2),
					Sy:    entityScreenPos.Y - (float64(assets.TileHeight + 25)),
					Speed: 5,
				}),
	)
	anim1.AnimationQueue = append(anim1.AnimationQueue,
		core.NewMoveAnimationFromParam(core.MoveParam{Tx: entityScreenPos.X - (float64(assets.TileWidth) / 2), Ty: topLimit, Speed: 5}),
		core.NewMoveAnimationFromParam(core.MoveParam{Tx: targetScreenPosX - (float64(assets.TileWidth) / 2), Ty: topLimit, Speed: 5}),
		core.NewMoveAnimationFromParam(core.MoveParam{Tx: targetScreenPosX - (float64(assets.TileWidth) / 2), Ty: targetScreenPosY - (float64(assets.TileHeight + 25)), Speed: 5}),
	)
	fxEntry := ecs_.World.Create(component.Fx)
	hh := ecs_.World.Entry(fxEntry)
	component.Fx.Set(hh, &component.FxData{anim1})
	anim1.Done = func() {
		ecs_.World.Remove(fxEntry)
		ecs_.World.Remove(targetGrid)
		sprite.Image = assets.Bloombomber
		dmgTile := ecs_.World.Create(
			component.GridPos,
			component.Damage,
			component.OnHit,
		)
		ent := ecs_.World.Entry(dmgTile)
		component.GridPos.Set(ent, &component.GridPosComponentData{Col: targetGridPos.Col, Row: targetGridPos.Row})
		component.Damage.Set(ent, &component.DamageData{Damage: param.Damage})
		component.OnHit.SetValue(ent, func(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
			damage := component.Damage.Get(projectile).Damage
			component.Health.Get(receiver).HP -= damage
			ecs.World.Remove(dmgTile)
		})
		query2 := donburi.NewQuery(filter.Contains(component.Health, component.GridPos))
		query2.Each(ecs_.World, func(e *donburi.Entry) {
			loc := component.GridPos.Get(e)
			if loc.Col == targetGridPos.Col && loc.Row == targetGridPos.Row {
				component.OnHit.GetValue(ent)(ecs_, ent, e)
			}
		})
		ecs_.World.Remove(dmgTile)

	}
}
