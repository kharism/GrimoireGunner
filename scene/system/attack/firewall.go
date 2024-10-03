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
	ModEntry     *loadout.CasterModifierData
}

func NewFirewallCaster() *FirewallCaster {
	return &FirewallCaster{Cost: 200, nextCooldown: time.Now(), Cooldown: 2 * time.Second, Damage: 10}
}
func (l *FirewallCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *FirewallCaster) SetModifier(e *loadout.CasterModifierData) {
	l.ModEntry = e
}
func (f *FirewallCaster) GetDamage() int {
	if f.ModEntry != nil {
		// mod := component.CasterModifier.Get(f.ModEntry)
		return f.Damage + f.ModEntry.DamageModifier
	}
	return f.Damage
}
func (l *FirewallCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nCreate firewall which damage %d if stepped on.\nCooldown %.1fs", l.Cost/100, l.Damage, l.Cooldown.Seconds())
}
func (l *FirewallCaster) GetName() string {
	return "Firewall"
}
func (f *FirewallCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
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
		if f.ModEntry != nil {
			// l := component.PostAtkModifier.GetValue(f.ModEntry)
			if f.ModEntry.PostAtk != nil {
				f.ModEntry.PostAtk(ecs, ensource)
			}
		}
	}
}
func (f *FirewallCaster) GetCost() int {
	if f.ModEntry != nil {
		// mod := component.CasterModifier.Get(f.ModEntry)
		if f.Cost+f.ModEntry.CostModifier < 0 {
			return 0
		}
		return f.Cost + f.ModEntry.CostModifier
	}
	return f.Cost
}
func (l *FirewallCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (f *FirewallCaster) GetIcon() *ebiten.Image {
	return assets.FirewallIcon
}
func (f *FirewallCaster) GetCooldown() time.Time {
	return f.nextCooldown
}
func (f *FirewallCaster) GetCooldownDuration() time.Duration {
	if f.ModEntry != nil {
		// mod := component.CasterModifier.Get(f.ModEntry)
		return f.Cooldown + f.ModEntry.CooldownModifer
	}
	return f.Cooldown
}
