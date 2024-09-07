package attack

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type ShockWaveCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
}

func NewShockwaveCaster() *ShockWaveCaster {
	return &ShockWaveCaster{Cost: 200, nextCooldown: time.Now(), Damage: 40}
}

var queryHP = donburi.NewQuery(
	filter.Contains(
		component.Health,
		component.GridPos,
	),
)

// check whether there are obstacle on row-col grid
// for now it checks another character
func validMove(ecs *ecs.ECS, row, col int) bool {
	if col >= 8 || col < 0 || row < 0 || row >= 4 {
		return false
	}
	ObstacleExist := false
	queryHP.Each(ecs.World, func(e *donburi.Entry) {
		pos := component.GridPos.Get(e)
		if pos.Col == col && pos.Row == row {
			ObstacleExist = true
		}
	})
	return !ObstacleExist
}

func shockWaveOnAtkHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	damage := component.Damage.Get(projectile).Damage
	healthComp := component.Health.Get(receiver)
	healthComp.HP -= damage
	if !receiver.HasComponent(archetype.ConstructTag) {
		healthComp.InvisTime = time.Now().Add(400 * time.Millisecond)
	} else {
		projectileGridPos := component.GridPos.Get(projectile)
		projectileGridPos.Row += 1
		scrPosProjectile := component.ScreenPos.Get(projectile)
		scrPosProjectile.X += float64(assets.TileWidth)
	}

	receiverPos := component.GridPos.Get(receiver)
	if validMove(ecs, receiverPos.Row, receiverPos.Col+1) {
		receiverPos.Col += 1
		scrPos := component.ScreenPos.Get(receiver)
		scrPos.X += 100
	}
}
func (l *ShockWaveCaster) GetDamage() int {
	return l.Damage
}

// cost 2 EN and inflict 40 DMG, slow moving projectile. Push back on enemy when hit
// cooldown for 2sec
func (c *ShockWaveCaster) Cast(ensource ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= 200 {
		c.nextCooldown = time.Now().Add(2 * time.Second)
		ensource.SetEn(en - c.Cost)
		query := donburi.NewQuery(
			filter.Contains(
				archetype.PlayerTag,
			),
		)

		playerId, ok := query.First(ecs.World)
		if !ok {
			return
		}
		gridPos := component.GridPos.Get(playerId)
		shockwave := ecs.World.Create(
			component.GridPos,
			component.ScreenPos,
			component.Speed,
			component.Damage,
			// component.Sprite,
			component.OnHit,
			component.Fx,
			component.TargetLocation,
		)
		shockwaveEntry := ecs.World.Entry(shockwave)
		component.GridPos.Set(shockwaveEntry, &component.GridPosComponentData{
			Row: gridPos.Row,
			Col: gridPos.Col + 1,
		})
		screenTargetX, screenTargetY := assets.GridCoord2Screen(gridPos.Row, 8)
		component.TargetLocation.Set(shockwaveEntry, &component.MoveTargetData{
			Tx: screenTargetX,
			Ty: screenTargetY,
		})
		component.OnHit.SetValue(shockwaveEntry, shockWaveOnAtkHit)
		SPEED := 5.0
		component.Speed.Set(shockwaveEntry, &component.SpeedData{Vx: SPEED, Vy: 0})
		component.Damage.Set(shockwaveEntry, &component.DamageData{Damage: c.Damage})
		screenX, screenY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col+1)
		screenX = screenX - 50
		screenY = screenY - 100
		shockwaveAnim := assets.NewShockwaveAnim(assets.SpriteParam{
			ScreenX: screenX,
			ScreenY: screenY,
			Modulo:  10,
			Done: func() {
				// ecs.World.Remove(shockwave)
			},
		})
		shockwaveAnim.AddAnimation(
			core.NewMoveAnimationFromParam(
				core.MoveParam{
					Tx:    screenTargetX,
					Ty:    screenY,
					Speed: 5,
				}),
		)
		// shockwaveAnim.MovableImage
		component.Fx.Set(shockwaveEntry, &component.FxData{Animation: shockwaveAnim})
	}

}
func (c *ShockWaveCaster) GetCost() int {
	return c.Cost
}
func (c *ShockWaveCaster) GetIcon() *ebiten.Image {
	return assets.ShockwaveIcon
}
func (c *ShockWaveCaster) GetCooldown() time.Time {
	return c.nextCooldown
}
