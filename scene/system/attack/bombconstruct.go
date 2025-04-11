package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type BombConstructCaster struct {
	Cost             int
	Damage           int
	ShotAmount       int
	nextCooldown     time.Time
	CooldownDuration time.Duration
	ModEntry         *loadout.CasterModifierData
	OnHit            component.OnAtkHit
}

func (l *BombConstructCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func NewBombConstructCaster() *BombConstructCaster {
	return &BombConstructCaster{Cost: 200, Damage: 200, ShotAmount: 10, nextCooldown: time.Now(), CooldownDuration: 3 * time.Second, OnHit: SingleHitProjectile}
}
func (l *BombConstructCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nCreate bomb (200HP) which explode for %d damage when destroyed.\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.CooldownDuration.Seconds())
}
func (l *BombConstructCaster) GetName() string {
	return "Bomb"
}
func (l *BombConstructCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		playerGridPos, _ := GetPlayerGridPos(ecs)

		if !validMove(ecs, playerGridPos.Row, playerGridPos.Col+4) {
			// entry := getEntryAtGridPos(ecs, playerGridPos.Row, playerGridPos.Col+4)
			dmgEntity := ecs.World.Create(component.GridPos, component.Damage, component.OnHit)
			dmgEntry := ecs.World.Entry(dmgEntity)
			component.GridPos.Set(dmgEntry, &component.GridPosComponentData{
				Col: playerGridPos.Col + 4,
				Row: playerGridPos.Row,
			})
			component.Damage.Set(dmgEntry, &component.DamageData{
				Damage: 200,
			})
			component.OnHit.SetValue(dmgEntry, SingleHitProjectile)
			return
		}
		bombEntity := archetype.NewConstruct(ecs.World, assets.Bomb2)
		bombEntry := ecs.World.Entry(*bombEntity)
		bombEntry.AddComponent(component.OnDestroy)
		bombEntry.AddComponent(component.OnHit)

		bombGridPos := component.GridPos.Get(bombEntry)
		bombGridPos.Col = playerGridPos.Col + 4
		bombGridPos.Row = playerGridPos.Row
		component.OnDestroy.SetValue(bombEntry, OnBombDestroyed)
		component.OnHit.SetValue(bombEntry, SingleHitProjectile)
	}

}
func OnBombDestroyed(ecs *ecs.ECS, entry *donburi.Entry) {
	gridPOs := component.GridPos.Get(entry)
	onHit := component.OnHit.Get(entry)
	AtkSfxQueue.QueueSFX(assets.ExplosionFx)
	for col := gridPOs.Col - 1; col <= gridPOs.Col+1; col++ {
		for row := gridPOs.Row - 1; row <= gridPOs.Row+1; row++ {
			if col >= 8 || col < 0 || row < 0 || row >= 4 {
				continue
			}
			damageTile := ecs.World.Create(component.GridPos, component.OnHit, component.ScreenPos, component.Damage, component.Fx)
			dmgTile := ecs.World.Entry(damageTile)
			// component.Transient.Set(dmgTile, &component.TransientData{
			// 	Start:    time.Now(),
			// 	Duration: 300 * time.Millisecond,
			// })
			component.Damage.Set(dmgTile, &component.DamageData{
				Damage: 200,
			})
			component.GridPos.Set(dmgTile, &component.GridPosComponentData{
				Col: col,
				Row: row,
			})

			scrPos := component.ScreenPos.Get(dmgTile)
			scrPos.X, scrPos.Y = assets.GridCoord2Screen(row, col)
			// gridMap[gridPos.Row][gridPos.Col] = nil
			// ecs.World.Remove(damageableEntity.Entity())

			explosionAnim := assets.NewExplosionAnim(assets.SpriteParam{
				ScreenX: scrPos.X - float64(assets.TileWidth)/2,
				ScreenY: scrPos.Y - 75,
				Modulo:  5,
			})
			explosionAnim.Done = func() {
				ecs.World.Remove(damageTile)
			}
			component.OnHit.Set(dmgTile, onHit)
			component.Fx.Set(dmgTile, &component.FxData{explosionAnim})
		}
	}
}
func (l *BombConstructCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *BombConstructCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *BombConstructCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}
func (l *BombConstructCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *BombConstructCaster) GetIcon() *ebiten.Image {
	return assets.BombIcon
}
func (l *BombConstructCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *BombConstructCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CooldownDuration + l.ModEntry.CooldownModifer
	}
	return l.CooldownDuration
}
