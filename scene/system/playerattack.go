package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type PlayerAttackSystem struct {
	PlayerIndex *donburi.Entity
}

func NewPlayerAttackSystem(player *donburi.Entity) *PlayerAttackSystem {
	return &PlayerAttackSystem{PlayerIndex: player}
}
func GenerateMagibullet(ecs *ecs.ECS, row, col int, xspeed float64) {

	archetype.NewProjectile(ecs.World, archetype.ProjectileParam{
		Vx:     xspeed,
		Vy:     0,
		Col:    col,
		Row:    row,
		Damage: 40,
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
func (s *PlayerAttackSystem) Update(ecs *ecs.ECS) {
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		playerId := ecs.World.Entry(*s.PlayerIndex)
		gridPos := component.GridPos.Get(playerId)
		GenerateMagibullet(ecs, gridPos.Row, gridPos.Col+1, 15)
	}
}
