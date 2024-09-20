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
	"github.com/yohamta/donburi/filter"
)

type LightnigAtkParam struct {
	StartRow, StartCol int
	//+1 to create lightning on right of starting point, and -1 to create on left
	Direction int
	Damage    int
	Actor     *donburi.Entry
}

func LightningBoltOnHitfunc(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	damage := component.Damage.Get(projectile).Damage
	health := component.Health.Get(receiver)
	health.HP -= damage
	ecs.World.Remove(projectile.Entity())
}
func NewLigtningAttack(ecs *ecs.ECS, param LightnigAtkParam) {
	startCol := param.StartCol

	for i := startCol; i >= 0 && i <= 7; i += param.Direction {
		entity := ecs.World.Create(
			component.Damage,
			component.GridPos,
			component.OnHit,
			component.Transient,
			component.Fx,
		)
		entry := ecs.World.Entry(entity)
		component.Damage.Set(entry, &component.DamageData{Damage: param.Damage})
		component.GridPos.Set(entry, &component.GridPosComponentData{Row: param.StartRow, Col: i})
		scrX, scrY := assets.GridCoord2Screen(param.StartRow, i)
		fxHeight := assets.LightningBolt.Bounds().Dy()
		anim1 := core.NewMovableImage(assets.LightningBolt, core.
			NewMovableImageParams().
			WithMoveParam(core.MoveParam{Sx: scrX - (float64(assets.TileWidth) / 2), Sy: scrY - float64(fxHeight)}))
		anim1.Done = func() {}
		component.Fx.Set(entry, &component.FxData{Animation: anim1})
		component.Transient.Set(entry, &component.TransientData{
			Start:    time.Now(),
			Duration: 200 * time.Millisecond,
		})
		component.OnHit.SetValue(entry, LightningBoltOnHitfunc)
	}
}

// cost 3 EN and cast lightning bolt
type LightingBoltCaster struct {
	Cost         int
	nextCooldown time.Time
	CoolDown     time.Duration
	Damage       int
}

func NewLightningBolCaster() *LightingBoltCaster {
	return &LightingBoltCaster{Cost: 300, Damage: 60, nextCooldown: time.Now(), CoolDown: 5 * time.Second}
}

func (l *LightingBoltCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= 300 {
		ensource.SetEn(en - l.Cost)
		query := donburi.NewQuery(
			filter.Contains(
				archetype.PlayerTag,
			),
		)

		playerId, ok := query.First(ecs.World)
		if !ok {
			return
		}
		gridPos := component.GridPos.Get(playerId)
		param := LightnigAtkParam{
			StartRow:  gridPos.Row,
			StartCol:  gridPos.Col + 1,
			Direction: 1,
			Actor:     playerId,
			Damage:    l.Damage,
		}
		NewLigtningAttack(ecs, param)
		l.nextCooldown = time.Now().Add(l.CoolDown)
	}
}
func (l *LightingBoltCaster) GetDamage() int {
	return l.Damage
}
func (l *LightingBoltCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\nShoots %d damage Lightning bolt instantly.\nCooldown %.1fs", l.Cost/100, l.Damage, l.CoolDown.Seconds())
}
func (l *LightingBoltCaster) GetName() string {
	return "LightningBolt"
}
func (l *LightingBoltCaster) GetCost() int {
	return l.Cost
}
func (l *LightingBoltCaster) GetIcon() *ebiten.Image {
	return assets.LightningIcon
}
func (l *LightingBoltCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *LightingBoltCaster) GetCooldownDuration() time.Duration {
	return l.CoolDown
}
