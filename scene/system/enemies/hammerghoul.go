package enemies

import (
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func NewHammerghoul(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Hammerghoul)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 400, Name: "HammerGhoul"})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[CURRENT_STRATEGY] = ""
	data[WARM_UP] = nil
	data[CUR_DMG] = 50
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: HammerGhoulRoutine, Memory: data})
}

func HammerGhoulRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(3 * time.Second)
	}
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "WARMUP"
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
			component.Sprite.Get(entity).Image = assets.HammerghoulWarmup
		}
	}
	if memory[CURRENT_STRATEGY] == "WARMUP" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.HammerghoulCooldown
			memory[CURRENT_STRATEGY] = ""
			gridPos := component.GridPos.Get(entity)
			startCol := gridPos.Col - 1
			for ; startCol >= 0; startCol-- {
				ff := time.Duration(200 * (gridPos.Col - 1 - startCol))
				CreateMovingExplosion(ecs, dmg, startCol, gridPos.Row, time.Now().Add(ff*time.Millisecond))
			}

		}
	}
}
func OnAtkHitExplosion(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	gg := component.Fx.Get(projectile)
	dmg := component.Damage.Get(projectile).Damage
	component.Health.Get(receiver).HP -= dmg
	anim := gg.Animation.(*core.AnimatedImage)
	anim.Done()
}

type createExplosionQueue struct {
	row    int
	col    int
	time   time.Time
	damage int
}

func (c *createExplosionQueue) Execute(ecs *ecs.ECS) {
	scrX, scrY := assets.GridCoord2Screen(c.row, c.col)
	hitbox1 := ecs.World.Create(component.Damage, component.GridPos, component.OnHit, component.Fx)
	entry1 := ecs.World.Entry(hitbox1)
	anim := assets.NewExplosionAnim(assets.SpriteParam{
		ScreenX: scrX - float64(assets.TileWidth)/2,
		ScreenY: scrY - 75,
		Modulo:  5,
	})

	component.Damage.Set(entry1, &component.DamageData{
		Damage: c.damage,
	})
	component.GridPos.Set(entry1, &component.GridPosComponentData{Row: c.row, Col: c.col})
	anim.Done = func() {
		//CreateMovingExplosion(ecs, col-1, row)
		ecs.World.Remove(entry1.Entity())
	}
	component.OnHit.SetValue(entry1, OnAtkHitExplosion)
	component.Fx.Set(entry1, &component.FxData{Animation: anim})
}
func (c *createExplosionQueue) GetTime() time.Time {
	return c.time
}
func CreateMovingExplosion(ecs *ecs.ECS, dmg, col, row int, time time.Time) {
	if col == -1 {
		return
	}
	j := &createExplosionQueue{row: row, col: col, time: time, damage: dmg}
	component.EventQueue.AddEvent(j)
}
