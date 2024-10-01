package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/joelschutz/stagehand"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/system/attack"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/kharism/hanashi/core"
)

type WorkshopScene struct {
	data           *SceneData
	sm             *stagehand.SceneDirector[*SceneData]
	selectedIndex  int
	cursorIsMoving bool
	moveLR         MoveCursorState
	curstate       string
}

var WorkshopSceneInstance = &WorkshopScene{}
var WORKSHOP_STATE_INVENTORY = "INV"
var WORKSHOP_STATE_UPGRADE = "UPGRADE"

func (r *WorkshopScene) Update() error {
	switch r.curstate {
	case WORKSHOP_STATE_INVENTORY:
		UpdateInventory(r)
	case WORKSHOP_STATE_UPGRADE:
		UpdateUpgrade(r)
	}

	return nil
}

func UpdateUpgrade(r *WorkshopScene) {
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		currentPick = (currentPick + 1) % 3
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if currentPick == 0 {
			currentPick = 2
		} else {
			currentPick = (currentPick - 1) % 3
		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		selectedCaster := decoratedCaster[currentPick]
		SetItem2(swapPayloadInstance.source, r.data, selectedCaster)
		swapPayloadInstance.source = nil
		r.sm.ProcessTrigger(TriggerToCombat)
	}
}
func SetItem2(pos *positionSwap, data *SceneData, i loadout.Caster) {
	if pos.itemCursorYPos == 0 {
		switch pos.XPos {
		case 0:
			data.MainLoadout[0] = i

		case 1:
			data.MainLoadout[1] = i

		case 2:
			data.SubLoadout1[0] = i

		case 3:
			data.SubLoadout1[1] = i

		case 4:
			data.SubLoadout2[0] = i

		case 5:
			data.SubLoadout2[1] = i

		}
	} else {
		caster := GetItem2(pos, data)
		for idx, cc := range data.Inventory {
			if cc == caster {
				data.Inventory[idx] = i
			}
		}
	}

}
func GetItem2(pos *positionSwap, data *SceneData) ItemInterface {
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
		return casterSlot[itemIdx]
	}
}

type CasterDecor func(loadout.Caster) loadout.Caster

var CasterDecorList = []CasterDecor{
	attack.DecorateWithHeal,
	attack.DecorateWithDoubleCast,
	attack.DecorateWithCooldownReduce,
	attack.DecorateWithBonus10,
	attack.DecorateWithCostReducer,
}
var decoratedCaster = []loadout.Caster{}

func UpdateInventory(r *WorkshopScene) {
	if !r.cursorIsMoving && inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if swapPayloadInstance.source == nil {
			swapPayloadInstance.source = &positionSwap{}
		}
		if itemCursorYPos == 0 {
			swapPayloadInstance.source.XPos = loadoutIdx
			swapPayloadInstance.source.itemCursorYPos = 0
		} else {
			swapPayloadInstance.source.XPos = itemIdx
			swapPayloadInstance.source.itemCursorYPos = 1
		}

		caster := GetItem2(swapPayloadInstance.source, r.data)
		cards = []*core.MovableImage{nil, nil, nil}
		decoratedCaster = []loadout.Caster{}
		// rand.Shuffle(len(CasterDecorList), func(i, j int) {
		// 	CasterDecorList[i], CasterDecorList[j] = CasterDecorList[j], CasterDecorList[i]
		// })
		for i := 0; i < 3; i++ {
			newCaster := CasterDecorList[i](caster.(loadout.Caster))
			decoratedCaster = append(decoratedCaster, newCaster)
			card1 := GenerateCard(newCaster)
			cards[i] = core.NewMovableImage(card1, core.NewMovableImageParams().WithMoveParam(core.MoveParam{Sx: CardStartX + CardDistX*float64(i), Sy: CardPicStartY}))
		}
		r.curstate = WORKSHOP_STATE_UPGRADE

	}
	if !r.cursorIsMoving && inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if itemCursorYPos == 1 {
			itemCursorYPos = 0
			targetX := ArrXposLoadout[loadoutIdx]
			targetY := 100.0
			scaleX := float64(32.0 / 195.0)
			scaleY := float64(32.0 / 250.0)
			scaleAnimation := core.ScaleAnimation{Tsx: scaleX, Tsy: scaleY, SpeedX: -0.02, SpeedY: -0.03}
			scaleAnimation.Apply(cardPickInventory)
			speed := 6.0
			if loadoutIdx == 4 || loadoutIdx == 5 {
				speed = 10
			}
			anim := core.NewMoveAnimationFromParam(core.MoveParam{
				Tx: targetX, Ty: targetY, Speed: speed,
			})
			r.cursorIsMoving = true
			cardPickInventory.Done = func() {
				r.cursorIsMoving = false
			}
			cardPickInventory.AddAnimation(anim)
			inventoryDesc = GetDescOfLoadout(r.data)
			r.moveLR = MoveCursorLeftRightLoadout
		}
	}
	if !r.cursorIsMoving && inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if itemCursorYPos == 0 && len(r.data.Inventory) > 0 {
			itemCursorYPos = 1
			targetX := 20.0 - 2
			targetY := CardStartY - 2
			scaleX := 1.0 //float64(32.0 / 195.0)
			scaleY := 1.0 //float64(32.0 / 250.0)
			scaleAnimation := core.ScaleAnimation{Tsx: scaleX, Tsy: scaleY, SpeedX: 0.02, SpeedY: 0.03}
			scaleAnimation.Apply(cardPickInventory)
			speed := 6.0
			if loadoutIdx == 4 || loadoutIdx == 5 {
				speed = 10
			}
			anim := core.NewMoveAnimationFromParam(core.MoveParam{
				Tx: targetX, Ty: targetY, Speed: speed,
			})
			r.cursorIsMoving = true
			cardPickInventory.Done = func() {
				r.cursorIsMoving = false
			}
			cardPickInventory.AddAnimation(anim)
			inventoryDesc = GetDescOfItem(r.data)
			r.moveLR = MoveCursorLeftRightInv
		}
	}

	if cardPickInventory == nil {
		if len(r.data.Inventory) > 0 && itemCursorYPos == 1 {
			cardPickInventory = core.NewMovableImage(assets.CardPick, core.NewMovableImageParams().WithMoveParam(core.MoveParam{
				Sx: 20 - 5,
				Sy: CardStartY - 10,
			}))
			itemCursorYPos = 1
		} else {
			cardPickInventory = core.NewMovableImage(assets.CardPick, core.NewMovableImageParams().WithMoveParam(core.MoveParam{
				Sx: 140 - 2,
				Sy: 100 - 2,
			}).WithScale(&core.ScaleParam{Sx: float64(34.0 / 195), Sy: float64(34.0 / 250.0)}))
			itemCursorYPos = 0
		}

	}
	if !r.cursorIsMoving && inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		r.moveLR(r.data, +1)
	}
	if !r.cursorIsMoving && inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		r.moveLR(r.data, -1)
	}
	cardPickInventory.Update()
	for _, j := range ItemSlot {
		j.Update()
	}
}
func DrawInventoryState(screen *ebiten.Image, data *SceneData) {
	screen.DrawImage(assets.BgWorkbench, &ebiten.DrawImageOptions{})

	textTranslate := ebiten.GeoM{}
	textTranslate.Translate(512, 50)
	text.Draw(screen, "Upgrade", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: textTranslate,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	})
	textTranslate.Translate(500, -20)
	text.Draw(screen, "Press q to upgrade", assets.FontFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: textTranslate,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignEnd,
		},
	})
	textBg := ebiten.NewImage(1200, 120)
	textBg.Fill(color.RGBA{R: 0, G: 0x97, B: 0xA4, A: 255})
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
	if data.MainLoadout[0] != nil {
		icon = data.MainLoadout[0].GetIcon()
	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Translate(40, 0)
	if data.MainLoadout[1] != nil {
		icon = data.MainLoadout[1].GetIcon()
	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Reset()
	Geom.Translate(380, 100)
	if data.SubLoadout1[0] != nil {
		icon = data.SubLoadout1[0].GetIcon()

	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Translate(40, 0)
	if data.SubLoadout1[1] != nil {
		icon = data.SubLoadout1[1].GetIcon()
	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Reset()
	Geom.Translate(380+240, 100)
	if data.SubLoadout2[0] != nil {
		icon = data.SubLoadout2[0].GetIcon()
	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})
	Geom.Translate(40, 0)
	if data.SubLoadout2[1] != nil {
		icon = data.SubLoadout2[1].GetIcon()
	} else {
		icon = assets.NAIcon
	}
	screen.DrawImage(icon, &ebiten.DrawImageOptions{
		GeoM: Geom,
	})

	if cardPickInventory != nil {
		cardPickInventory.Draw(screen)
	}
	// ItemSlot = []*ebiten.Image{}
	for _, j := range ItemSlot {
		j.Draw(screen)
	}
	Geom.Reset()
	Geom.Translate(10, 435)
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
func DrawUpgradeState(screen *ebiten.Image, data *SceneData) {
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
	Desctext := decoratedCaster[currentPick].GetDescription()
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

func (r *WorkshopScene) Draw(screen *ebiten.Image) {
	switch r.curstate {
	case WORKSHOP_STATE_INVENTORY:
		DrawInventoryState(screen, r.data)
	case WORKSHOP_STATE_UPGRADE:
		DrawUpgradeState(screen, r.data)
	}

}

var casterSlot []loadout.Caster

func (r *WorkshopScene) Load(state *SceneData, manager stagehand.SceneController[*SceneData]) {
	r.sm = manager.(*stagehand.SceneDirector[*SceneData]) // This type assertion is important
	r.data = state
	cardPickInventory = nil
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
	r.curstate = WORKSHOP_STATE_INVENTORY
	ItemSlot = []*core.MovableImage{}
	casterSlot = []loadout.Caster{}
	for inventoryIdx, j := range r.data.Inventory {
		if vv, ok := j.(loadout.Caster); ok {
			c := GenerateCard(vv)
			dim := c.Bounds()
			newMvImage := core.NewMovableImage(c, core.NewMovableImageParams().WithMoveParam(core.MoveParam{
				Sx: 23 + float64(inventoryIdx*(dim.Dx()+30)),
				Sy: CardStartY + 2,
			}))
			casterSlot = append(casterSlot, vv)
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
}

func (s *WorkshopScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 600
}
func (s *WorkshopScene) Unload() *SceneData {
	return s.data
}
