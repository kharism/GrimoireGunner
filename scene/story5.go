package scene

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/hanashi/core"
)

func Scene5(layouter core.GetLayouter) *core.Scene {
	scene := core.NewScene()
	scene.SetLayouter(layouter)
	portraitMoveParam := core.MoveParam{Sx: 10, Sy: 450, Tx: 10, Ty: 450}
	portraitScaleParam := &core.ScaleParam{Sx: 2, Sy: 2}

	JackOverlay := ebiten.NewImage(64, 64)
	JackOverlay.DrawImage(assets.Jack, &ebiten.DrawImageOptions{})
	JackOverlay.DrawImage(assets.Overlay, &ebiten.DrawImageOptions{})
	scene.FontFace = assets.FontFace
	scene.Characters = []*core.Character{
		core.NewCharacterImage("Sven", assets.Sven2),
		core.NewCharacterImage("Shizuku", assets.Shizuku),
		core.NewCharacterImage("Jack", JackOverlay),
		core.NewCharacterImage("Landon", assets.Landon),
		core.NewCharacterImage("<unknown spirit beast>", assets.Lupus),
		core.NewCharacterImage("Lupus", assets.Lupus),
	}
	blackBg := ebiten.NewImage(1200, 600)
	blackBg.Fill(color.Black)
	scene.Events = []core.Event{
		&core.ComplexEvent{Events: []core.Event{
			&core.PlayBgmEvent{Audio: &assets.Unease1, Type: core.TypeMP3},
			core.NewBgChangeEvent(assets.Cave, core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: 0, Speed: 3}, nil),
			&core.DialogueEvent{Name: "", Dialogue: "*Sven managed to overpower Landon*", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{

			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Surrender now Landon!!", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "Or else...*Points gun at landon*", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Landon", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Landon", Dialogue: "You won this time sven", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Landon"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Scram!!!", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "Or else", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			&core.PlaySfxEvent{Audio: &assets.MagibulletFx3x, Type: core.TypeMP3}, // play shoot audio here
			&core.StopBgmEvent{},
			&core.DialogueEvent{Name: "", Dialogue: ""},
		}},
		&core.DialogueEvent{Name: "", Dialogue: "*Landon escape from sven*", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			// core.NewCharacterRemoveEvent("Landon"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Now what to do with you little one?*look at the spirit beast*", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("<unknown spirit beast>", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "<unknown spirit beast>", Dialogue: "*looking at sven with puppy eyes*", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("<unknown spirit beast>"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Firstly give it a name,\nat least for now", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "so we have some easier way to refer it", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "It looks like a wolf, so lupus it is", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "It seems to enjoy the name isn't it?", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			&core.PlaySfxEvent{Audio: &assets.GrowlFx, Type: core.TypeMP3},
			core.NewCharacterAddEvent("Lupus", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Lupus", Dialogue: "*bark*", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Lupus"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "We'll talk about it once you get home", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "For now just get the carcass and go home safely", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewBgChangeEvent(blackBg, core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: 0, Speed: 3}, nil),
			core.NewCharacterRemoveEvent("Jack"),
			&core.CharacterImageSwapEvent{Name: "Sven", NewImage: assets.Sven},
			&core.DialogueEvent{Name: "", Dialogue: "*After that, back at home*", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.PlayBgmEvent{Audio: &assets.Downtimemusic, Type: core.TypeMP3},
			core.NewBgChangeEvent(assets.Kitchen, core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: 0, Speed: 3}, nil),
			&core.PlaySfxEvent{Audio: &assets.GrowlFx, Type: core.TypeMP3},
			&core.DialogueEvent{Name: "Lupus", Dialogue: "*Bark*", FontFace: assets.FontFace},
			&core.CharacterImageSwapEvent{Name: "Jack", NewImage: assets.Jack},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Lupus"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Good boy!!!!", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "I already like him", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Jack", Dialogue: "You know, Your mother and I had a plan to adopt\na spirit beast once you married and move out. But...", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "....", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "Father!!!!", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Sorry, My bad, My bad", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "Anyway", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Jack", Dialogue: "I've talked with inspector Townsend about Lupus", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Jack", Dialogue: "He said to bring Lupus to the Hunter office tommorow for inspection", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Jack", Dialogue: "If there are no issues, we can adopt lupus", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Jack", Dialogue: "Let's have dinner shall we?", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "Hope you like today's dinner", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Shizuku", Dialogue: "It's wyrm meat spaghetti", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "Hope you like today's dinner", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "*shows nauseous pale face*", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "*(it reminds me of the baby wyrm wrigling out of their mother womb)", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "What's wrong?", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Shizuku", Dialogue: "Now, you're the one looking pale", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Nothing. I'm just....", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "not hungry", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "The spaghetti probably just reminded him off baby wyrm\nbursting out their mother stomach", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "*turns pale too* Eeewww", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Shizuku", Dialogue: "*Hold to the chair and her head*", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Welp, Should I make fried rice for us?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "*Sit on the chair*Yes, please do", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Well more for me then", FontFace: assets.FontFace},
		}},
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
