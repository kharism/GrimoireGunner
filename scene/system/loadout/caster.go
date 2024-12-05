package loadout

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi/ecs"
)

// cast our spell.
// any combat or non combat will utilize this
type Caster interface {
	Cast(ensource ENSetGetter, ecs *ecs.ECS)
	GetCost() int
	GetIcon() *ebiten.Image
	GetCooldown() time.Time
	GetCooldownDuration() time.Duration
	GetDamage() int

	ResetCooldown()

	GetDescription() string
	GetName() string
}
type ElementalCaster interface {
	Caster
	GetElement() component.Elemental
}
type ENSetGetter interface {
	SetEn(val int)
	GetEn() int
	GetMaxEn() int
}
