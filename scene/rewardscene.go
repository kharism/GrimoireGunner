package scene

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"
)

type RewardScene struct {
	data        *SceneData
	sm          *stagehand.SceneDirector[*SceneData]
	musicPlayer *assets.AudioPlayer
	loopMusic   bool
}

func (r *RewardScene) Update() error {
	if r.loopMusic && !r.musicPlayer.AudioPlayer().IsPlaying() {
		r.musicPlayer.AudioPlayer().Rewind()
		r.musicPlayer.AudioPlayer().Play()
	}
	if r.musicPlayer != nil {
		r.musicPlayer.Update()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		r.musicPlayer.QueueSFX(assets.CursorFx)
		currentPick = (currentPick + 1) % 3
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		r.musicPlayer.QueueSFX(assets.CursorFx)
		if currentPick == 0 {
			currentPick = 2
		} else {
			currentPick = (currentPick - 1) % 3
		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		r.musicPlayer.QueueSFX(assets.CursorFx)
		for i := 0; i < len(cards); i++ {
			if i == currentPick {
				moveAnim := core.NewMoveAnimationFromParam(core.MoveParam{Tx: CardStartX + 800, Ty: CardStartY, Speed: 20})
				cards[i].Done = func() {
					// trigger event to go to
					// just to test manually assign to subweapon1 slot0
					// r.data.SubLoadout1[0] = casterPick[currentPick]
					if cc, ok := casterPick[currentPick].(OnAquireDoer); ok {
						cc.OnAquireDo(r.data)
					}
					r.data.Inventory = append(r.data.Inventory, casterPick[currentPick])
					r.sm.ProcessTrigger(TriggerToCombat)
				}
				cards[i].AddAnimation(moveAnim)

			} else {
				cards[i].AddAnimation(core.NewMoveAnimationFromParam(
					core.MoveParam{Tx: -80, Ty: CardStartY, Speed: 20},
				))
			}
			cards[i].Update()
		}
	}
	for i := 0; i < len(cards); i++ {
		if cards[i] != nil {
			cards[i].Update()
		}

	}

	return nil
}

var RewardSceneInstance = &RewardScene{}
var currentPick = 0

var CardStartX = 184.0
var CardStartY = 178.0
var CardDistX = 200.0
var CardPicStartX = CardStartX - 6
var CardPicStartY = CardStartY - 10

var cards []*core.MovableImage

func (r *RewardScene) Draw(screen *ebiten.Image) {
	bg := ebiten.NewImage(1024, 600)
	bg.Fill(color.RGBA{R: 0x21, G: 0x43, B: 0x58, A: 255})
	screen.DrawImage(bg, &ebiten.DrawImageOptions{})
	textTranslate := ebiten.GeoM{}
	textTranslate.Translate(512, 90)
	text.Draw(screen, "UPGRADE", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: textTranslate,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	})

	Geom := ebiten.GeoM{}
	// Geom.Translate(CardStartX, CardStartY)
	// screen.DrawImage(card1, &ebiten.DrawImageOptions{
	// 	GeoM: Geom,
	// })
	// Geom.Translate(CardDistX, 0)
	// screen.DrawImage(card2, &ebiten.DrawImageOptions{
	// 	GeoM: Geom,
	// })
	// Geom.Translate(CardDistX, 0)
	// screen.DrawImage(card3, &ebiten.DrawImageOptions{
	// 	GeoM: Geom,
	// })
	textBg := ebiten.NewImage(1200, 120)
	textBg.Fill(color.RGBA{R: 0, G: 0x97, B: 0xA4, A: 255})
	for i := 0; i < len(cards); i++ {
		if cards[i] != nil {
			cards[i].Draw(screen)
		}

	}
	//draw cursor
	Geom.Reset()
	Geom.Translate(CardPicStartX+float64(currentPick)*CardDistX, CardPicStartY)
	screen.DrawImage(assets.CardPick, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Reset()
	Geom.Translate(0, 435)
	screen.DrawImage(textBg, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Translate(10, 0)
	Desctext := casterPick[currentPick].GetDescription()
	text.Draw(screen, Desctext, assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: Geom,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignStart,
			LineSpacing:  15,
		},
	})

}

var casterPick = [3]ItemInterface{}

func GenerateCard(caster ItemInterface) *ebiten.Image {

	cardBounds := assets.CardTemplate.Bounds()
	newImage := ebiten.NewImage(cardBounds.Dx(), cardBounds.Dy())
	newImage.DrawImage(assets.CardTemplate, &ebiten.DrawImageOptions{})
	icon := caster.GetIcon()
	geom := ebiten.GeoM{}
	geom.Scale(2.1, 2.1)
	geom.Translate(50, 75)
	newImage.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
	geom.Reset()
	geom.Scale(1.2, 1.2)
	geom.Translate(13, 17)
	colorScale := &ebiten.ColorScale{}
	colorScale.Scale(1, 1, 1, 1)
	if cc, ok := caster.(loadout.Caster); ok {
		Cost := cc.GetCost() / 100
		text.Draw(newImage, fmt.Sprintf("%d", Cost), assets.FontFace, &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{
				GeoM:       geom,
				ColorScale: *colorScale,
			},
		})
	}

	geom.Reset()
	geom.Scale(0.8, 0.8)
	geom.Translate(float64(cardBounds.Dx())/2, 45)
	Name := caster.GetName()

	text.Draw(newImage, Name, assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM:       geom,
			ColorScale: *colorScale,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	})
	geom.Reset()
	geom.Scale(1.3, 1.3)
	geom.Translate(180, 190)
	if cc, ok := caster.(loadout.Caster); ok {
		Damage := cc.GetDamage()
		text.Draw(newImage, fmt.Sprintf("%d", Damage), assets.FontFace, &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{
				GeoM:       geom,
				ColorScale: *colorScale,
			},
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: text.AlignEnd,
			},
		})
	}

	geom.Reset()
	geom.Translate(10, 200)
	if cc, ok := caster.(loadout.Caster); ok {
		cooldown := cc.GetCooldownDuration()
		text.Draw(newImage, fmt.Sprintf("%.1fs", cooldown.Seconds()), assets.FontFace, &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{
				GeoM:       geom,
				ColorScale: *colorScale,
			},
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: text.AlignStart,
			},
		})
	}

	return newImage
}

func (r *RewardScene) Load(state *SceneData, manager stagehand.SceneController[*SceneData]) {
	r.sm = manager.(*stagehand.SceneDirector[*SceneData]) // This type assertion is important
	r.data = state
	if state.rewards == nil {
		casterPick[0] = GenerateReward()
		casterPick[1] = GenerateReward()
		casterPick[2] = GenerateReward()
	} else {
		casterPick[0] = state.rewards[0]
		casterPick[1] = state.rewards[1]
		casterPick[2] = state.rewards[2]
	}

	cards = []*core.MovableImage{nil, nil, nil}

	card1 := GenerateCard(casterPick[0])
	cards[0] = core.NewMovableImage(card1, core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: CardStartX, Sy: CardPicStartY}))
	card2 := GenerateCard(casterPick[1])
	cards[1] = core.NewMovableImage(card2, core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: CardStartX + CardDistX, Sy: CardPicStartY}))
	card3 := GenerateCard(casterPick[2])
	cards[2] = core.NewMovableImage(card3, core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: CardStartX + 2*CardDistX, Sy: CardPicStartY}))

	r.loopMusic = true
	var err error
	if r.musicPlayer == nil {
		r.musicPlayer, err = assets.NewAudioPlayer(assets.IntermissionMusic, assets.TypeMP3)
		if err != nil {
			fmt.Println(err.Error())
		}
		r.musicPlayer.AudioPlayer().Play()
	}
}
func (s *RewardScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *RewardScene) Unload() *SceneData {
	// your unload code
	s.data.CurrentLevel.SelectedStage = nil
	s.data.SceneDecor = nil
	s.loopMusic = false
	s.data.MusicSeek = s.musicPlayer.AudioPlayer().Position()
	s.musicPlayer.AudioPlayer().Rewind()
	s.musicPlayer.AudioPlayer().Pause()
	return s.data
}
