package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/hanashi/core"
)

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
			// &core.PlayBgmEvent{Audio: &assets.SmoothJazz, Type: core.TypeMP3},
			core.NewBgChangeEvent(assets.Portal1, core.MoveParam{Sx: 0, Sy: 0, Tx: 0, Ty: 0, Speed: 3}, nil),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "I'm near the first portal"},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "There were some weird encounter but I'm fine"},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "What encounter?"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "I fought some people with grimoire system"},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "They seems better equipped to combat a stalker unit\nthan a simple bandit"},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "...."},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "I always keep track with factions on the underworld,\nbut I can't seem to recognize their insignia"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Can you describe the insignia?"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "A winged rhino in front of a triangle"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "That's fallen knight of crimson king"},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "It's been a long time since I cleaned them"},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "How long since you cleaned them?"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "I remember you crapped on my hand when I bathed you\the day after I cleaned them up"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "DAD!!!"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Anyway, Shizuku here would like to see you"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "How's going over there?"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Don't worry. Nothing out of typical monster here."},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "Nothing I don't deal with almost everyday"},
		&core.DialogueEvent{Name: "Sven", Dialogue: "You still dizzy?"},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "Yes. I just vomit a moment ago"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "You still eat regularly don't you?"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "NO!!!"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "FATHER!!!"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Yep, that's my father"},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "Just try to eat as much as you can"},
		&core.DialogueEvent{Name: "Sven", Dialogue: "Once I back, we'll make smoked grilled meat"},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Shizuku", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Shizuku", Dialogue: "*pick a bag and vomit off screen*"},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Shizuku"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Or not. We'll decide on how to process the meat later then"},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "I'll get going"},
		&core.DialogueEvent{Name: "Sven", Dialogue: "I leave shizuku to you dad"},
	}
	scene.TxtBg = ebiten.NewImage(1024-128, 128)
	scene.TxtBg.Fill(color.RGBA{R: 0x4f, G: 0x8f, B: 0xba, A: 255})
	// pp, err := core.NewDefaultAudioInterfacer()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(-1)
	// }
	scene.AudioInterface = nil
	// scene.Events[0].Execute(scene)
	return scene
}
