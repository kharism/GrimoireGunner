package attack

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// cost 2 EN and create a drill construct
type DrillCaster struct {
	Cost         int
	Damage       int
	nextCooldown time.Time
	CoolDown     time.Duration
	ModEntry     *loadout.CasterModifierData
	OnHit        component.OnAtkHit
}

func (l *DrillCaster) GetModifierEntry() *loadout.CasterModifierData {
	return l.ModEntry
}
func (l *DrillCaster) SetModifier(e *loadout.CasterModifierData) {
	if l.ModEntry != e && e.OnHit != nil {
		if l.OnHit == nil {
			l.OnHit = e.OnHit
		} else {
			l.OnHit = JoinOnAtkHit(l.OnHit, e.OnHit)
		}
	}
	l.ModEntry = e
}
func NewDrillCaster() *DrillCaster {
	return &DrillCaster{Cost: 200, nextCooldown: time.Now(), Damage: 150, CoolDown: 3 * time.Second, OnHit: PullbackOnHit}
}

func (l *DrillCaster) GetDescription() string {
	return fmt.Sprintf("Cost:%d EN\n%d Create Drill construct that moves forward.\nCooldown %.1fs", l.Cost/100, l.GetDamage(), l.GetCooldownDuration().Seconds())
}
func (l *DrillCaster) GetName() string {
	return "DrillConstruct"
}
func (l *DrillCaster) GetDamage() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.Damage + l.ModEntry.DamageModifier
	}
	return l.Damage
}
func (l *DrillCaster) GetCost() int {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		if l.Cost+l.ModEntry.CostModifier < 0 {
			return 0
		}
		return l.Cost + l.ModEntry.CostModifier
	}
	return l.Cost
}
func (l *DrillCaster) GetIcon() *ebiten.Image {
	return assets.DrillIcon
}
func (l *DrillCaster) GetCooldown() time.Time {
	return l.nextCooldown
}
func (l *DrillCaster) GetCooldownDuration() time.Duration {
	if l.ModEntry != nil {
		// mod := component.CasterModifier.Get(l.ModEntry)
		return l.CoolDown + l.ModEntry.CooldownModifer
	}
	return l.CoolDown
}
func (l *DrillCaster) ResetCooldown() {
	l.nextCooldown = time.Now()
}
func (l *DrillCaster) Cast(ensource loadout.ENSetGetter, ecs *ecs.ECS) {
	en := ensource.GetEn()
	if en >= l.GetCost() {
		ensource.SetEn(en - l.GetCost())
		l.nextCooldown = time.Now().Add(l.GetCooldownDuration())
		playerPos, _ := GetPlayerGridPos((ecs))
		if validMove(ecs, playerPos.Row, playerPos.Col+1) {
			drillCons := archetype.NewConstruct(ecs.World, assets.Drill)
			kk := ecs.World.Entry(*drillCons)

			gridPos1 := component.GridPos.Get(kk)
			gridPos1.Col = playerPos.Col + 1
			gridPos1.Row = playerPos.Row
			kk.AddComponent(component.EnemyRoutine)
			kk.AddComponent(component.TargetLocation)
			kk.AddComponent(component.Speed)
			component.Speed.Set(kk, &component.SpeedData{V: 2})
			tX, tY := assets.GridCoord2Screen(gridPos1.Row, 6)
			component.TargetLocation.Set(kk, &component.MoveTargetData{Tx: tX, Ty: tY})
			data := map[string]any{}

			ll := ecs.World.Create(component.GridPos, component.Damage, component.OnHit)
			dmgGrid := ecs.World.Entry(ll)
			gridPos2 := component.GridPos.Get(dmgGrid)
			gridPos2.Col = playerPos.Col + 2
			gridPos2.Row = playerPos.Row
			component.Damage.Set(dmgGrid, &component.DamageData{Damage: l.GetDamage()})
			component.OnHit.SetValue(dmgGrid, SingleHitProjectile)

			data["ALREADY_FIRED"] = false
			data["WARM_UP"] = time.Now()
			data["CURRENT_STRATEGY"] = ""
			data["MOVE_COUNT"] = 0
			data["CUR_DMG"] = l.GetDamage()
			data["DMG_GRID"] = dmgGrid
			component.EnemyRoutine.Set(kk, &component.EnemyRoutineData{Routine: DrillRoutine, Memory: data})

		}

		if l.ModEntry != nil {
			// l := component.PostAtkModifier.GetValue(l.ModEntry)
			if l.ModEntry.PostAtk != nil {
				l.ModEntry.PostAtk(ecs, ensource)
			}
		}
	}
}

// func ondrillhit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
// 	gridPos := component.GridPos.Get(projectile)
// 	damage := component.Damage.Get(projectile)
// 	component.EventQueue.AddEvent()
// 	SingleHitProjectile(ecs, projectile, receiver)
// }

type createNewDrillDmgGrid struct {
	time time.Time
}

func DrillRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	gridPos := component.GridPos.Get(entity)
	dmg := memory["CUR_DMG"].(int)
	dmgGrid := memory["DMG_GRID"].(*donburi.Entry)
	if gridPos.Col == 7 {
		ecs.World.Remove(entity.Entity())
		ecs.World.Remove(dmgGrid.Entity())
		return
	}
	if !ecs.World.Valid(dmgGrid.Entity()) {
		memory["CURRENT_STRATEGY"] = "WAIT"
		// memory["WARM_UP"] = time.Now().Add(200 * time.Millisecond)
	}
	if memory["CURRENT_STRATEGY"] == "" {
		memory["CURRENT_STRATEGY"] = "MOVE"
	}
	if memory["CURRENT_STRATEGY"] == "WAIT" {
		if waitTime, ok := memory["WARM_UP"].(time.Time); ok && waitTime.Before(time.Now()) {
			ll := ecs.World.Create(component.GridPos, component.Damage, component.OnHit)
			dmgGrid := ecs.World.Entry(ll)
			gridPos2 := component.GridPos.Get(dmgGrid)
			gridPos2.Col = gridPos.Col + 1
			gridPos2.Row = gridPos.Row
			component.Damage.Set(dmgGrid, &component.DamageData{Damage: dmg})
			component.OnHit.SetValue(dmgGrid, SingleHitProjectile)
			memory["CURRENT_STRATEGY"] = "MOVE"
			memory["DMG_GRID"] = dmgGrid
		}
	}
	if memory["CURRENT_STRATEGY"] == "MOVE" {
		if dmgGrid.HasComponent(component.GridPos) {
			cc := component.GridPos.Get(dmgGrid)
			cc.Col = gridPos.Col + 1
		}
		if !validMove(ecs, gridPos.Row, gridPos.Col+1) {
			component.Speed.Get(entity).V = 0
			component.Speed.Get(entity).Vx = 0
		}
		if component.Speed.Get(entity).Vx == 0 && validMove(ecs, gridPos.Row, gridPos.Col+1) {
			component.Speed.Get(entity).Vx = 2
		}
	}
}
