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

type GatlingCaster struct {
	Cost             int
	Damage           int
	ShotAmount       int
	nextCooldown     time.Time
	CooldownDuration time.Duration
}

func NewGatlingCastor() *GatlingCaster {
	return &GatlingCaster{Cost: 200, Damage: 5, ShotAmount: 10, nextCooldown: time.Now(), CooldownDuration: 3 * time.Second}
}
func (l *GatlingCaster) GetDamage() int {
	return l.Damage
}
func (l *GatlingCaster) Cast(ensource ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.Cost {
		ensource.SetEn(en - l.Cost)
		l.nextCooldown = time.Now().Add(l.CooldownDuration)
		now := time.Now()
		for i := 0; i < l.ShotAmount; i++ {
			ev := &addGatlingShot{Damage: l.Damage, Time: now.Add(time.Duration(100*i) * time.Millisecond)}
			component.EventQueue.Queue = append(component.EventQueue.Queue, ev)
		}
	}

}

type addGatlingShot struct {
	Damage int
	Time   time.Time
}

func (a *addGatlingShot) Execute(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.Contains(
			archetype.PlayerTag,
		),
	)

	playerEntry, ok := query.First(ecs.World)
	if !ok {
		return
	}
	// playerScrLoc := component.ScreenPos.GetValue(playerEntry)
	playerGridLoc := component.GridPos.GetValue(playerEntry)
	archetype.NewProjectile(ecs.World, archetype.ProjectileParam{
		Vx:     15,
		Vy:     0,
		Col:    playerGridLoc.Col + 1,
		Row:    playerGridLoc.Row,
		Sprite: assets.Projectile1,
		Damage: a.Damage,
		OnHit:  SingleHitProjectile,
	})
}
func (a *addGatlingShot) GetTime() time.Time {
	return a.Time
}
func (l *GatlingCaster) GetCost() int {
	return l.Cost
}
func (l *GatlingCaster) GetIcon() *ebiten.Image {
	return assets.GatlingIcon
}
func (l *GatlingCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
