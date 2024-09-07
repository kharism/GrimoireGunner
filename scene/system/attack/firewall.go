package attack

import (
	"fmt"
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

func NewFirewallAttack(ecs *ecs.ECS, sourceRow, sourceCol, damage int) {
	sourceScrX, sourceSrcY := assets.GridCoord2Screen(sourceRow, sourceCol)
	fmt.Println(sourceRow, sourceCol, sourceScrX, sourceSrcY)
	sourceScrX -= 50
	sourceSrcY -= 50
	for i := 0; i < 4; i++ {
		targetScrX, targetScrY := assets.GridCoord2Screen(i, sourceCol+4)
		targetScrX -= 50
		targetScrY -= 100
		row := i
		movableImg := core.NewMovableImage(assets.Projectile1, core.NewMovableImageParams().WithMoveParam(
			core.MoveParam{Sx: sourceScrX, Sy: sourceSrcY, Tx: targetScrX, Ty: targetScrY, Speed: 10}))
		fx := ecs.World.Create(component.Fx)
		entry := ecs.World.Entry(fx)
		movableImg.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{Tx: targetScrX, Ty: targetScrY, Speed: 10}))
		movableImg.Done = func() {
			ecs.World.Remove(fx)
			entity := ecs.World.Create(component.Damage, component.GridPos, component.Transient, component.OnHit, component.Fx)
			entry := ecs.World.Entry(entity)
			component.Damage.Set(entry, &component.DamageData{Damage: damage})
			component.GridPos.Set(entry, &component.GridPosComponentData{Col: sourceCol + 4, Row: row})
			component.Transient.Set(entry, &component.TransientData{Start: time.Now(), Duration: 5 * time.Second})
			component.OnHit.SetValue(entry, OnTowerHit)
			flameTower := core.NewMovableImage(assets.FlametowerRaw, core.NewMovableImageParams().
				WithMoveParam(core.MoveParam{Sx: targetScrX, Sy: targetScrY, Speed: 3}))
			component.Fx.Set(entry, &component.FxData{Animation: flameTower})
		}
		component.Fx.Set(entry, &component.FxData{Animation: movableImg})
	}
}
func OnTowerHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	health := component.Health.Get(receiver)
	damage := component.Damage.Get(projectile)
	health.HP -= damage.Damage
	health.InvisTime = time.Now().Add(1 * time.Second)
}

type FirewallCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	Cooldown     time.Duration
}

func NewFirewallCaster() *FirewallCaster {
	return &FirewallCaster{Cost: 200, nextCooldown: time.Now(), Cooldown: 2 * time.Second}
}
func (f *FirewallCaster) GetDamage() int {
	return f.Damage
}
func (f *FirewallCaster) Cast(ensource ENSetGetter, ecs *ecs.ECS) {
	curEn := ensource.GetEn()
	if curEn >= f.Cost {
		ensource.SetEn(curEn - f.Cost)
		f.nextCooldown = time.Now().Add(f.Cooldown)
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
		NewFirewallAttack(ecs, gridPos.Row, gridPos.Col, f.Damage)
	}
}
func (f *FirewallCaster) GetCost() int {
	return f.Cost
}
func (f *FirewallCaster) GetIcon() *ebiten.Image {
	return assets.FirewallIcon
}
func (f *FirewallCaster) GetCooldown() time.Time {
	return f.nextCooldown
}
