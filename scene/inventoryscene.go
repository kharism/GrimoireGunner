package scene

import (
	"errors"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"
)

type InventoryScene struct {
	data           *SceneData
	sm             *stagehand.SceneDirector[*SceneData]
	cursorIsMoving bool
	moveLR         MoveCursorState
	musicPlayer    *assets.AudioPlayer
	loopMusic      bool
}

var itemCursorYPos int // fill this with 0 or 1. 0 means the loadout, 1 means the items

var loadoutIdx int //index for loadout

var itemIdx int //index for item
var ItemSlot []*core.MovableImage
var inventoryDesc string
var inventoryName string

type positionSwap struct {
	itemCursorYPos int
	XPos           int
}
type swapPayload struct {
	source *positionSwap
	dest   *positionSwap
}

func GetItem(pos *positionSwap, data *SceneData) ItemInterface {
	if pos.itemCursorYPos == 0 {
		switch pos.XPos {
		case 0:
			return data.MainLoadout[0]
		case 1:
			return data.MainLoadout[1]
		case 2:
			return data.SubLoadout1[0]
		case 3:
			return data.SubLoadout1[1]
		case 4:
			return data.SubLoadout2[0]
		case 5:
			return data.SubLoadout2[1]
		}
		return nil
	} else {
		return data.Inventory[itemIdx]
	}
}
func SetItem(pos *positionSwap, r *InventoryScene, i ItemInterface) error {
	if pos.itemCursorYPos == 0 {
		ii, ok := i.(loadout.Caster)
		if i != nil && !ok {
			return errors.New("incompatible")
		}
		switch pos.XPos {
		case 0:
			r.data.MainLoadout[0] = ii
			return nil
		case 1:
			r.data.MainLoadout[1] = ii
			return nil
		case 2:
			r.data.SubLoadout1[0] = ii
			return nil
		case 3:
			r.data.SubLoadout1[1] = ii
			return nil
		case 4:
			r.data.SubLoadout2[0] = ii
			return nil
		case 5:
			r.data.SubLoadout2[1] = ii
			return nil
		}
		return nil
	} else {
		r.data.Inventory[itemIdx] = i
		return nil
	}
}

var swapPayloadInstance = &swapPayload{}

func (s *swapPayload) Swap(r *InventoryScene) {
	source1 := GetItem(s.source, r.data)
	source2 := GetItem(s.dest, r.data)
	if SetItem(s.source, r, source2) == nil {
		if SetItem(s.dest, r, source1) != nil {
			//rollback
			SetItem(s.source, r, source1)
		} else {
			newInv := []ItemInterface{}
			for _, v := range r.data.Inventory {
				if v != nil {
					newInv = append(newInv, v)
				}
			}

			r.data.Inventory = newInv
			// update invlist
			ItemSlot = []*core.MovableImage{}
			for inventoryIdx, j := range r.data.Inventory {
				c := GenerateCard(j)
				dim := c.Bounds()
				newMvImage := core.NewMovableImage(c, core.NewMovableImageParams().WithMoveParam(core.MoveParam{
					Sx: 23 + float64(inventoryIdx*(dim.Dx()+30)),
					Sy: CardStartY + 2,
				}))
				ItemSlot = append(ItemSlot, newMvImage)
				itemIdx = 0
				// GenerateCard()
			}
		}
	}

}

func MoveCursorLeftRightLoadout(data *SceneData, LR int) {
	_, targetY := cardPickInventory.GetPos()
	loadoutIdx += LR
	if loadoutIdx == -1 {
		loadoutIdx = 0
		return
	}
	if loadoutIdx == 6 {
		loadoutIdx = 5
		return
	}
	inventoryDesc = GetDescOfLoadout(data)
	inventoryName = GetNameOfLoadout(data)
	targetX := float64(ArrXposLoadout[loadoutIdx]) - 4
	cardPickInventory.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{
		Tx: targetX, Ty: targetY, Speed: 6,
	}))

}
func MoveCursorLeftRightInv(data *SceneData, LR int) {
	itemIdx += LR
	if itemIdx == -1 {
		itemIdx = 0
		return
	}
	if itemIdx >= len(ItemSlot) {
		itemIdx = len(ItemSlot) - 1
		return
	}
	inventoryDesc = GetDescOfItem(data)
	inventoryName = GetNameOfItem(data)
	for _, j := range ItemSlot {
		posX, Ty := j.GetPos()
		width, _ := j.GetSize()
		Tx := posX + float64(-LR)*(width+30)
		j.AddAnimation(core.NewMoveAnimationFromParam(core.MoveParam{Tx: Tx, Ty: Ty, Speed: 10}))
	}
}

// abstract function to handle the left/right movement of the cursor
// LR is integer to tell whether the user move cursor left or right
// 1 is for right, -1 is for left
type MoveCursorState func(r *SceneData, LR int)

func GetDescOfLoadout(data *SceneData) string {
	switch loadoutIdx {
	case 0:
		if data.MainLoadout[0] != nil {
			return data.MainLoadout[0].GetDescription()
		}
	case 1:
		if data.MainLoadout[1] != nil {
			return data.MainLoadout[1].GetDescription()
		}
	case 2:
		if data.SubLoadout1[0] != nil {
			return data.SubLoadout1[0].GetDescription()
		}
	case 3:
		if data.SubLoadout1[1] != nil {
			return data.SubLoadout1[1].GetDescription()
		}
	case 4:
		if data.SubLoadout2[0] != nil {
			return data.SubLoadout2[0].GetDescription()
		}
	case 5:
		if data.SubLoadout2[1] != nil {
			return data.SubLoadout2[1].GetDescription()
		}
	}
	return ""

}
func GetNameOfLoadout(data *SceneData) string {
	switch loadoutIdx {
	case 0:
		if data.MainLoadout[0] != nil {
			return data.MainLoadout[0].GetName()
		}
	case 1:
		if data.MainLoadout[1] != nil {
			return data.MainLoadout[1].GetName()
		}
	case 2:
		if data.SubLoadout1[0] != nil {
			return data.SubLoadout1[0].GetName()
		}
	case 3:
		if data.SubLoadout1[1] != nil {
			return data.SubLoadout1[1].GetName()
		}
	case 4:
		if data.SubLoadout2[0] != nil {
			return data.SubLoadout2[0].GetName()
		}
	case 5:
		if data.SubLoadout2[1] != nil {
			return data.SubLoadout2[1].GetName()
		}
	}
	return ""

}
func GetDescOfItem(data *SceneData) string {
	return data.Inventory[itemIdx].GetDescription()
}
func GetNameOfItem(data *SceneData) string {
	return data.Inventory[itemIdx].GetName()
}

func (r *InventoryScene) Update() error {
	if r.loopMusic && !r.musicPlayer.AudioPlayer().IsPlaying() {
		r.musicPlayer.AudioPlayer().Rewind()
		r.musicPlayer.AudioPlayer().Play()
	}
	if r.musicPlayer != nil {
		r.musicPlayer.Update()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		r.sm.ProcessTrigger(TriggerToCombat)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		r.musicPlayer.QueueSFX(assets.CursorFx)
		if swapPayloadInstance.source == nil {
			swapPayloadInstance.source = &positionSwap{itemCursorYPos: itemCursorYPos}
			if itemCursorYPos == 0 {
				swapPayloadInstance.source.XPos = loadoutIdx
			} else {
				swapPayloadInstance.source.XPos = itemIdx
			}
		} else {
			swapPayloadInstance.dest = &positionSwap{itemCursorYPos: itemCursorYPos}
			if itemCursorYPos == 0 {
				swapPayloadInstance.dest.XPos = loadoutIdx
			} else {
				swapPayloadInstance.dest.XPos = itemIdx
			}
			swapPayloadInstance.Swap(r)
			swapPayloadInstance.source = nil
		}
	}
	if cardPickInventory == nil {
		if len(r.data.Inventory) > 0 && itemCursorYPos == 1 {
			cardPickInventory = core.NewMovableImage(assets.CardPickRed, core.NewMovableImageParams().WithMoveParam(core.MoveParam{
				Sx: 20 - 5,
				Sy: CardStartY - 10,
			}))
			itemCursorYPos = 1
		} else {
			cardPickInventory = core.NewMovableImage(assets.CardPickRed, core.NewMovableImageParams().WithMoveParam(core.MoveParam{
				Sx: 140 - 4,
				Sy: 100 - 4,
			}).WithScale(&core.ScaleParam{Sx: float64(40.0 / 195), Sy: float64(40.0 / 250.0)}))
			itemCursorYPos = 0
		}

	}
	if !r.cursorIsMoving && inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		r.musicPlayer.QueueSFX(assets.CursorFx)
		if itemCursorYPos == 1 {
			itemCursorYPos = 0
			targetX := ArrXposLoadout[loadoutIdx] - 4
			targetY := 100.0 - 4
			scaleX := float64(40.0 / 195.0)
			scaleY := float64(40.0 / 250.0)
			time := 5.0
			posX, posY := cardPickInventory.GetPos()
			xSpeed := (targetX - posX) / time
			ySpeed := (targetY - posY) / time
			varSpeed := math.Sqrt(xSpeed*xSpeed + ySpeed*ySpeed)
			scaleAnimation := core.ScaleAnimation{Tsx: scaleX, Tsy: scaleY, SpeedX: -0.02, SpeedY: -0.03}
			scaleAnimation.Apply(cardPickInventory)
			anim := core.NewMoveAnimationFromParam(core.MoveParam{
				Tx: targetX, Ty: targetY, Sx: xSpeed, Sy: ySpeed, Speed: varSpeed,
			})
			r.cursorIsMoving = true
			cardPickInventory.Done = func() {
				r.cursorIsMoving = false
			}
			cardPickInventory.AddAnimation(anim)
			inventoryDesc = GetDescOfLoadout(r.data)
			inventoryName = GetNameOfLoadout(r.data)
			r.moveLR = MoveCursorLeftRightLoadout
		}
	}
	if !r.cursorIsMoving && inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		r.musicPlayer.QueueSFX(assets.CursorFx)
		if itemCursorYPos == 0 && len(r.data.Inventory) > 0 {
			itemCursorYPos = 1
			targetX := 20.0 - 2
			targetY := CardStartY - 2
			scaleX := 1.0 //float64(32.0 / 195.0)
			scaleY := 1.0 //float64(32.0 / 250.0)
			scaleAnimation := core.ScaleAnimation{Tsx: scaleX, Tsy: scaleY, SpeedX: 0.02, SpeedY: 0.03}
			scaleAnimation.Apply(cardPickInventory)
			time := 5.0
			posX, posY := cardPickInventory.GetPos()
			xSpeed := (targetX - posX) / time
			ySpeed := (targetY - posY) / time
			varSpeed := math.Sqrt(xSpeed*xSpeed + ySpeed*ySpeed)
			anim := core.NewMoveAnimationFromParam(core.MoveParam{
				Tx: targetX, Ty: targetY, Sx: xSpeed, Sy: ySpeed, Speed: varSpeed,
			})
			r.cursorIsMoving = true
			cardPickInventory.Done = func() {
				r.cursorIsMoving = false
			}
			cardPickInventory.AddAnimation(anim)
			inventoryDesc = GetDescOfItem(r.data)
			inventoryName = GetNameOfItem(r.data)
			r.moveLR = MoveCursorLeftRightInv
		}
	}
	if !r.cursorIsMoving && itemCursorYPos == 0 && inpututil.IsKeyJustPressed(ebiten.KeyE) {
		r.musicPlayer.QueueSFX(assets.CursorFx)
		item := GetItem(&positionSwap{0, loadoutIdx}, r.data)
		SetItem(&positionSwap{0, loadoutIdx}, r, nil)
		r.data.Inventory = append(r.data.Inventory, item)
		ItemSlot = []*core.MovableImage{}
		for inventoryIdx, j := range r.data.Inventory {
			if vv, ok := j.(loadout.Caster); ok {
				c := GenerateCard(vv)
				dim := c.Bounds()
				newMvImage := core.NewMovableImage(c, core.NewMovableImageParams().WithMoveParam(core.MoveParam{
					Sx: 23 + float64(inventoryIdx*(dim.Dx()+30)),
					Sy: CardStartY + 2,
				}))
				ItemSlot = append(ItemSlot, newMvImage)
				itemIdx = 0
			}
			// GenerateCard()
		}

	}
	if !r.cursorIsMoving && inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		r.musicPlayer.QueueSFX(assets.CursorFx)
		r.moveLR(r.data, +1)
	}
	if !r.cursorIsMoving && inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		r.musicPlayer.QueueSFX(assets.CursorFx)
		r.moveLR(r.data, -1)
	}
	cardPickInventory.Update()
	for _, j := range ItemSlot {
		j.Update()
	}
	return nil
}

var cardPickInventory *core.MovableImage
var LoadoutInvStartX = 140.0
var LoadoutInvStartY = 100.0

var ArrXposLoadout = []float64{140, 140 + 40, 380, 380 + 40, 380 + 240, 380 + 240 + 40}

func (r *InventoryScene) Draw(screen *ebiten.Image) {
	bg := ebiten.NewImage(1024, 600)
	bg.Fill(color.RGBA{R: 0x21, G: 0x43, B: 0x58, A: 255})
	screen.DrawImage(bg, &ebiten.DrawImageOptions{})
	textTranslate := ebiten.GeoM{}
	textTranslate.Translate(512, 50)
	text.Draw(screen, "Inventory", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: textTranslate,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	})
	textTranslate.Translate(500, -20)
	text.Draw(screen, "Press q to swap", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: textTranslate,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignEnd,
		},
	})
	textTranslate.Translate(0, 20)
	text.Draw(screen, "Press w to return", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: textTranslate,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignEnd,
		},
	})
	textTranslate.Translate(0, 20)
	if itemCursorYPos == 0 {
		text.Draw(screen, "Press e to unequip", assets.FontFace, &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{
				GeoM: textTranslate,
			},
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: text.AlignEnd,
			},
		})
	}

	textBg := ebiten.NewImage(1200, 120)
	textBg.Fill(color.RGBA{R: 0, G: 0x97, B: 0xA4, A: 255})

	if itemCursorYPos == 0 {

	}
	Geom := ebiten.GeoM{}
	Geom.Reset()
	Geom.Translate(0, 435)
	screen.DrawImage(textBg, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})

	Geom.Reset()
	Geom.Translate(140, 80)
	text.Draw(screen, "MainLoadout", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: Geom,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignStart,
		},
	})

	// Geom.Reset()
	Geom.Translate(240, 0)
	text.Draw(screen, "SubLoadout1", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: Geom,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignStart,
		},
	})
	Geom.Translate(240, 0)
	text.Draw(screen, "SubLoadout2", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: Geom,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignStart,
		},
	})
	Geom.Reset()
	Geom.Translate(LoadoutInvStartX, LoadoutInvStartY)
	var icon *ebiten.Image
	if r.data.MainLoadout[0] != nil {
		icon = r.data.MainLoadout[0].GetIcon()
	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Translate(40, 0)
	if r.data.MainLoadout[1] != nil {
		icon = r.data.MainLoadout[1].GetIcon()
	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Reset()
	Geom.Translate(380, 100)
	if r.data.SubLoadout1[0] != nil {
		icon = r.data.SubLoadout1[0].GetIcon()

	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Translate(40, 0)
	if r.data.SubLoadout1[1] != nil {
		icon = r.data.SubLoadout1[1].GetIcon()
	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Reset()
	Geom.Translate(380+240, 100)
	if r.data.SubLoadout2[0] != nil {
		icon = r.data.SubLoadout2[0].GetIcon()
	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Translate(40, 0)
	if r.data.SubLoadout2[1] != nil {
		icon = r.data.SubLoadout2[1].GetIcon()
	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})

	if cardPickInventory != nil {
		cardPickInventory.Draw(screen)
	}
	if swapPayloadInstance.source != nil {
		if swapPayloadInstance.source.itemCursorYPos == 0 {
			// draw white box on loadout
			geom := ebiten.GeoM{}
			// posX, posY := cardPickInventory.GetPos()
			targetX := ArrXposLoadout[swapPayloadInstance.source.XPos] - 4
			targetY := 100.0 - 4
			scaleX := float64(40.0 / 195.0)
			scaleY := float64(40.0 / 250.0)
			geom.Scale(scaleX, scaleY)
			geom.Translate(targetX, targetY)
			opt := &ebiten.DrawImageOptions{
				GeoM: geom,
			}
			screen.DrawImage(assets.CardPick, opt)
		} else {
			// draw white box on inventory
			geom := ebiten.GeoM{}
			// posX, posY := cardPickInventory.GetPos()
			targetX := 20.0 - 2
			targetY := CardStartY - 2
			scaleX := 1.0
			scaleY := 1.0
			geom.Scale(scaleX, scaleY)
			geom.Translate(targetX, targetY)
			opt := &ebiten.DrawImageOptions{
				GeoM: geom,
			}
			screen.DrawImage(assets.CardPick, opt)
		}
	}
	// ItemSlot = []*ebiten.Image{}
	for _, j := range ItemSlot {
		j.Draw(screen)
	}
	Geom.Reset()
	Geom.Translate(10, 435)
	text.Draw(screen, inventoryName, assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: Geom,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignStart,
			LineSpacing:  15,
		},
	})
	Geom.Translate(0, 15)
	text.Draw(screen, inventoryDesc, assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: Geom,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignStart,
			LineSpacing:  15,
		},
	})
}

// interface for items
// Just a subset of caster
type ItemInterface interface {
	GetIcon() *ebiten.Image
	GetDescription() string
	GetName() string
}

type OnAquireDoer interface {
	OnAquireDo(*SceneData)
}

var InventorySceneInstance = &InventoryScene{}

func (r *InventoryScene) Load(state *SceneData, manager stagehand.SceneController[*SceneData]) {
	r.sm = manager.(*stagehand.SceneDirector[*SceneData]) // This type assertion is important
	r.data = state
	cardPickInventory = nil
	loadoutIdx = 0
	itemIdx = 0
	swapPayloadInstance.source = nil
	if len(state.Inventory) == 0 {
		itemCursorYPos = 0
		r.moveLR = MoveCursorLeftRightLoadout
		inventoryDesc = GetDescOfLoadout(r.data)

	} else {
		// itemCursorYPos = 1
		if itemCursorYPos == 1 {
			r.moveLR = MoveCursorLeftRightInv
		} else {
			r.moveLR = MoveCursorLeftRightLoadout
		}

		inventoryDesc = GetDescOfItem(r.data)
	}
	ItemSlot = []*core.MovableImage{}
	for inventoryIdx, j := range r.data.Inventory {
		if vv, ok := j.(ItemInterface); ok {
			c := GenerateCard(vv)
			dim := c.Bounds()
			newMvImage := core.NewMovableImage(c, core.NewMovableImageParams().WithMoveParam(core.MoveParam{
				Sx: 23 + float64(inventoryIdx*(dim.Dx()+30)),
				Sy: CardStartY + 2,
			}))
			ItemSlot = append(ItemSlot, newMvImage)
			// dim := c.Bounds()
			// Geom := ebiten.GeoM{}
			// Geom.Translate(20+float64(inventoryIdx*dim.Dx()), CardStartY)
			// screen.DrawImage(c, &ebiten.DrawImageOptions{
			// 	GeoM: Geom,
			// })
		}
		// GenerateCard()
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
	// fmt.Println("Load Inventory")
}
func (s *InventoryScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *InventoryScene) Unload() *SceneData {
	s.loopMusic = false
	s.data.MusicSeek = s.musicPlayer.AudioPlayer().Position()
	s.musicPlayer.AudioPlayer().Rewind()
	s.musicPlayer.AudioPlayer().Pause()
	return s.data
}
