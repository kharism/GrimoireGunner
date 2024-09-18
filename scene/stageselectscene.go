package scene

import (
	"fmt"
	"image/color"

	// "image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/hanashi/core"
)

type StageSelect struct {
	data        *SceneData
	sm          *stagehand.SceneDirector[*SceneData]
	LevelLayout *Level
}
type LevelNode struct {
	Id string

	Icon      *ebiten.Image
	Tier      int
	Decorator CombatSceneDecorator
	NextNode  []*LevelNode
}
type Level struct {
	Root *LevelNode
}

func GenerateLayout1() *Level {
	var LevelLayout1 = &Level{
		Root: &LevelNode{
			Id:        "0",
			Tier:      0,
			Decorator: nil,
			Icon:      assets.BattleIcon,
		},
	}
	CurNode1 := &LevelNode{Id: "1", Tier: 1, Decorator: RandDecorator(), NextNode: []*LevelNode{}, Icon: assets.BattleIcon}
	CurNode2 := &LevelNode{Id: "2", Tier: 1, Decorator: RandDecorator(), NextNode: []*LevelNode{}, Icon: assets.BattleIcon}
	LevelLayout1.Root.NextNode = []*LevelNode{
		CurNode1, CurNode2,
	}
	for i := 0; i < 4; i++ {
		NewNodeA := &LevelNode{Id: fmt.Sprintf("%d", 2*i+3), Icon: assets.BattleIcon, Tier: CurNode1.Tier + 1, Decorator: RandDecorator(), NextNode: []*LevelNode{}}
		NewNodeB := &LevelNode{Id: fmt.Sprintf("%d", 2*i+4), Icon: assets.BattleIcon, Tier: CurNode2.Tier + 1, Decorator: RandDecorator(), NextNode: []*LevelNode{}}
		CurNode1.NextNode = append(CurNode1.NextNode, NewNodeA)
		CurNode2.NextNode = append(CurNode2.NextNode, NewNodeB)
		if rand.Int()%10 <= 3 {
			// add cros section
			if rand.Int()%2 == 0 {
				//create line from lower path to uper path
				CurNode2.NextNode = append(CurNode2.NextNode, NewNodeA)
			} else {
				//create line from upper path to lower path
				CurNode1.NextNode = append(CurNode1.NextNode, NewNodeB)
			}
		}
		CurNode1 = NewNodeA
		CurNode2 = NewNodeB
	}
	BossNode := LevelNode{Id: "AA", Tier: CurNode1.Tier + 1, Icon: assets.BattleIcon}
	CurNode1.NextNode = append(CurNode1.NextNode, &BossNode)
	CurNode2.NextNode = append(CurNode2.NextNode, &BossNode)
	//upLine :=
	return LevelLayout1
}
func contains(haystack []*LevelNode, needle *LevelNode) bool {
	exist := false
	for _, c := range haystack {
		if c == needle {
			return true
		}
	}
	return exist
}
func TraverseBfs(node *LevelNode, tiers *[][]*LevelNode) {
	if len(*tiers) <= node.Tier {
		*tiers = append(*tiers, []*LevelNode{node})
	} else {
		if !contains((*tiers)[node.Tier], node) {
			(*tiers)[node.Tier] = append((*tiers)[node.Tier], node)
		}
	}
	for _, c := range node.NextNode {
		TraverseBfs(c, tiers)
	}
}

var StartPositionX = 10.0
var XDist = 40 * 3
var YDist = 40 * 3
var StartPositionY = 80.0
var tiers [][]*LevelNode

func (r *StageSelect) Draw(screen *ebiten.Image) {

	// fmt.Println(tiers)
	bg := ebiten.NewImage(1024, 600)
	bg.Fill(color.RGBA{R: 0x21, G: 0x43, B: 0x58, A: 255})
	screen.DrawImage(bg, &ebiten.DrawImageOptions{})
	textTranslate := ebiten.GeoM{}
	textTranslate.Translate(512, 50)

	text.Draw(screen, "Select Stage", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: textTranslate,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	})
	//TODO: bugfix display
	for idx1, c := range tiers {
		for idx2, d := range c {

			transform := ebiten.GeoM{}
			transform.Scale(3, 3)
			transform.Translate(StartPositionX+float64(XDist*idx1), StartPositionY+float64(YDist*idx2))
			if d == r.data.CurrentLevel || contains(r.data.CurrentLevel.NextNode, d) {
				opt := ebiten.DrawImageOptions{
					GeoM: transform,
				}
				screen.DrawImage(d.Icon, &opt)
			} else {
				opts := &ebiten.DrawRectShaderOptions{
					GeoM: transform,
				}
				opts.Images[0] = d.Icon
				iconBound := d.Icon.Bounds()
				// bounds := .Bounds()
				screen.DrawRectShader(iconBound.Dx(), iconBound.Dy(), assets.DarkerShader, opts)
			}

		}
	}
	stagePick.Draw(screen)
}
func (r *StageSelect) Update() error {

	curLevel := r.data.CurrentLevel
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		stageCursorIndex += 1
		if stageCursorIndex == len(curLevel.NextNode) {
			stageCursorIndex -= 1
		} else {
			CursYPos := StartPositionY + float64(YDist*stageCursorIndex)
			if len(tiers[curLevel.Tier]) > 1 && tiers[curLevel.Tier][1] == curLevel {
				if len(curLevel.NextNode) == 1 {
					CursYPos += float64(YDist)
				}
			}
			stagePick.AddAnimation(core.NewMoveAnimationFromParam(
				core.MoveParam{
					Tx:    StartPositionX + float64(XDist*(curLevel.Tier+1)),
					Ty:    CursYPos,
					Speed: 10,
				},
			))
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		stageCursorIndex -= 1
		if stageCursorIndex == -1 {
			stageCursorIndex += 1
		} else {
			CursYPos := StartPositionY + float64(YDist*stageCursorIndex)
			if len(tiers[curLevel.Tier]) > 1 && tiers[curLevel.Tier][1] == curLevel {
				if len(curLevel.NextNode) == 1 {
					CursYPos += float64(YDist)
				}

			}
			stagePick.AddAnimation(core.NewMoveAnimationFromParam(
				core.MoveParam{
					Tx:    StartPositionX + float64(XDist*(curLevel.Tier+1)),
					Ty:    CursYPos,
					Speed: 10,
				},
			))
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		pickedStage = r.data.CurrentLevel.NextNode[stageCursorIndex]
		r.sm.ProcessTrigger(TriggerToCombat)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		// pickedStage = tiers[r.data.CurrentLevel.Tier+1][stageCursorIndex]
		r.sm.ProcessTrigger(TriggerToCombat)
	}
	stagePick.Update()
	return nil
}

var StageSelectInstance = &StageSelect{}
var stageCursorIndex int
var stagePick *core.MovableImage
var pickedStage *LevelNode

func (r *StageSelect) Load(state *SceneData, manager stagehand.SceneController[*SceneData]) {
	r.sm = manager.(*stagehand.SceneDirector[*SceneData]) // This type assertion is important
	r.data = state
	stageCursorIndex = 0
	stagePick = nil
	pickedStage = nil
	tiers = [][]*LevelNode{}
	r.LevelLayout = state.LevelLayout
	TraverseBfs(r.LevelLayout.Root, &tiers)
	curLevel := r.data.CurrentLevel
	bounds := assets.CardPick.Bounds()
	stagePick = core.NewMovableImage(assets.CardPick,
		core.NewMovableImageParams().WithScale(&core.ScaleParam{
			Sx: float64(32*3) / float64(bounds.Dx()),
			Sy: float64(32*3) / float64(bounds.Dy()),
		}),
	)
	startY := StartPositionY
	if len(tiers[curLevel.Tier]) > 1 && tiers[curLevel.Tier][1] == curLevel {
		if len(curLevel.NextNode) == 1 {
			startY += float64(YDist)
		}
	}
	stagePick.SetPos(StartPositionX+float64(XDist*(curLevel.Tier+1)), startY)

}
func (s *StageSelect) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *StageSelect) Unload() *SceneData {
	if pickedStage != nil {
		RegisterCombatClear = false
		s.data.CurrentLevel = pickedStage
		s.data.SceneDecor = pickedStage.Decorator
	}
	// s.data.CurrentLevel = tiers[s.data.CurrentLevel.Tier][stageCursorIndex]
	return s.data
}
