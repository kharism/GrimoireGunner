package scene

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/hanashi/core"
)

// this executed at the end of level1. Describe sven's journey so far
func Scene2(layouter core.GetLayouter) *core.Scene {
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
	}
	scene.Events = []core.Event{
		&core.ComplexEvent{Events: []core.Event{
			&core.PlayBgmEvent{Audio: &assets.MidMusic, Type: core.TypeMP3},
			core.NewBgChangeEvent(assets.Portal1, core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: 0, Speed: 3}, nil),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "I'm near the first portal", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "There were some weird encounter but I'm fine", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "What encounter?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "I fought some people with grimoire system", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "They seems better equipped to combat a stalker unit\nthan a simple bandit", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "....", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "I always keep track with factions on the underworld,\nbut I can't seem to recognize their insignia", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Can you describe the insignia?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "A winged rhino in front of a triangle", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Assuming that was a unicorn and not rhino\nIt's probably fallen knight of crimson king", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "It's been a long time since I cleaned them up", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "How long since you cleaned them?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "I remember you crapped on my hand when I bathed you\nthe day after I cleaned them up", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "DAD!!!", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Anyway, Shizuku here would like to see you", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "How's going over there?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Don't worry. Just regular monster here", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "Just as regular as breakfast", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Sven", Dialogue: "You still dizzy?", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "Yes. I just vomit a moment ago", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "You still eat regularly don't you?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "NO!!!", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "FATHER!!!", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Yep, that's my father", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "Just try to eat as much as you can", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Sven", Dialogue: "Once I back, we'll make smoked grilled meat", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "*pick a bag and vomit off screen*", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Or not. We'll decide on how to process the meat later then", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "I'll get going", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Sven", Dialogue: "I leave shizuku to you dad", FontFace: assets.FontFace},
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
