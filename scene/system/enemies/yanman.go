package enemies

import (
	"math/rand"
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/hanashi/core"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func NewYanman(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Yanma)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	component.Health.Set(entry, &component.HealthData{HP: 500, MaxHP: 500, Name: "Yanman"})

	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	data[OPTION_LIST] = []donburi.Entity{}
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: YanmanRoutine, Memory: data})
}

const OPTION_LIST = "OPT_LIST"

func BambooAtk(ecs *ecs.ECS, entity *donburi.Entry) {
	gridPos := component.GridPos.Get(entity)
	now := time.Now()
	dmgGridEnt := ecs.World.Create(component.GridPos, component.Damage, component.OnHit, component.Transient)
	dmgGrid := ecs.World.Entry(dmgGridEnt)
	component.GridPos.Set(dmgGrid, gridPos)
	component.Damage.Set(dmgGrid, &component.DamageData{Damage: 50})
	component.OnHit.SetValue(dmgGrid, attack.SingleHitProjectile)
	component.Transient.Set(dmgGrid, &component.TransientData{
		Start:    now,
		Duration: 300 * time.Millisecond,
	})
	scrX, scrY := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
	scrX -= 50
	scrY -= 50
	fxBamboo := ecs.World.Create(component.Fx, component.Transient)
	fxBambooEntry := ecs.World.Entry(fxBamboo)
	bambooLance := core.NewMovableImage(assets.BambooLance,
		core.NewMovableImageParams().WithMoveParam(core.MoveParam{
			Sx: scrX,
			Sy: scrY,
		}),
	)
	component.Transient.Set(fxBambooEntry, &component.TransientData{
		Start:    now,
		Duration: 300 * time.Millisecond,
	})
	component.Fx.Set(fxBambooEntry, &component.FxData{
		Animation: bambooLance,
	})

}
func YanmanRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)

	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	if memory[CURRENT_STRATEGY] == "SLEEP" {
		optionList := memory[OPTION_LIST].([]donburi.Entity)
		newOptionList := []donburi.Entity{}
		for _, i := range optionList {
			if ecs.World.Valid(i) {
				newOptionList = append(newOptionList, i)
			}
		}
		if len(newOptionList) == 0 {
			memory[CURRENT_STRATEGY] = "MOVE"
		}
		memory[OPTION_LIST] = newOptionList
		playerPos, _ := attack.GetPlayerGridPos(ecs)
		if playerPos.Col == 0 {
			for i := 0; i < 4; i++ {
				targetGrid := ecs.World.Create(component.GridPos, component.GridTarget, component.Transient)
				targetGridEnt := ecs.World.Entry(targetGrid)
				component.GridPos.Set(targetGridEnt, &component.GridPosComponentData{
					Col: 0,
					Row: i,
				})
				component.Transient.Set(targetGridEnt, &component.TransientData{
					Start:            time.Now(),
					Duration:         900 * time.Millisecond,
					OnRemoveCallback: BambooAtk,
				})
			}
		}
	}
	gridPos := component.GridPos.Get(entity)
	if memory[CURRENT_STRATEGY] == "WAIT" {
		component.Sprite.Get(entity).Image = assets.Yanma
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			// memory[CURRENT_STRATEGY] = "MOVE"
			memory[CURRENT_STRATEGY] = "MOVE"
			memory[WARM_UP] = time.Now().Add(300 * time.Millisecond)
			// if gridPos.Col != 4 {
			// 	memory[CURRENT_STRATEGY] = "MOVE"
			// } else {
			// 	memory[CURRENT_STRATEGY] = "SUMMON_CONSTRUCT"
			// 	memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
			// }
		}
	}
	moveCount := memory[MOVE_COUNT].(int)
	if memory[CURRENT_STRATEGY] == "MOVE" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			if moveCount < 4 {
				scrPos := component.ScreenPos.Get(entity)
			checkMove:
				for {
					// rndMove := rand.Int() % 4
					newCol := 4 + rand.Int()%4
					newRow := rand.Int() % 4
					if validMove(ecs, newRow, newCol) {
						gridPos.Row = newRow
						gridPos.Col = newCol
						scrPos.Y = 0
						scrPos.X = 0
						memory[MOVE_COUNT] = moveCount + 1
						memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)

						break checkMove
					}

				}
			} else {
				// warmup for attack
				memory[MOVE_COUNT] = 0
				moveCount = 0
				// ii := rand.Int() % 2
				memory[CURRENT_STRATEGY] = "ATTACK"
				component.Sprite.Set(entity, &component.SpriteData{Image: assets.YanmaAttack})
			}
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK" {
		rows := []int{0, 1, 2, 3}
		rand.Shuffle(4, func(i, j int) {
			rows[i], rows[j] = rows[j], rows[i]
		})
		optionList := []donburi.Entity{}
		for i := 0; i < 2; i++ {
			jj := GenerateYanmaOption(ecs, gridPos.Col-1, rows[i], dmg, i)
			optionList = append(optionList, *jj)
		}
		memory[OPTION_LIST] = optionList
		memory[CURRENT_STRATEGY] = "SLEEP"
	}
}
func GenerateYanmaOption(ecs *ecs.ECS, col, row, damage, multiplier int) *donburi.Entity {
	yanmaOption := archetype.NewNPC(ecs.World, assets.YanmaOption)
	entry := ecs.World.Entry(*yanmaOption)
	component.Health.Set(entry, &component.HealthData{
		HP:    100,
		MaxHP: 100,
		Name:  "YanmaOption",
	})
	component.GridPos.Set(entry, &component.GridPosComponentData{
		Row: row,
		Col: col,
	})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = damage
	data[MULTIPLIER] = multiplier
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: YanmaOptionRoutine, Memory: data})
	return yanmaOption
}

const MULTIPLIER = "MUL"

func YanmaOptionRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)
	mul := memory[MULTIPLIER].(int)
	gridpos := component.GridPos.Get(entity)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "ATTACK"
		memory[WARM_UP] = time.Now().Add(300 * time.Millisecond)
	}
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "ATTACK"
			memory[WARM_UP] = time.Now().Add(300 * time.Millisecond)
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			atkCount := memory[MOVE_COUNT].(int)
			if atkCount < 3 {
				archetype.NewProjectile(ecs.World, archetype.ProjectileParam{
					Vx:     -5,
					Vy:     0,
					Col:    gridpos.Col - 1,
					Row:    gridpos.Row,
					Damage: dmg,
					Sprite: assets.Projectile1,
					OnHit:  attack.SingleHitProjectile,
				})
				memory[CURRENT_STRATEGY] = "WAIT"
				memory[WARM_UP] = time.Now().Add(150 * time.Duration(mul) * time.Millisecond)
				memory[MOVE_COUNT] = atkCount + 1
				component.EnemyRoutine.Get(entity).Memory = memory
			} else {
				//become projectile
				memory[CURRENT_STRATEGY] = "KAMIKAZE"
				// memory[WARM_UP] = time.Now().Add(2000 * time.Millisecond)
				entity.AddComponent(component.Speed)
				entity.AddComponent(component.TargetLocation)
				entity.AddComponent(component.KamikazeTag)
				entity.AddComponent(component.Damage)
				component.Speed.Set(entity, &component.SpeedData{
					Vx: -5,
					Vy: 0,
				})
				component.Damage.Set(entity, &component.DamageData{Damage: dmg})
			}

		}

	}
	if memory[CURRENT_STRATEGY] == "KAMIKAZE" {

	}
}
