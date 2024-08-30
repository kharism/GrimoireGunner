package system

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func RenderLoadOut(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(
		filter.Contains(
			archetype.PlayerTag,
		),
	)

	playerEntry, _ := query.First(ecs.World)
	playerSprite := component.Sprite.Get(playerEntry).Image
	playerScrPos := component.ScreenPos.Get(playerEntry)
	bounds := playerSprite.Bounds()
	if CurLoadOut[0] != nil {
		icon := CurLoadOut[0].GetIcon()
		transformation := ebiten.GeoM{}
		transformation.Scale(0.75, 0.75)
		iconBound := icon.Bounds()
		transformation.Translate(playerScrPos.X-float64(bounds.Dx()/2), playerScrPos.Y-float64(bounds.Dy())-float64(iconBound.Dy())*0.75)
		DrawOp := ebiten.DrawImageOptions{
			GeoM: transformation,
		}
		if CurLoadOut[0].GetCooldown().After(time.Now()) {
			opts := &ebiten.DrawRectShaderOptions{
				GeoM: transformation,
			}
			opts.Images[0] = icon
			// bounds := .Bounds()
			screen.DrawRectShader(iconBound.Dx(), iconBound.Dy(), assets.DakkaShader, opts)
		} else {
			screen.DrawImage(icon, &DrawOp)
		}

	}
	if CurLoadOut[1] != nil {
		icon := CurLoadOut[1].GetIcon()
		transformation := ebiten.GeoM{}
		transformation.Scale(0.75, 0.75)
		iconBound := icon.Bounds()
		transformation.Translate(playerScrPos.X-float64(bounds.Dx()/2)+float64(iconBound.Dx())*0.75, playerScrPos.Y-float64(bounds.Dy())-float64(iconBound.Dy())*0.75)
		DrawOp := ebiten.DrawImageOptions{
			GeoM: transformation,
		}
		if CurLoadOut[1].GetCooldown().After(time.Now()) {
			opts := &ebiten.DrawRectShaderOptions{
				GeoM: transformation,
			}
			opts.Images[0] = icon
			// bounds := .Bounds()
			screen.DrawRectShader(iconBound.Dx(), iconBound.Dy(), assets.DakkaShader, opts)
		} else {
			screen.DrawImage(icon, &DrawOp)
		}
		// screen.DrawImage(icon, &DrawOp)
	}
}
