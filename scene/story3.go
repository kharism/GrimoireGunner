package scene

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/hanashi/core"
)

func Scene3(layouter core.GetLayouter) *core.Scene {
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
			&core.DialogueEvent{Name: "Sven", Dialogue: "(Alright, this is the last portal I need to pass)", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "(Self-diagnostic detect minor anomaly on the suit's spirit flow)", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Sven", Dialogue: "(I can brush it off for now and keep going)", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Still alive son?", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			// &core.PlayBgmEvent{Audio: &assets.MidMusic, Type: core.TypeMP3},
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Alive and kicking", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "For now, I guess", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Sven", Dialogue: "These guys are getting relentless, attacking me\nleft and right", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Great, based on my informant, the best has yet to come", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			// &core.PlayBgmEvent{Audio: &assets.MidMusic, Type: core.TypeMP3},
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "What??", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "Remember the fallen knight stuff I mentioned?", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "Apparently they want to retrieve the suit you use right now", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			// &core.PlayBgmEvent{Audio: &assets.MidMusic, Type: core.TypeMP3},
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "Why?? This suit is just an old decommisioned suit", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "It's the same suit I used to fight their man-made spirit beast", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "Some of its spirit was absorbed by the suit when I defeated it", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Jack", Dialogue: "And they want to take it back", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			// &core.PlayBgmEvent{Audio: &assets.MidMusic, Type: core.TypeMP3},
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "So I have a bounty on my suit?", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Sven", Dialogue: "What a thrill", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			core.NewCharacterRemoveEvent("Sven"),
			core.NewCharacterAddEvent("Jack", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Jack", Dialogue: "I know right?", FontFace: assets.FontFace},
		}},
		&core.DialogueEvent{Name: "Jack", Dialogue: "Also, be careful to the dormant spirit beast on your suit", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Jack", Dialogue: "It's been decades, It might grow into independent spirit beast\non its own", FontFace: assets.FontFace},
		&core.DialogueEvent{Name: "Jack", Dialogue: "Nobody knows what will happen if it wakes.", FontFace: assets.FontFace},
		&core.ComplexEvent{Events: []core.Event{
			// &core.PlayBgmEvent{Audio: &assets.MidMusic, Type: core.TypeMP3},
			core.NewCharacterRemoveEvent("Jack"),
			core.NewCharacterAddEvent("Sven", portraitMoveParam, portraitScaleParam),
			&core.DialogueEvent{Name: "Sven", Dialogue: "...", FontFace: assets.FontFace},
		}},
		&core.ComplexEvent{Events: []core.Event{
			&core.StopBgmEvent{},
			&core.DialogueEvent{Name: "Sven", Dialogue: "Understood", FontFace: assets.FontFace},
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
