package system

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/yohamta/donburi/ecs"
)

// cast our spell.
// any combat or non combat will utilize this
type Caster interface {
	Cast(ensource attack.ENSetGetter, ecs *ecs.ECS)
	GetCost() int
	GetIcon() *ebiten.Image
	GetCooldown() time.Time
	GetDamage() int
}
