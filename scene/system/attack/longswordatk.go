package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

// func NewLongSwordAttack(ensetgetter ENSetGetter, ecs *ecs.ECS, playerScrLoc component.ScreenPosComponentData, playerGridLoc component.GridPosComponentData) {
// 	en := ensetgetter.GetEn()
// 	if en >= 200 {
// 		ensetgetter.SetEn(en - 200)
// 		newLongSwordAttack(ecs, playerScrLoc, playerGridLoc, 80)
// 	}
// }

type LongSwordCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
}

// cost 2 EN to execute. 80 dmg 2 tiles in front
func NewLongSwordCaster() *LongSwordCaster {
	return &LongSwordCaster{Cost: 200, Damage: 80, nextCooldown: time.Now()}
}
func (l *LongSwordCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nHit 2 grid in front for %d damage.\nNo cooldown", l.Cost/100, l.Damage)
}
func (l *LongSwordCaster) GetName() string {
	return "LongSwrd"
}
func (l *LongSwordCaster) GetDamage() int {
	return l.Damage
}
func (l *LongSwordCaster) Cast(ensource ENSetGetter, ecs *ecs.ECS) {
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
		playerScrLoc := component.ScreenPos.GetValue(playerEntry)
		playerGridLoc := component.GridPos.GetValue(playerEntry)
		newLongSwordAttack(ecs, playerScrLoc, playerGridLoc, l.Damage)
		l.nextCooldown = time.Now().Add(750 * time.Millisecond)
	}
}
func (l *LongSwordCaster) GetCost() int {
	return l.Cost
}
func (l *LongSwordCaster) GetIcon() *ebiten.Image {
	return assets.LongSwordIcon
}
func (l *LongSwordCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *LongSwordCaster) GetCooldownDuration() time.Duration {
	return 0
}
