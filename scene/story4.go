package scene

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/hanashi/core"
)

func Scene4(layouter core.GetLayouter) *core.Scene {
	scene := core.NewScene()
	scene.SetLayouter(layouter)
	shizukuOverlay := ebiten.NewImage(64, 64)
	shizukuOverlay.DrawImage(assets.Shizuku, &ebiten.DrawImageOptions{})
	shizukuOverlay.DrawImage(assets.Overlay, &ebiten.DrawImageOptions{})

	JackOverlay := ebiten.NewImage(64, 64)
	JackOverlay.DrawImage(assets.Jack, &ebiten.DrawImageOptions{})
	JackOverlay.DrawImage(assets.Overlay, &ebiten.DrawImageOptions{})
	portraitMoveParam := core.MoveParam{Sx: 10, Sy: 450, Tx: 10, Ty: 450}
	portraitScaleParam := &core.ScaleParam{Sx: 2, Sy: 2}
	// sceneWidth, sceneHeight := layouter.GetLayout()
	scene.Characters = []*core.Character{
		core.NewCharacterImage("Sven", assets.Sven2),
		core.NewCharacterImage("Shizuku", shizukuOverlay),
		core.NewCharacterImage("Jack", JackOverlay),
		core.NewCharacterImage("Landon", assets.Landon),
		core.NewCharacterImage("<unknown spirit beast>", assets.Lupus),
	}
	scene.Events = []core.Event{
		&core.ComplexEvent{Events: []core.Event{
			&core.PlayBgmEvent{Audio: &assets.Unease1, Type: core.TypeMP3},
			core.NewBgChangeEvent(assets.Cave, core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: 0, Speed: 3}, nil),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "*pant, pant*(Finally!!!)", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "*kneel on the ground*(Just in time the anomalous spirit flow\nbecome unbearable)", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Sven", Dialogue: "(I'll check it later, for now...)*gets up*", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Sven", Dialogue: "(time to get the meat)*walks to the corpse*", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Sven", Dialogue: "(wait, what?)", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Sven", Dialogue: "*smaller baby wyrm rips out off the stomach of the dead wyrm*", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Sven", Dialogue: "(So, it was pregnant as well huh?)", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Sven", Dialogue: "(Sorry, but it has to be this way)", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Oh you managed to kill it?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			// &core.PlayBgmEvent{Audio: &assets.MidMusic, Type: core.TypeMP3},
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Yea, now how should I cut the meat ?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Careful with the cuts. The skin are still usable\nand you might hit the dragon stone inside", FontFace: assets.FontFace},
		}},

		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			&core.StopBgmEvent{},

			&core.PlaySfxEvent{Audio: &assets.MagibulletFx3x, Type: core.TypeMP3}, // play shoot audio here
			&core.DialogueEvent{Name: "", Dialogue: "", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			// &core.PlayBgmEvent{Audio: &assets.MidMusic, Type: core.TypeMP3},
			&core.PlayBgmEvent{Audio: &assets.Unease2, Type: core.TypeMP3},
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "What the....?", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "Landon...", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Landon", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Landon", Dialogue: "Been a while sven", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Landon", Dialogue: "Miss me?", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Landon"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Considering your aiming,\nThat should be my question", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Landon", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Landon", Dialogue: "Not changing even a bit, I see", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Landon", Dialogue: "Now, Would you kindly surrender your suit to me?", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Landon", Dialogue: "Or don't", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Landon"),
			&core.PlaySfxEvent{Audio: &assets.MagibulletFx3x, Type: core.TypeMP3}, // play shoot audio here
			&core.DialogueEvent{Name: "", Dialogue: ""},
		}},

		&core.ComplexEvent{Events: []core.Event{
			&core.PlaySfxEvent{Audio: &assets.GrowlFx, Type: core.TypeMP3}, // play shoot audio here
			core.NewCharacterAddEvent("<unknown spirit beast>", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "", Dialogue: "*A spirit beast appear out of the armor and parry the bullets*", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("<unknown spirit beast>"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "What the??"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "It appears the partially absorbed spirit has become its own things"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "And it coming out of the suit makes the spirit flow in the suit\nnormal again"},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "Time to get serious!!!"},
	}

	scene.TxtBg = ebiten.NewImage(1024-128, 128)
	scene.TxtBg.Fill(color.RGBA{R: 0x4f, G: 0x8f, B: 0xba, A: 255})
	pp, err := core.NewDefaultAudioInterfacer()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	scene.AudioInterface = pp
	return scene
}
