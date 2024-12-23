package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type ShockWaveCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	Cooldown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func (l *ShockWaveCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *ShockWaveCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}
func NewShockwaveCaster() *ShockWaveCaster {
	return &ShockWaveCaster{Cost: 200, nextCooldown: time.Now(), Damage: 40, Cooldown: 6 * time.Second, OnHit: shockWaveOnAtkHit}
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
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *ShockWaveCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nHit 2 grid in front for %d damage.\nCooldown %.1f seconds", l.Cost/100, l.GetDamage(), l.GetCooldownDuration().Seconds())
}
func (l *ShockWaveCaster) GetName() string {
	return "ShockWave"
}

// cost 2 EN and inflict 40 DMG, slow moving projectile. Push back on enemy when hit
// cooldown for 2sec
func (c *ShockWaveCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= c.GetCost() {
		c.nextCooldown = time.Now().Add(c.GetCooldownDuration())
		ensource.SetEn(en - c.GetCost())
		query := donburi.NewQuery(
			filter.Contains(
				archetype.PlayerTag,
			),
		)
		AtkSfxQueue.QueueSFX(assets.HitscanFx)
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
		component.OnHit.SetValue(shockwaveEntry, c.OnHit)
		SPEED := 5.0
		component.Speed.Set(shockwaveEntry, &component.SpeedData{Vx: SPEED, Vy: 0})
		component.Damage.Set(shockwaveEntry, &component.DamageData{Damage: c.GetDamage()})
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
		if c.ModEntry != nil {
			// l := component.PostAtkModifier.GetValue(l.ModEntry)
			if c.ModEntry.PostAtk != nil {
				c.ModEntry.PostAtk(ecs, ensource)
			}
		}
	}

}
func (l *ShockWaveCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (c *ShockWaveCaster) GetCost() int {
	if c.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if c.Cost+c.ModEntry.CostModifier < 0 {
			return 0
		}
		return c.Cost + c.ModEntry.CostModifier
	}
	return c.Cost
}
func (c *ShockWaveCaster) GetIcon() *ebiten.Image {
	return assets.ShockwaveIcon
}
func (c *ShockWaveCaster) GetCooldown() time.Time {
	return c.nextCooldown
}
func (l *ShockWaveCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Cooldown + l.ModEntry.CooldownModifer
	}
	return l.Cooldown
}
