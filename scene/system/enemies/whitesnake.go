package enemies

import (
	"math/rand"
	"time"

	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func NewWhiteSnake(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Whitesnake)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	entry.AddComponent(component.Shader)
	component.Shader.Set(entry, nil)
	component.Health.Set(entry, &component.HealthData{HP: 1200, MaxHP: 1200, Name: "WhiteSnake"})

	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})
	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	data[OPTION_LIST] = []donburi.Entity{}
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: WhiteSnakeRoutine, Memory: data})
}

var filterEnemy = donburi.NewQuery(
	filter.Contains(
		component.EnemyTag,
	),
)

func WhiteSnakeRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)
	health := component.Health.Get(entity)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
	}
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			i := rand.Int() % 4
			var el component.Elemental
			switch i {
			case 0:
				el = component.FIRE
				memory[CURRENT_STRATEGY] = "FIRE_BREATH_1"
			case 1:
				el = component.ELEC
				memory[CURRENT_STRATEGY] = "ELECTROCUTE_1"
			case 2:
				el = component.WATER
				memory[CURRENT_STRATEGY] = "ICE_SPICKE_1"
			case 3:
				el = component.WOOD
				memory[CURRENT_STRATEGY] = "SUMMON_BUZZER"
			}
			health.Element = el
			shader := assets.Element2Shader(el)
			if shader != nil {
				component.Shader.Set(entity, shader)
			}
			memory[WARM_UP] = time.Now().Add(700 * time.Millisecond)
		}

	}
	playerGrid, _ := attack.GetPlayerGridPos(ecs)
	if playerGrid == nil {
		return
	}
	demonPos := component.GridPos.Get(entity)
	if playerGrid == nil {
		return
	}
	if memory[CURRENT_STRATEGY] == "FIRE_BREATH_1" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.WhitesnakeWarmup
			//move to front

			tempRow := playerGrid.Row
			tempCol := 4
			for {
				if validMove(ecs, tempRow, tempCol) {
					demonPos.Row = tempRow
					demonPos.Col = tempCol
					scrPos := component.ScreenPos.Get(entity)
					scrPos.X = 0
					scrPos.Y = 0
					break
				} else if tempCol < 8 {
					tempCol += 1
				} else {
					tempRow += 1
					tempCol = 0
				}
			}
			memory[CURRENT_STRATEGY] = "FIRE_BREATH_2"
			memory[WARM_UP] = time.Now().Add(1 * time.Second)
		}

	}
	if memory[CURRENT_STRATEGY] == "FIRE_BREATH_2" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			PyroEyesAttack(ecs, entity)
			memory[CURRENT_STRATEGY] = "WAIT"
			memory[WARM_UP] = time.Now().Add(600 * time.Millisecond)
			memory[CUR_DMG] = dmg + 10
		}
	}
	if memory[CURRENT_STRATEGY] == "ELECTROCUTE_1" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			component.Sprite.Get(entity).Image = assets.WhitesnakeWarmup
			//move to front
			// demonPos := component.GridPos.Get(entity)
			tempRow := playerGrid.Row
			tempCol := 4
			for {
				if validMove(ecs, tempRow, tempCol) {
					demonPos.Row = tempRow
					demonPos.Col = tempCol
					scrPos := component.ScreenPos.Get(entity)
					scrPos.X = 0
					scrPos.Y = 0
					break
				} else if tempCol < 8 {
					tempCol += 1
				} else {
					tempRow += 1
					tempCol = 0
				}
			}
			memory[CURRENT_STRATEGY] = "ELECTROCUTE_2"
			memory[WARM_UP] = time.Now().Add(700 * time.Millisecond)
		}
	}
	if memory[CURRENT_STRATEGY] == "ELECTROCUTE_2" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			rows := []int{0, 1, 2, 3}
			rand.Shuffle(4, func(i, j int) {
				rows[i], rows[j] = rows[j], rows[i]
			})
			shootTime := time.Now().Add(1 * time.Second)

			for i := 0; i < 3; i++ {
				entity := archetype.NewProjectile(ecs.World, archetype.ProjectileParam{
					Vx:     0,
					Vy:     0,
					Col:    demonPos.Col - 1,
					Row:    rows[i],
					Damage: dmg,
					Sprite: assets.ElecSphere,
					OnHit:  attack.SingleHitProjectile,
				})
				jj := &moveDagger{entry: ecs.World.Entry(*entity), Time: shootTime}
				component.EventQueue.AddEvent(jj)
			}
			memory[CURRENT_STRATEGY] = "WAIT"
			memory[WARM_UP] = time.Now().Add(600 * time.Millisecond)
		}
	}
	if memory[CURRENT_STRATEGY] == "SUMMON_BUZZER" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			enemyCount := filterEnemy.Count(ecs.World)
			// filterEnemy.Each(ecs.World, func(e *donburi.Entry) {
			// 	hh := component.Health.Get(e)
			// 	fmt.Println(hh.Name)
			// })
			if enemyCount <= 1 {
				tempRow := playerGrid.Row
				tempCol := 6
				for {
					if validMove(ecs, tempRow, tempCol) {
						demonPos.Row = tempRow
						demonPos.Col = tempCol
						scrPos := component.ScreenPos.Get(entity)
						scrPos.X = 0
						scrPos.Y = 0
						break
					} else if tempCol < 8 {
						tempCol += 1
					} else {
						tempRow += 1
						tempCol = 0
					}

				}
				memory[CURRENT_STRATEGY] = "SUMMON_DRGNFLY"
				memory[WARM_UP] = time.Now().Add(1 * time.Second)
				NewMiniBuzzer(ecs, demonPos.Col-1, demonPos.Row)
			} else {
				memory[CURRENT_STRATEGY] = "SUMMON_DRGNFLY"
			}

		}
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
			memory[CURRENT_STRATEGY] = "WAIT"
		}
		memory[OPTION_LIST] = newOptionList
		if playerGrid.Col == 0 {
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
	if memory[CURRENT_STRATEGY] == "SUMMON_DRGNFLY" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			enemyCount := filterEnemy.Count(ecs.World)
			if enemyCount <= 2 {
				rows := []int{0, 1, 2, 3}
				rand.Shuffle(4, func(i, j int) {
					rows[i], rows[j] = rows[j], rows[i]
				})

				optionList := []donburi.Entity{}
				for i := 0; i < 2; i++ {
					jj := GenerateYanmaOption(ecs, demonPos.Col-1, rows[i], dmg, i)
					optionList = append(optionList, *jj)
				}
				memory[OPTION_LIST] = optionList
				memory[OPTION_LIST] = optionList
				memory[CURRENT_STRATEGY] = "SLEEP"
			}

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
						demonPos.Row = newRow
						demonPos.Col = newCol
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
				memory[CURRENT_STRATEGY] = "WAIT"

			}
		}
	}
	if memory[CURRENT_STRATEGY] == "ICE_SPICKE_1" {
		playerPos := playerGrid
		now := time.Now()
		for col := playerPos.Col - 1; col <= playerPos.Col+1; col++ {
			for row := playerPos.Row - 1; row <= playerPos.Row+1; row++ {
				if col < 0 || col > 7 || row < 0 || row > 3 {
					continue
				}
				target := ecs.World.Create(component.GridPos, component.Damage, component.GridTarget, component.Transient)
				t := ecs.World.Entry(target)
				component.GridPos.Set(t, &component.GridPosComponentData{Col: col, Row: row})
				component.Damage.Set(t, &component.DamageData{Damage: dmg})
				component.Transient.Set(t, &component.TransientData{
					Start:            now,
					Duration:         200 * time.Millisecond,
					OnRemoveCallback: CreateAvalance,
				})
				// component.GridPos.Set(t,&component.GridPosComponentData{Col:col,Row: row})
			}
		}
		memory[CURRENT_STRATEGY] = "MOVE"
		memory[WARM_UP] = time.Now().Add(800 * time.Millisecond)
	}
}
