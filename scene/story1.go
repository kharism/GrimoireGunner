package scene

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/hanashi/core"
)

// intro. Set up the characters and their motivation
func Scene1(layouter core.GetLayouter) *core.Scene {
	scene := core.NewScene()
	scene.SetLayouter(layouter)

	scene.Characters = []*core.Character{
		core.NewCharacterImage("Sven", assets.Sven),
		core.NewCharacterImage("Shizuku", assets.Shizuku),
		core.NewCharacterImage("Jack", assets.Jack),
	}
	scene.FontFace = assets.FontFace
	fmt.Println(scene.FontFace)
	portraitMoveParam := core.MoveParam{Sx: 10, Sy: 450, Tx: 10, Ty: 450}
	portraitScaleParam := &core.ScaleParam{Sx: 2, Sy: 2}
	sceneWidth, sceneHeight := layouter.GetLayout()
	blackBg := ebiten.NewImage(sceneWidth, sceneHeight)
	scene.Events = []core.Event{
		&core.ComplexEvent{Events: []core.Event{
			&core.PlayBgmEvent{Audio: &assets.SmoothJazz, Type: core.TypeMP3},
			core.NewBgChangeEvent(assets.Bedroom, core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: 0, Speed: 3}, nil),
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			core.NewDialogueEvent("Sven", "Rest here. Don't move around too much.", assets.FontFace),
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			core.NewDialogueEvent("Shizuku", "....", assets.FontFace),
		}},
		&core.DialogueEvent{Name: "Shizuku", Dialogue: "*Pick a paper bag and vomit into it*", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			core.NewDialogueEvent("Sven", "*pat her back* there, there let it out", assets.FontFace),
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "....", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Are you hungry?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "I'm not hungry", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "You need to eat regularly for our babies", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "*irritated* I know. I just......", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Don't worry about it. I'll get you some food okay?\nJust rest for now.", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewBgChangeEvent(assets.BedroomDoor, core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: 0, Speed: 3}, nil),
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "How its going? the gynecologyst said something?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Worse than we expect.", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.DialogueEvent{Name: "Sven", Dialogue: "We're expecting twins and each\n baby have different elemental affinity than their mother's,", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.DialogueEvent{Name: "Sven", Dialogue: "significantly increasing risk of miscarriage. ", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "That's the risk of with female yokai pregnancy.\nIt is small risk but a risk nonetheless.", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "Your mother also had similar complication when having you.\nI had to get some meat of jade wyrm to make a medicine for her", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Isn't that illegal?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Not back then.....well not today either.", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "Jade wyrm hunting is strictly regulated but not outright banned.", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Really?? I've never knew about the regulation allows it before", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "It is written in lawyerspeak.\nTo put it simply a B-rank or lower Stalker-unit\nlike you get the chance to hunt a mystical animal\nonce every 5 years", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "The catch is, you'll need to start with basic equipment", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Jack", Dialogue: "I'll prepare the paperwork while you prepare my old suit from\nstorage unit", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Thanks dad!!!!", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewBgChangeEvent(blackBg, core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: 0, Speed: 3}, nil),
			core.NewCharacterRemoveEvent("Sven"),
			&core.StopBgmEvent{},
			&core.PlayBgmEvent{Audio: &assets.Sax, Type: core.TypeMP3},
			&core.DialogueEvent{Name: "", Dialogue: "Next Morning", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewBgChangeEvent(assets.Workshop1, core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: 0, Speed: 3}, nil),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Done with the suit?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Yes. The self diagnostic does not show any anomaly\nNeither does the external diagnostics", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Ah The good ol' reliable Falken model.\nAs expected from well preserved military model", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "Brings back memory when I hunt jade wyrm for your mother", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "So, you're going??", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Yes", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "Don't worry. It won't be long.\nfather will accompany you while I'm away", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: ".....", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "*Hugs shizuku*", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "I'll be back with the cure okay.", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "*hugs him tigher*", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Shizuku", Dialogue: "*I'm not ready to become a widow*", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "You won't be.\nI'll get the jade wyrm and we can have a chance\ndelivering our babies safely", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewBgChangeEvent(assets.Workshop2, core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: 0, Speed: 3}, nil),
			&core.CharacterImageSwapEvent{Name: "Sven", NewImage: assets.Sven2},
			&core.DialogueEvent{Name: "Sven", Dialogue: "Off I go", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Godspeed, son", FontFace: assets.FontFace},
			&core.StopBgmEvent{},
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
	// scene.Events[0].Execute(scene)
	return scene
}
