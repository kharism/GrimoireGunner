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

func NewHealslime(ecs *ecs.ECS, col, row int) {
	entity := archetype.NewNPC(ecs.World, assets.Slime)
	entry := ecs.World.Entry(*entity)
	entry.AddComponent(component.EnemyTag)
	entry.AddComponent(component.Elements)
	component.Health.Set(entry, &component.HealthData{HP: 300, MaxHP: 300, Name: "HealSlime", Element: component.WOOD})
	component.GridPos.Set(entry, &component.GridPosComponentData{Row: row, Col: col})
	component.ScreenPos.Set(entry, &component.ScreenPosComponentData{})

	data := map[string]any{}
	data[ALREADY_FIRED] = false
	data[WARM_UP] = nil
	data[CURRENT_STRATEGY] = ""
	data[MOVE_COUNT] = 0
	data[CUR_DMG] = 50
	component.EnemyRoutine.Set(entry, &component.EnemyRoutineData{Routine: SlimeRoutine, Memory: data})
}

var MOVE_COUNT = "MOVE_COUNT"

func SlimeRoutine(ecs *ecs.ECS, entity *donburi.Entry) {
	memory := component.EnemyRoutine.Get(entity).Memory
	dmg := memory[CUR_DMG].(int)
	if memory[CURRENT_STRATEGY] == "" {
		memory[CURRENT_STRATEGY] = "WAIT"
		memory[WARM_UP] = time.Now().Add(1 * time.Second)
	}
	gridPos := component.GridPos.Get(entity)
	if memory[CURRENT_STRATEGY] == "WAIT" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "MOVE"
		}
	}
	if memory[CURRENT_STRATEGY] == "MOVE" {
		component.Sprite.Set(entity, &component.SpriteData{Image: assets.Slime})
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			moveCount := memory[MOVE_COUNT].(int)
			if moveCount < 4 {

				scrPos := component.ScreenPos.Get(entity)
			checkMove:
				for {
					rndMove := rand.Int() % 4
					switch rndMove {
					case 0:
						if validMove(ecs, gridPos.Row-1, gridPos.Col) {
							gridPos.Row -= 1
							scrPos.Y -= float64(assets.TileHeight)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					case 1:
						if validMove(ecs, gridPos.Row, gridPos.Col-1) && gridPos.Col-1 >= 4 {
							gridPos.Col -= 1
							scrPos.X -= float64(assets.TileWidth)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					case 2:
						if validMove(ecs, gridPos.Row+1, gridPos.Col) {
							gridPos.Row += 1
							scrPos.Y += float64(assets.TileHeight)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					case 3:
						if validMove(ecs, gridPos.Row, gridPos.Col+1) {
							gridPos.Col += 1
							scrPos.X += float64(assets.TileWidth)
							memory[MOVE_COUNT] = moveCount + 1
							memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
							break checkMove
						}
					}
				}
			} else {
				memory[MOVE_COUNT] = 0
				memory[CURRENT_STRATEGY] = "ATTACK"
				component.Sprite.Set(entity, &component.SpriteData{Image: assets.Slime2})
				memory[WARM_UP] = time.Now().Add(1000 * time.Millisecond)
			}
		}
	}
	if memory[CURRENT_STRATEGY] == "ATTACK" {
		now := time.Now()
		component.Sprite.Set(entity, &component.SpriteData{Image: assets.Slime2})
		for i := 0; i < 10; i++ {
			component.EventQueue.AddEvent(NewSlimeShoot(now.Add(time.Duration(i*300)*time.Millisecond),
				dmg, gridPos.Col, gridPos.Row))
		}
		memory[WARM_UP] = time.Now().Add(1200 * time.Millisecond)
		memory[CURRENT_STRATEGY] = "COOLDOWN"
	}
	if memory[CURRENT_STRATEGY] == "COOLDOWN" {
		if waitTime, ok := memory[WARM_UP].(time.Time); ok && waitTime.Before(time.Now()) {
			memory[CURRENT_STRATEGY] = "MOVE"
			ll := component.Health.Get(entity)
			ll.HP += 50
			if ll.HP > ll.MaxHP {
				ll.HP = ll.MaxHP
			}

			fxEntity := ecs.World.Create(component.Fx, component.Transient)

			fx := ecs.World.Entry(fxEntity)

			x, y := assets.GridCoord2Screen(gridPos.Row, gridPos.Col)
			x -= 50
			y -= 100

			dmg += 10
			memory[CUR_DMG] = dmg + 10

			anim := core.NewMovableImage(assets.HealFx, core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: x, Sy: y}))
			component.Fx.Set(fx, &component.FxData{Animation: anim})
			component.Transient.Set(fx, &component.TransientData{Start: time.Now(), Duration: 500 * time.Millisecond})
			memory[WARM_UP] = time.Now().Add(500 * time.Millisecond)
		}
	}

}

// shoot something to player. Col and row are the start of the bullet
func NewSlimeShoot(shTime time.Time, damage, col, row int) *SlimeShoot {
	return &SlimeShoot{time: shTime, StartCol: col, StartRow: row, Damage: damage}
}

type SlimeShoot struct {
	time     time.Time
	StartCol int
	StartRow int
	Damage   int
}

func (s *SlimeShoot) Execute(ecs *ecs.ECS) {
	gridData, _ := attack.GetPlayerGridPos(ecs)
	if gridData == nil {
		return
	}
	// fmt.Println(gridData.Col, gridData.Row)
	bullet := ecs.World.Create(component.Fx)
	scrX, srcY := assets.GridCoord2Screen(s.StartRow, s.StartCol)
	bomb := assets.Bomb1
	bombBound := bomb.Bounds()
	scrX -= float64(assets.TileWidth) / 2
	srcY -= float64(bombBound.Dy())
	targX, targY := assets.GridCoord2Screen(gridData.Row, gridData.Col)
	targX -= float64(assets.TileWidth) / 2
	targY -= float64(bombBound.Dy())
	targRow := gridData.Row
	targCol := gridData.Col
	anim := core.NewMovableImage(
		assets.Bomb1,
		core.NewMovableImageParams().WithMoveParam(core.MoveParam{
			Sx:    scrX,
			Sy:    srcY,
			Speed: 10,
		}),
	)
	anim.Done = func() {
		Damage := s.Damage
		damageGrid := ecs.World.Create(component.GridPos, component.Elements, component.Damage, component.Transient, component.OnHit)
		dd := ecs.World.Entry(damageGrid)
		component.GridPos.Set(dd, &component.GridPosComponentData{Col: targCol, Row: targRow})
		component.Damage.Set(dd, &component.DamageData{
			Damage: Damage,
		})
		component.Elements.SetValue(dd, component.WOOD)
		component.OnHit.SetValue(dd, onHealSlimeHit)
		component.Transient.Set(dd, &component.TransientData{
			Start:    time.Now(),
			Duration: 100 * time.Millisecond,
		})
		ecs.World.Remove(bullet)
	}
	anim.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{Tx: targX, Ty: targY, Speed: 5}))
	component.Fx.Set(ecs.World.Entry(bullet), &component.FxData{Animation: anim})
}
func onHealSlimeHit(ecs *ecs.ECS, projectile, receiver *donburi.Entry) {
	damage := component.Damage.Get(projectile).Damage
	health := component.Health.Get(receiver)
	health.HP -= damage
}
func (s *SlimeShoot) GetTime() time.Time {
	return s.time
}
