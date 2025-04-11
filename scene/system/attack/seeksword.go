package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type SeekswordCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func (l *SeekswordCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *SeekswordCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}
func NewSeekswordCaster() *SeekswordCaster {
	return &SeekswordCaster{Cost: 100, nextCooldown: time.Now(), Damage: 50, CoolDown: 3 * time.Second, OnHit: PushbackOnHit}
}

func (l *SeekswordCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Go to nearest enemies then use wide sword.\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.GetCooldownDuration().Seconds())
}
func (l *SeekswordCaster) GetName() string {
	return "SeekswordCaster"
}
func (l *SeekswordCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *SeekswordCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *SeekswordCaster) GetIcon() *ebiten.Image {
	return assets.SeekSwordIcon
}
func (l *SeekswordCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *SeekswordCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
func (l *SeekswordCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *SeekswordCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		closestTarget := GetNearestTarget(ecs)
		AtkSfxQueue.QueueSFX(assets.SlashFx)
		if closestTarget != nil {
			targetGridPos := component.GridPos.Get(closestTarget)
			oldGridPos, player := GetPlayerGridPos(ecs)
			player.AddComponent(component.Root)
			foundMove := false
			var newGridPos *component.GridPosComponentData

			for i := -1; i <= 1; i++ {
				if validMove(ecs, targetGridPos.Row+i, targetGridPos.Col-1) {
					newGridPos = &component.GridPosComponentData{Row: targetGridPos.Row + i, Col: targetGridPos.Col - 1}
					foundMove = true
				}
			}
			if foundMove {
				component.GridPos.Set(player, newGridPos)
				entries := []*donburi.Entry{}
				for i := -1; i <= 1; i++ {
					if newGridPos.Row+i == -1 || newGridPos.Row+i == 4 {
						continue
					}
					pp := ecs.World.Create(component.GridPos, component.OnHit, component.Damage, component.Transient)
					entry := ecs.World.Entry(pp)
					component.GridPos.Set(entry, &component.GridPosComponentData{Row: newGridPos.Row + i, Col: newGridPos.Col + 1})
					component.Damage.Set(entry, &component.DamageData{Damage: l.GetDamage()})
					component.OnHit.SetValue(entry, SingleHitProjectile)
					component.Transient.Set(entry, &component.TransientData{Start: time.Now(), Duration: 200 * time.Millisecond})
					entries = append(entries, entry)
				}
				fxEntity := ecs.World.Create(component.Fx)
				AtkSfxQueue.QueueSFX(assets.SlashFx)
				fx := ecs.World.Entry(fxEntity)
				scrX, scrY := assets.GridCoord2Screen(newGridPos.Row-1, newGridPos.Col+1)
				scrX -= 50
				scrY -= 50
				wideSword := assets.NewWideSlashAtkAnim(assets.SpriteParam{
					ScreenX: scrX,
					ScreenY: scrY,
					Modulo:  5,
					Done: func() {
						ecs.World.Remove(fxEntity)
						player.RemoveComponent(component.Root)
						component.GridPos.Set(player, oldGridPos)
					},
				})
				wideSword.FlipHorizontal = true
				component.Fx.Set(fx, &component.FxData{Animation: wideSword})
			}
		}
		if l.ModEntry != nil {
			// l := component.PostAtkModifier.GetValue(l.ModEntry)
			if l.ModEntry.PostAtk != nil {
				l.ModEntry.PostAtk(ecs, ensource)
			}
		}
	}
}
