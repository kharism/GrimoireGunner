package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/hanashi/core"
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
func NewLongSwordAttack(ecs *ecs.ECS, playerScrLoc component.ScreenPosComponentData, playerGridLoc component.GridPosComponentData) {
	param := DamageGridParam{}
	loc1 := &component.GridPosComponentData{Row: playerGridLoc.Row, Col: playerGridLoc.Col + 1}
	loc2 := &component.GridPosComponentData{Row: playerGridLoc.Row, Col: playerGridLoc.Col + 2}
	param.Locations = []*component.GridPosComponentData{loc1, loc2}
	param.Damage = []int{20, 20}
	param.OnHit = SingleHitProjectile

	param.Fx = assets.NewSwordAtkAnim(assets.SpriteParam{
		ScreenX: playerScrLoc.X + float64(tileWidth)/2,
		ScreenY: playerScrLoc.Y - float64(tileHeight),
		Modulo:  2,
	})
	NonProjectileAtk(ecs, param)
}

type DamageGridParam struct {
	Locations []*component.GridPosComponentData
	Damage    []int
	Fx        *core.AnimatedImage
	OnHit     component.OnAtkHit
}

func NonProjectileAtk(ecs *ecs.ECS, param DamageGridParam) {
	removeList := []donburi.Entity{}
	for idx, l := range param.Locations {
		entity := ecs.World.Create(
			component.Damage,
			component.GridPos,
			component.OnHit,
		)
		entry := ecs.World.Entry(entity)
		removeList = append(removeList, entity)
		component.Damage.Set(entry, &component.DamageData{param.Damage[idx]})
		component.GridPos.Set(entry, l)
		component.OnHit.Set(entry, &param.OnHit)
	}

	entity := ecs.World.Create(component.Fx)
	entry := ecs.World.Entry(entity)
	param.Fx.Done = func() {
		ecs.World.Remove(removeList[0])
		ecs.World.Remove(removeList[1])
		ecs.World.Remove(entity)
	}
	component.Fx.Set(entry, &param.Fx)
}
func (s *PlayerAttackSystem) Update(ecs *ecs.ECS) {
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		playerId := ecs.World.Entry(*s.PlayerIndex)
		gridPos := component.GridPos.Get(playerId)
		GenerateMagibullet(ecs, gridPos.Row, gridPos.Col+1, 15)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		playerId := ecs.World.Entry(*s.PlayerIndex)
		gridPos := component.GridPos.Get(playerId)
		scrPos := component.ScreenPos.Get(playerId)
		NewLongSwordAttack(ecs, *scrPos, *gridPos)
	}
}
