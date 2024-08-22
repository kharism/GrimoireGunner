package enemies

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func NewBloombomber(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Bloombomber)
	entry := ecs.World.Entry(*entity)
	component.Health.Set(entry, &component.HealthData{HP: 140, Name: "Cannoneer"})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[LAST_FIRED] = nil

	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: BloombomberRoutine, Memory: data})
}

const LAST_FIRED = "LAST_FIRED"

func BloombomberRoutine(ecs_ *ecs.ECS, entity *donburi.Entry) {

	// playerScreenPos := component.ScreenPos.Get(player)
	// Damage := 20
	memory := component.EnemyRoutine.Get(entity).Memory
	// sprite := component.Sprite.Get(entity)
	curTime := time.Now()
	// var lastFired *time.Time
	if _, ok := memory[LAST_FIRED].(time.Time); !ok {
		memory[LAST_FIRED] = curTime
		return
	}
	lastFired := memory[LAST_FIRED].(time.Time)
	kk := curTime.Sub(lastFired)
	vv := 2 * time.Second
	dist := kk <= vv
	if dist {
		return
	}
	player, _ := archetype.PlayerTag.First(ecs_.World)
	playerGridPos := component.GridPos.GetValue(player)
	memory[LAST_FIRED] = time.Now()
	bombParam := attack.LobAtkParams{
		EntrySource: entity,
		TargetRow:   playerGridPos.Row,
		TargetCol:   playerGridPos.Col,
		Sprite:      assets.Bomb1,
		Damage:      20,
	}
	attack.NewLobAttack(ecs_, bombParam)
	// targetGrid := ecs_.World.Create(component.GridPos, component.GridTarget)
	// ii := ecs_.World.Entry(targetGrid)
	// component.GridPos.Set(ii, &component.GridPosComponentData{
	// 	Row: playerGridPos.Row,
	// 	Col: playerGridPos.Col,
	// })
	// sprite := component.Sprite.Get(entity)
	// sprite.Image = assets.BloombomberAtk
	// topLimit := -110.0
	// entityScreenPos := component.ScreenPos.Get(entity)
	// // ii := ecs.World.Entry(targetGrid)
	// targetGridPos := component.GridPos.Get(ii)
	// targetScreenPosX, targetScreenPosY := assets.GridCoord2Screen(targetGridPos.Row, targetGridPos.Col)
	// anim1 := core.NewMovableImage(assets.Bomb1,
	// 	core.NewMovableImageParams().
	// 		WithMoveParam(
	// 			core.MoveParam{
	// 				Sx:    entityScreenPos.X - (float64(assets.TileWidth) / 2),
	// 				Sy:    entityScreenPos.Y - (float64(assets.TileHeight + 25)),
	// 				Speed: 5,
	// 			}),
	// )
	// anim1.AnimationQueue = append(anim1.AnimationQueue,
	// 	core.NewMoveAnimationFromParam(core.MoveParam{Tx: entityScreenPos.X - (float64(assets.TileWidth) / 2), Ty: topLimit, Speed: 5}),
	// 	core.NewMoveAnimationFromParam(core.MoveParam{Tx: targetScreenPosX - (float64(assets.TileWidth) / 2), Ty: topLimit, Speed: 5}),
	// 	core.NewMoveAnimationFromParam(core.MoveParam{Tx: targetScreenPosX - (float64(assets.TileWidth) / 2), Ty: targetScreenPosY - (float64(assets.TileHeight + 25)), Speed: 5}),
	// )
	// fxEntry := ecs_.World.Create(component.Fx)
	// hh := ecs_.World.Entry(fxEntry)
	// component.Fx.Set(hh, &component.FxData{anim1})
	// anim1.Done = func() {
	// 	ecs_.World.Remove(fxEntry)
	// 	ecs_.World.Remove(targetGrid)
	// 	sprite.Image = assets.Bloombomber
	// 	dmgTile := ecs_.World.Create(
	// 		component.GridPos,
	// 		component.Damage,
	// 		component.OnHit,
	// 	)
	// 	ent := ecs_.World.Entry(dmgTile)
	// 	component.GridPos.Set(ent, &component.GridPosComponentData{Col: targetGridPos.Col, Row: targetGridPos.Row})
	// 	component.Damage.Set(ent, &component.DamageData{Damage: Damage})
	// 	component.OnHit.SetValue(ent, func(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	// 		damage := component.Damage.Get(projectile).Damage
	// 		component.Health.Get(receiver).HP -= damage
	// 		ecs.World.Remove(dmgTile)
	// 	})
	// 	query2 := donburi.NewQuery(filter.Contains(component.Health, component.GridPos))
	// 	query2.Each(ecs_.World, func(e *donburi.Entry) {
	// 		loc := component.GridPos.Get(e)
	// 		if loc.Col == targetGridPos.Col && loc.Row == targetGridPos.Row {
	// 			component.OnHit.GetValue(ent)(ecs_, ent, e)
	// 		}
	// 	})
	// 	ecs_.World.Remove(dmgTile)

	// }
	// sprite.Image = assets.BloombomberAtk
	// memory[ALREADY_FIRED] = true
}
