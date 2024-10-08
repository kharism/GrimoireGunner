package attack

import (
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func GenerateMagibullet(ecs *ecs.ECS, row, col int, xspeed float64) *donburi.Entity {

	return archetype.NewProjectile(ecs.World, archetype.ProjectileParam{
		Vx:     xspeed,
		Vy:     0,
		Col:    col,
		Row:    row,
		Damage: 2,
		Sprite: assets.Projectile1,
		OnHit:  SingleHitProjectile,
	})
}

// use this as single hit projectile. Once a projectile hit,
// apply damage then disappear. Basically the default behaviour of any projectile based attack
func SingleHitProjectile(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	damage := component.Damage.Get(projectile).Damage
	component.Health.Get(receiver).HP -= damage
	ecs.World.Remove(projectile.Entity())
}
