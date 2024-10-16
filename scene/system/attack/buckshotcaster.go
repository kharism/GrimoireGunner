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
	"github.com/yohamta/donburi/filter"
)

// cost 2 EN and cast bullet in T-shaped cone
type BuckshotCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func (l *BuckshotCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *BuckshotCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}
func NewBuckshotCaster() *BuckshotCaster {
	return &BuckshotCaster{Cost: 200, nextCooldown: time.Now(), Damage: 150, CoolDown: 2 * time.Second, OnHit: SingleHitProjectile}
}
func (l *BuckshotCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Damage in T-shaped cone in front. Hit the target in front 3 times\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.GetCooldownDuration().Seconds())
}
func (l *BuckshotCaster) GetName() string {
	return "BuckShot"
}
func (l *BuckshotCaster) GetDamage() int {
	return l.Damage
}
func (l *BuckshotCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.Cost {
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
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
		grid1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
		grid1Entry := ecs.World.Entry(grid1)
		grid1b := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
		grid1bEntry := ecs.World.Entry(grid1b)
		grid1c := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
		grid1cEntry := ecs.World.Entry(grid1c)
		grid2 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
		grid2Entry := ecs.World.Entry(grid2)
		component.Damage.Set(grid1Entry, &component.DamageData{Damage: l.Damage})
		component.Damage.Set(grid2Entry, &component.DamageData{Damage: l.Damage})
		component.Damage.Set(grid1bEntry, &component.DamageData{Damage: l.Damage})
		component.Damage.Set(grid1cEntry, &component.DamageData{Damage: l.Damage})
		component.OnHit.SetValue(grid1Entry, l.OnHit)
		component.OnHit.SetValue(grid2Entry, l.OnHit)
		component.OnHit.SetValue(grid1bEntry, l.OnHit)
		component.OnHit.SetValue(grid1cEntry, l.OnHit)
		component.GridPos.Set(grid1Entry, &component.GridPosComponentData{Col: gridPos.Col + 1, Row: gridPos.Row})
		component.GridPos.Set(grid1bEntry, &component.GridPosComponentData{Col: gridPos.Col + 1, Row: gridPos.Row})
		component.GridPos.Set(grid1cEntry, &component.GridPosComponentData{Col: gridPos.Col + 1, Row: gridPos.Row})
		component.GridPos.Set(grid2Entry, &component.GridPosComponentData{Col: gridPos.Col + 2, Row: gridPos.Row})
		var grid3Entry *donburi.Entry
		var grid4Entry *donburi.Entry
		if gridPos.Row > 0 {
			grid3 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
			grid3Entry = ecs.World.Entry(grid3)
			component.Damage.Set(grid3Entry, &component.DamageData{Damage: l.Damage})
			component.OnHit.SetValue(grid3Entry, l.OnHit)
			component.GridPos.Set(grid3Entry, &component.GridPosComponentData{Col: gridPos.Col + 2, Row: gridPos.Row - 1})
		}
		if gridPos.Row < 3 {
			grid4 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
			grid4Entry = ecs.World.Entry(grid4)
			component.Damage.Set(grid4Entry, &component.DamageData{Damage: l.Damage})
			component.OnHit.SetValue(grid4Entry, l.OnHit)
			component.GridPos.Set(grid4Entry, &component.GridPosComponentData{Col: gridPos.Col + 2, Row: gridPos.Row + 1})
		}
		fx := ecs.World.Create(component.Fx)
		fxEnt := ecs.World.Entry(fx)
		scrX, scrY := assets.GridCoord2Screen(gridPos.Row-1, gridPos.Col+1)
		fmt.Println(gridPos.Row-1, gridPos.Col+1, scrX, scrY)
		buckShotAnim := assets.NewBuckshotAtkAnim(assets.SpriteParam{
			ScreenX: scrX - 50,
			ScreenY: scrY - 50,
			Modulo:  3,
			Done: func() {
				ecs.World.Remove(grid1)
				ecs.World.Remove(grid1b)
				ecs.World.Remove(grid1c)
				ecs.World.Remove(grid2)
				if grid3Entry != nil {
					ecs.World.Remove(grid3Entry.Entity())
				}
				if grid4Entry != nil {
					ecs.World.Remove(grid4Entry.Entity())
				}
				ecs.World.Remove(fx)
			},
		})
		if l.ModEntry != nil {
			// l := component.PostAtkModifier.GetValue(l.ModEntry)
			if l.ModEntry.PostAtk != nil {
				l.ModEntry.PostAtk(ecs, ensource)
			}
		}
		component.Fx.Set(fxEnt, &component.FxData{Animation: buckShotAnim})
	}

}
func (l *BuckshotCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *BuckshotCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *BuckshotCaster) GetIcon() *ebiten.Image {
	return assets.BuckshotIcon
}
func (l *BuckshotCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *BuckshotCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
