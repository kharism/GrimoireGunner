package scene

import (
	"fmt"
	"image"
	"image/color"

	// "image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/hanashi/core"
)

type StageSelect struct {
	data        *SceneData
	sm          *stagehand.SceneDirector[*SceneData]
	LevelLayout *Level
	musicPlayer *assets.AudioPlayer
	loopMusic   bool
}
type LevelNode struct {
	Id string
	// Trigger     stagehand.SceneTransitionTrigger
	Icon          *ebiten.Image
	Tier          int
	SelectedStage nextStagePicker
	// CombatDecor CombatSceneDecorator
	NextNode []*LevelNode
	ScrX     float64
	SrcY     float64
}

type nextStagePicker interface {
	DecorSceneData(*SceneData)
	GetNextStageTrigger() stagehand.SceneTransitionTrigger
	GetIcon() *ebiten.Image
}
type Level struct {
	Root      *LevelNode
	BossLevel *LevelNode
}

func GenerateLayout1() *Level {
	var LevelLayout1 = &Level{
		Root: &LevelNode{
			Id:            "0",
			Tier:          0,
			SelectedStage: NewCombatNextStage(level1Decorator7),
			Icon:          assets.BattleIcon,
		},
	}
	CurNode1 := &LevelNode{Id: "1", Tier: 1, SelectedStage: NewCombatNextStage(nil), NextNode: []*LevelNode{}, Icon: assets.BattleIcon}
	CurNode2 := &LevelNode{Id: "2", Tier: 1, SelectedStage: NewCombatNextStage(nil), NextNode: []*LevelNode{}, Icon: assets.BattleIcon}
	LevelLayout1.Root.NextNode = []*LevelNode{
		CurNode1, CurNode2,
	}
	for i := 0; i < 4; i++ {
		NewNodeA := &LevelNode{Id: fmt.Sprintf("%d", 2*i+3), Icon: assets.BattleIcon, Tier: CurNode1.Tier + 1, SelectedStage: NewCombatNextStage(nil), NextNode: []*LevelNode{}}
		NewNodeB := &LevelNode{Id: fmt.Sprintf("%d", 2*i+4), Icon: assets.BattleIcon, Tier: CurNode2.Tier + 1, SelectedStage: NewCombatNextStage(nil), NextNode: []*LevelNode{}}
		CurNode1.NextNode = append(CurNode1.NextNode, NewNodeA)
		CurNode2.NextNode = append(CurNode2.NextNode, NewNodeB)
		if i == 3 {
			if rand.Int()%2 == 0 {
				NewNodeA.SelectedStage = &RestSceneNextStage{}
				NewNodeB.SelectedStage = &WorkshopSceneNextStage{}
			} else {
				NewNodeA.SelectedStage = &WorkshopSceneNextStage{}
				NewNodeB.SelectedStage = &RestSceneNextStage{}
			}
			NewNodeA.Icon = NewNodeA.SelectedStage.GetIcon()
			NewNodeB.Icon = NewNodeB.SelectedStage.GetIcon()
		}
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
	BossNode := LevelNode{Id: "AA", Tier: CurNode1.Tier + 1, Icon: assets.BattleIcon, SelectedStage: NewCombatNextStage(level1Decorator9)}
	CurNode1.NextNode = append(CurNode1.NextNode, &BossNode)
	CurNode2.NextNode = append(CurNode2.NextNode, &BossNode)
	LevelLayout1.BossLevel = &BossNode
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

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

var StartPositionX = 10.0
var XDist = 40 * 3
var YDist = 40 * 3
var StartPositionY = 80.0
var tiers [][]*LevelNode

func drawLine(dest *ebiten.Image, xStart, yStart, xEnd, yEnd float32, lColor color.Color) {
	var path vector.Path
	path.MoveTo(xStart, yStart)
	path.LineTo(xEnd, yEnd)
	path.Close()
	var vs []ebiten.Vertex
	var is []uint16

	op := &vector.StrokeOptions{}
	op.Width = 5
	op.LineJoin = vector.LineJoinRound
	R, G, B, _ := lColor.RGBA()

	vs, is = path.AppendVerticesAndIndicesForStroke(nil, nil, op)
	for i := range vs {
		// vs[i].DstX = (vs[i].DstX + float32(x))
		// vs[i].DstY = (vs[i].DstY + float32(y))
		// vs[i].SrcX = 1
		// vs[i].SrcY = 1
		vs[i].ColorR = float32(R) / float32(256)
		vs[i].ColorG = float32(G) / float32(256)
		vs[i].ColorB = float32(B) / float32(256)
		vs[i].ColorA = 1
	}
	op2 := &ebiten.DrawTrianglesOptions{}
	op2.AntiAlias = true
	op2.FillRule = ebiten.NonZero

	dest.DrawTriangles(vs, is, whiteSubImage, op2)
}
func getIndexOfLevelNode(hay []*LevelNode, ss *LevelNode) int {
	for idx, c := range hay {
		if c == ss {
			return idx
		}
	}
	return -1
}
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
			var lineColor color.Color
			if d == r.data.CurrentLevel || contains(r.data.CurrentLevel.NextNode, d) {
				lineColor = color.White
			} else {
				lineColor = color.RGBA{128, 128, 128, 255}
			}
			for _, e := range d.NextNode {
				idx3 := getIndexOfLevelNode(tiers[idx1+1], e)
				drawLine(screen,
					float32(StartPositionX+float64(XDist*idx1))+40,
					float32(StartPositionY+float64(YDist*idx2))+40,
					float32(StartPositionX+float64(XDist*(idx1+1)))+40,
					float32(StartPositionY+float64(YDist*idx3))+40,
					lineColor,
				)
			}
			transform := ebiten.GeoM{}
			transform.Scale(3, 3)
			transform.Translate(StartPositionX+float64(XDist*idx1), StartPositionY+float64(YDist*idx2))
			if d == r.data.CurrentLevel || contains(r.data.CurrentLevel.NextNode, d) {
				opt := ebiten.DrawImageOptions{
					GeoM: transform,
				}
				d.ScrX = StartPositionX + float64(XDist*idx1)
				d.SrcY = StartPositionY + float64(YDist*idx2)
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
	if r.loopMusic && !r.musicPlayer.AudioPlayer().IsPlaying() {
		r.musicPlayer.AudioPlayer().Rewind()
		r.musicPlayer.AudioPlayer().Play()
	}
	if r.musicPlayer != nil {
		r.musicPlayer.Update()
	}
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
		r.sm.ProcessTrigger(pickedStage.SelectedStage.GetNextStageTrigger())
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
	if curLevel != nil {
		if len(tiers[curLevel.Tier]) > 1 && tiers[curLevel.Tier][1] == curLevel {
			if len(curLevel.NextNode) == 1 && curLevel.NextNode[0].Id != "AA" {
				startY += float64(YDist)
			}
		}
		stagePick.SetPos(StartPositionX+float64(XDist*(curLevel.Tier+1)), startY)
	} else {
		r.data.CurrentLevel = &LevelNode{
			NextNode: []*LevelNode{r.data.LevelLayout.Root},
		}
		stagePick.SetPos(StartPositionX, startY)
	}
	r.loopMusic = true
	var err error
	if r.musicPlayer == nil {
		r.musicPlayer, err = assets.NewAudioPlayer(assets.IntermissionMusic, assets.TypeMP3)
		if err != nil {
			fmt.Println(err.Error())
		}
		r.musicPlayer.AudioPlayer().SetPosition(r.data.MusicSeek)
		r.musicPlayer.AudioPlayer().Play()
		// set interfaces for sfx
	}

}
func (s *StageSelect) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *StageSelect) Unload() *SceneData {
	if pickedStage != nil {
		RegisterCombatClear = false
		s.data.CurrentLevel = pickedStage
		pickedStage.SelectedStage.DecorSceneData(s.data)
		// s.data.SceneDecor = pickedStage.CombatDecor
	}
	s.loopMusic = false
	s.data.MusicSeek = s.musicPlayer.AudioPlayer().Position()
	s.musicPlayer.AudioPlayer().Rewind()
	s.musicPlayer.AudioPlayer().Pause()
	// s.data.CurrentLevel = tiers[s.data.CurrentLevel.Tier][stageCursorIndex]
	return s.data
}
