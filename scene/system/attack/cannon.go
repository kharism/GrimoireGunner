package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi/ecs"
)

// cost 1 EN and cast hitscan bullet
type CannonCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
}

func NewCannonCaster() *CannonCaster {
	return &CannonCaster{Cost: 100, nextCooldown: time.Now(), Damage: 80, CoolDown: 1 * time.Second}
}
func (l *CannonCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Damage 1 target on front immediately.\nCooldown %.1fs", l.Cost/100, l.Damage, l.CoolDown.Seconds())
}
func (l *CannonCaster) GetName() string {
	return "Cannon"
}
func (l *CannonCaster) GetDamage() int {
	return l.Damage
}
func (l *CannonCaster) Cast(ensource ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.Cost {
		l.nextCooldown = time.Now().Add(l.CoolDown)
		closestTarget := HitScanGetNearestTarget(ecs)
		if closestTarget != nil {
			grid1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit)
			grid1Entry := ecs.World.Entry(grid1)
			targetGridPos := component.GridPos.Get(closestTarget)
			component.GridPos.Set(grid1Entry, &component.GridPosComponentData{Col: targetGridPos.Col, Row: targetGridPos.Row})
			component.Damage.Set(grid1Entry, &component.DamageData{Damage: l.Damage})
			component.OnHit.SetValue(grid1Entry, SingleHitProjectile)
		}
	}
}

func (l *CannonCaster) GetCost() int {
	return l.Cost
}
func (l *CannonCaster) GetIcon() *ebiten.Image {
	return assets.CannonIcon
}
func (l *CannonCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *CannonCaster) GetCooldownDuration() time.Duration {
	return l.CoolDown
}
