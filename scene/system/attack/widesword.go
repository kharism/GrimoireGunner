package attack

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type WideSwordCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
}

func NewWideSwordCaster() *WideSwordCaster {
	return &WideSwordCaster{Cost: 200, Damage: 100, nextCooldown: time.Now()}
}
func (l *WideSwordCaster) GetDamage() int {
	return l.Damage
}
func (l *WideSwordCaster) Cast(ensource ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.Cost {
		ensource.SetEn(en - l.Cost)
		query := donburi.NewQuery(
			filter.Contains(
				archetype.PlayerTag,
			),
		)

		playerEntry, ok := query.First(ecs.World)
		if !ok {
			return
		}
		playerGridLoc := component.GridPos.GetValue(playerEntry)
		var entry1 *donburi.Entry
		var entry2 *donburi.Entry
		var entry3 *donburi.Entry
		if playerGridLoc.Row > 0 {
			hitbox1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
			entry1 = ecs.World.Entry(hitbox1)
			component.Damage.Set(entry1, &component.DamageData{Damage: l.Damage})
			component.GridPos.Set(entry1, &component.GridPosComponentData{Row: playerGridLoc.Row - 1, Col: playerGridLoc.Col + 1})
			component.OnHit.SetValue(entry1, OnWideswordHit)
		}
		hitbox2 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
		entry2 = ecs.World.Entry(hitbox2)
		component.Damage.Set(entry2, &component.DamageData{Damage: l.Damage})
		component.GridPos.Set(entry2, &component.GridPosComponentData{Row: playerGridLoc.Row, Col: playerGridLoc.Col + 1})
		component.OnHit.SetValue(entry2, OnWideswordHit)
		if playerGridLoc.Row < 3 {
			hitbox3 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
			entry3 = ecs.World.Entry(hitbox3)
			component.Damage.Set(entry3, &component.DamageData{Damage: l.Damage})
			component.GridPos.Set(entry3, &component.GridPosComponentData{Row: playerGridLoc.Row + 1, Col: playerGridLoc.Col + 1})
			component.OnHit.SetValue(entry3, OnWideswordHit)
		}
		fxEntity := ecs.World.Create(component.Fx)
		fx := ecs.World.Entry(fxEntity)
		scrX, scrY := assets.GridCoord2Screen(playerGridLoc.Row-1, playerGridLoc.Col+1)
		scrX -= 50
		scrY -= 50
		wideSword := assets.NewWideSlashAtkAnim(assets.SpriteParam{
			ScreenX: scrX,
			ScreenY: scrY,
			Modulo:  5,
			Done: func() {
				if entry1 != nil {
					ecs.World.Remove(entry1.Entity())
				}
				if entry2 != nil {
					ecs.World.Remove(entry2.Entity())
				}
				if entry3 != nil {
					ecs.World.Remove(entry3.Entity())
				}
				ecs.World.Remove(fxEntity)
			},
		})
		wideSword.FlipHorizontal = true
		component.Fx.Set(fx, &component.FxData{Animation: wideSword})

	}
}

func OnWideswordHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	DmgComponent := component.Damage.Get(projectile)
	Health := component.Health.Get(receiver)
	Health.HP -= DmgComponent.Damage
	ecs.World.Remove(projectile.Entity())
}
func (l *WideSwordCaster) GetCost() int {
	return l.Cost
}
func (l *WideSwordCaster) GetIcon() *ebiten.Image {
	return assets.LongSwordIcon
}
func (l *WideSwordCaster) GetCooldown() time.Time {
	return l.nextCooldown
}