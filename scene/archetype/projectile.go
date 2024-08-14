package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	mycomponent "github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
)

var ProjectileTag = donburi.NewTag("Projectile")

type ProjectileParam struct {
	Vx, Vy   float64
	Col, Row int
	Damage   int
	Sprite   *ebiten.Image
	OnHit    mycomponent.OnAtkHit
}

// create projectile that moves across the field
func NewProjectile(world donburi.World, param ProjectileParam) *donburi.Entity {
	entity := world.Create(
		mycomponent.GridPos,
		mycomponent.ScreenPos,
		mycomponent.Speed,
		mycomponent.Damage,
		mycomponent.OnHit,
		mycomponent.Sprite, ProjectileTag)
	entId := world.Entry(entity)
	mycomponent.Damage.Set(entId, &mycomponent.DamageData{Damage: param.Damage})
	mycomponent.Speed.Set(entId, &mycomponent.SpeedData{Vx: param.Vx, Vy: param.Vy})
	mycomponent.GridPos.Set(entId, &mycomponent.GridPosComponentData{Row: param.Row, Col: param.Col})
	mycomponent.Sprite.Set(entId, &mycomponent.SpriteData{Image: param.Sprite})
	mycomponent.OnHit.Set(entId, &param.OnHit)
	return &entity
}
