package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi/ecs"
)

type WallCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func NewWallCaster() *WallCaster {
	return &WallCaster{
		Cost:         100,
		Damage:       0,
		nextCooldown: time.Now(),
		CoolDown:     5 * time.Second,
	}
}
func (l *WallCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *WallCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}
func (l *WallCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nCreate a wall with 150HP.\nCooldown %.1fs", l.Cost/100, l.GetCooldownDuration().Seconds())
}
func (l *WallCaster) GetName() string {
	return "WallCaster"
}
func (l *WallCaster) GetDamage() int {
	return 0
}
func (l *WallCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *WallCaster) GetIcon() *ebiten.Image {
	return assets.WallIcon
}
func (l *WallCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *WallCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
func (l *WallCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *WallCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		gridPos, _ := GetPlayerGridPos(ecs)
		if validMove(ecs, gridPos.Row, gridPos.Col+1) {
			wallEntity := archetype.NewConstruct(ecs.World, assets.Wall)
			wallEntry := ecs.World.Entry(*wallEntity)
			component.GridPos.Set(wallEntry, &component.GridPosComponentData{Col: gridPos.Col + 1, Row: gridPos.Row})
			component.Health.Set(wallEntry, &component.HealthData{HP: 150})
		}
	}
}
