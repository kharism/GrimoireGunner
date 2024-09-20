package attack

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi/ecs"
)

type AtkBonus struct {
	nextCooldown time.Time
}

func (a *AtkBonus) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {

}
func (a *AtkBonus) GetCost() int {
	return 300
}
func (a *AtkBonus) GetIcon() *ebiten.Image {
	return nil
}
func (a *AtkBonus) GetCooldown() time.Time {
	return a.nextCooldown
}
func (a *AtkBonus) GetCooldownDuration() time.Duration {
	return 20 * time.Second
}
func (a *AtkBonus) GetDamage() int {
	return 0
}

func (a *AtkBonus) GetDescription() string {
	return "Give +10 damage to caster on other slot"
}
func (a *AtkBonus) GetName() string {
	return "AtkBoost"
}
