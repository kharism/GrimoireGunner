package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// delayed attack for 500ms
type ChargeshotCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func NewChargeshotCaster() *ChargeshotCaster {
	return &ChargeshotCaster{Cost: 300, nextCooldown: time.Now(), Damage: 200, CoolDown: 9 * time.Second, OnHit: SingleHitProjectile}
}
func (l *ChargeshotCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Charge for 500ms then shot projectile.\nCooldown %.1fs", l.Cost/100, l.Damage, l.CoolDown.Seconds())
}
func (l *ChargeshotCaster) GetName() string {
	return "Chargeshot"
}
func (l *ChargeshotCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		gridPos, playerentry := GetPlayerGridPos(ecs)
		chargeFxentity := ecs.World.Create(component.Fx, component.ScreenPos, component.Transient)
		chargeFxEntry := ecs.World.Entry(chargeFxentity)
		playerentry.AddComponent(component.Root)
		playerScrPosX, playerScrPosY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
		playerScrPosX = playerScrPosX - 50
		playerScrPosY = playerScrPosY - 100
		component.ScreenPos.Set(chargeFxEntry, &component.ScreenPosComponentData{
			X: playerScrPosX,
			Y: playerScrPosY,
		})
		component.Fx.Set(chargeFxEntry, &component.FxData{
			Animation: core.NewMovableImage(
				assets.ChargeshotFx,
				core.NewMovableImageParams().
					WithMoveParam(core.MoveParam{Sx: playerScrPosX, Sy: playerScrPosY})),
		})
		component.Transient.Set(chargeFxEntry, &component.TransientData{Start: time.Now(), Duration: 1000 * time.Millisecond, OnRemoveCallback: l.CreateChargedProjectile(ecs, *gridPos)})
	}

}
func (l *ChargeshotCaster) CreateChargedProjectile(ecs2 *ecs.ECS, gridPos component.GridPosComponentData) func(ecs *ecs.ECS, entity *donburi.Entry) {
	return func(ecs *ecs.ECS, entity *donburi.Entry) {
		chargeShotEntity := ecs.World.Create(component.Fx, component.GridPos, component.Damage, component.TargetLocation, component.OnHit, component.Speed, component.ScreenPos, archetype.ProjectileTag)
		chargeShotEntry := ecs.World.Entry(chargeShotEntity)
		_, playerentry := GetPlayerGridPos(ecs)
		playerentry.RemoveComponent(component.Root)
		component.Damage.Set(chargeShotEntry, &component.DamageData{Damage: l.Damage})
		component.Speed.Set(chargeShotEntry, &component.SpeedData{Vx: 10, Vy: 0})
		component.GridPos.Set(chargeShotEntry, &component.GridPosComponentData{Row: gridPos.Row, Col: gridPos.Col + 1})
		// component.Sprite.Set(chargeShotEntry, &component.SpriteData{Image: param.Sprite})
		screenX, screenY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col+1)
		screenX = screenX - 50
		screenY = screenY - 150
		animation := assets.NewChargeShotAnim(assets.SpriteParam{
			ScreenX: screenX,
			ScreenY: screenY,
			Modulo:  5,
			Done: func() {

			},
		})
		screenTargetX, screenTargetY := assets.GridCoord2Screen(gridPos.Row, 8)
		component.TargetLocation.Set(chargeShotEntry, &component.MoveTargetData{
			Tx: screenTargetX,
			Ty: screenTargetY,
		})
		animation.AddAnimation(
			core.NewMoveAnimationFromParam(
				core.MoveParam{
					Tx:    screenTargetX,
					Ty:    screenY,
					Speed: 10,
				}),
		)
		component.Fx.Set(chargeShotEntry, &component.FxData{Animation: animation})
		component.OnHit.SetValue(chargeShotEntry, l.OnHit)
	}
}
func (l *ChargeshotCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *ChargeshotCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *ChargeshotCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}

func (l *ChargeshotCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *ChargeshotCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *ChargeshotCaster) GetIcon() *ebiten.Image {
	return assets.ChargeshotIcon
}
func (l *ChargeshotCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *ChargeshotCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
