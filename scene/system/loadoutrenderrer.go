package system

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kharism/grimoiregunner/scene/archetype"
	"github.com/kharism/grimoiregunner/scene/assets"
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/kharism/grimoiregunner/scene/system/loadout"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var MonogramFace *text.GoTextFace

func init() {
	MonogramFace = &text.GoTextFace{
		Source: assets.MonogramFont,
		Size:   35,
	}
}

var LoadOutIconStartX = 40
var LoadOutIconStartY = 520

func RenderLoadOut(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(
		filter.Contains(
			archetype.PlayerTag,
		),
	)

	playerEntry, ok := query.First(ecs.World)
	if !ok {
		return
	}
	playerSprite := component.Sprite.Get(playerEntry).Image
	playerScrPos := component.ScreenPos.Get(playerEntry)
	bounds := playerSprite.Bounds()
	if loadout.CurLoadOut[0] != nil {
		icon := loadout.CurLoadOut[0].GetIcon()
		transformation := ebiten.GeoM{}
		transformation.Scale(0.75, 0.75)
		iconBound := icon.Bounds()
		transformation.Translate(playerScrPos.X-float64(bounds.Dx()/2), playerScrPos.Y-float64(bounds.Dy())-float64(iconBound.Dy())*0.75)
		DrawOp := ebiten.DrawImageOptions{
			GeoM: transformation,
		}
		transformation2 := ebiten.GeoM{}
		transformation2.Scale(2, 2)
		transformation2.Translate(float64(LoadOutIconStartX), float64(LoadOutIconStartY))
		DrawOp2 := ebiten.DrawImageOptions{
			GeoM: transformation2,
		}
		now := time.Now()
		if loadout.CurLoadOut[0].GetCooldown().After(now) {
			opts := &ebiten.DrawRectShaderOptions{
				GeoM: transformation,
			}
			opts.Images[0] = icon
			opts.Uniforms = make(map[string]interface{})

			// bounds := .Bounds()
			vv := float32(loadout.CurLoadOut[0].GetCooldown().Sub(now))
			timeLeftPercentage := 1 - (vv / float32(loadout.CurLoadOut[0].GetCooldownDuration()))
			fmt.Println(timeLeftPercentage)
			opts.Uniforms["Iter"] = float32(timeLeftPercentage)
			screen.DrawRectShader(iconBound.Dx(), iconBound.Dy(), assets.CooldownShader, opts)
			opts2 := &ebiten.DrawRectShaderOptions{
				GeoM: transformation2,
			}
			opts2.Images[0] = icon
			opts2.Uniforms = make(map[string]interface{})
			opts2.Uniforms["Iter"] = float32(timeLeftPercentage)
			screen.DrawRectShader(iconBound.Dx(), iconBound.Dy(), assets.CooldownShader, opts2)
			// dist := loadout.CurLoadOut[0].GetCooldown().Sub(now)
			// textTranslate := ebiten.GeoM{}
			// textTranslate.Translate(float64(LoadOutIconStartX)+10, float64(LoadOutIconStartY)+float64(iconBound.Dy())*1.5)

			// textDrawOpt := text.DrawOptions{
			// 	LayoutOptions: text.LayoutOptions{
			// 		PrimaryAlign: text.AlignCenter,
			// 	},
			// 	DrawImageOptions: ebiten.DrawImageOptions{
			// 		GeoM: textTranslate,
			// 	},
			// }

			// text.Draw(screen, fmt.Sprintf("%.0fs", dist.Seconds()), MonogramFace, &textDrawOpt)
		} else {
			screen.DrawImage(icon, &DrawOp)
			screen.DrawImage(icon, &DrawOp2)
		}
		dmgText := loadout.CurLoadOut[0].GetDamage()
		textTranslate := ebiten.GeoM{}
		// textTranslate.Translate()
		textTranslate.Translate(float64(LoadOutIconStartX)+float64(iconBound.Dx())*2, float64(LoadOutIconStartY)+float64(iconBound.Dy())*1.5)
		textDrawOpt := text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: text.AlignEnd,
			},
			DrawImageOptions: ebiten.DrawImageOptions{
				GeoM: textTranslate,
			},
		}
		text.Draw(screen, fmt.Sprintf("%.d", dmgText), MonogramFace, &textDrawOpt)

	}
	if loadout.CurLoadOut[1] != nil {
		icon := loadout.CurLoadOut[1].GetIcon()
		transformation := ebiten.GeoM{}
		transformation.Scale(0.75, 0.75)
		iconBound := icon.Bounds()
		transformation.Translate(playerScrPos.X-float64(bounds.Dx()/2)+float64(iconBound.Dx())*0.75, playerScrPos.Y-float64(bounds.Dy())-float64(iconBound.Dy())*0.75)
		DrawOp := ebiten.DrawImageOptions{
			GeoM: transformation,
		}
		transformation2 := ebiten.GeoM{}
		transformation2.Scale(2, 2)
		transformation2.Translate(float64(LoadOutIconStartX+iconBound.Dx()*2), float64(LoadOutIconStartY))
		DrawOp2 := ebiten.DrawImageOptions{
			GeoM: transformation2,
		}
		now := time.Now()
		if loadout.CurLoadOut[1].GetCooldown().After(now) {
			opts := &ebiten.DrawRectShaderOptions{
				GeoM: transformation,
			}
			opts.Uniforms = make(map[string]interface{})
			opts.Images[0] = icon
			vv := float32(loadout.CurLoadOut[1].GetCooldown().Sub(now))
			timeLeftPercentage := 1 - (vv / float32(loadout.CurLoadOut[0].GetCooldownDuration()))
			fmt.Println(timeLeftPercentage)
			opts.Uniforms["Iter"] = float32(timeLeftPercentage)
			// bounds := .Bounds()
			opts2 := &ebiten.DrawRectShaderOptions{
				GeoM: transformation2,
			}
			opts2.Images[0] = icon
			opts2.Uniforms = make(map[string]interface{})

			opts2.Uniforms["Iter"] = float32(timeLeftPercentage)
			screen.DrawRectShader(iconBound.Dx(), iconBound.Dy(), assets.CooldownShader, opts2)
			// dist := loadout.CurLoadOut[1].GetCooldown().Sub(now)
			// textTranslate := ebiten.GeoM{}
			// textTranslate.Translate(float64(LoadOutIconStartX+iconBound.Dx()*2+10), float64(LoadOutIconStartY)+float64(iconBound.Dy())*1.5)

			// textDrawOpt := text.DrawOptions{
			// 	LayoutOptions: text.LayoutOptions{
			// 		PrimaryAlign: text.AlignCenter,
			// 	},
			// 	DrawImageOptions: ebiten.DrawImageOptions{
			// 		GeoM: textTranslate,
			// 	},
			// }

			// text.Draw(screen, fmt.Sprintf("%.0fs", dist.Seconds()), MonogramFace, &textDrawOpt)
			screen.DrawRectShader(iconBound.Dx(), iconBound.Dy(), assets.CooldownShader, opts)
		} else {
			screen.DrawImage(icon, &DrawOp)
			screen.DrawImage(icon, &DrawOp2)
		}
		dmgText := loadout.CurLoadOut[1].GetDamage()
		textTranslate := ebiten.GeoM{}
		// textTranslate.Translate()
		textTranslate.Translate(float64(LoadOutIconStartX+iconBound.Dx()*4), float64(LoadOutIconStartY)+float64(iconBound.Dy())*1.5)
		textDrawOpt := text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: text.AlignEnd,
			},
			DrawImageOptions: ebiten.DrawImageOptions{
				GeoM: textTranslate,
			},
		}
		text.Draw(screen, fmt.Sprintf("%.d", dmgText), MonogramFace, &textDrawOpt)
		// screen.DrawImage(icon, &DrawOp)
	}
	Sub1StartIconX := LoadOutIconStartX
	if loadout.SubLoadOut1[0] != nil {
		icon := loadout.SubLoadOut1[0].GetIcon()
		transformation := ebiten.GeoM{}
		transformation.Scale(2, 2)
		iconBound := icon.Bounds()
		transformation.Translate(float64(Sub1StartIconX+2*iconBound.Dx()*2), float64(LoadOutIconStartY))
		DrawOp := ebiten.DrawImageOptions{
			GeoM: transformation,
		}
		now := time.Now()
		if loadout.SubLoadOut1[0].GetCooldown().After(now) {
			opts2 := &ebiten.DrawRectShaderOptions{
				GeoM: transformation,
			}
			vv := float32(loadout.CurLoadOut[0].GetCooldown().Sub(now))
			timeLeftPercentage := 1 - (vv / float32(loadout.CurLoadOut[0].GetCooldownDuration()))
			opts2.Images[0] = icon
			opts2.Uniforms = make(map[string]interface{})
			opts2.Uniforms["Iter"] = float32(timeLeftPercentage)
			screen.DrawRectShader(iconBound.Dx(), iconBound.Dy(), assets.CooldownShader, opts2)
			// dist := loadout.SubLoadOut1[0].GetCooldown().Sub(now)
			// textTranslate := ebiten.GeoM{}
			// textTranslate.Translate(float64(Sub1StartIconX+2*iconBound.Dx()*2+10), float64(LoadOutIconStartY)+float64(iconBound.Dy())*1.5)

			// textDrawOpt := text.DrawOptions{
			// 	LayoutOptions: text.LayoutOptions{
			// 		PrimaryAlign: text.AlignCenter,
			// 	},
			// 	DrawImageOptions: ebiten.DrawImageOptions{
			// 		GeoM: textTranslate,
			// 	},
			// }

			// text.Draw(screen, fmt.Sprintf("%.0fs", dist.Seconds()), MonogramFace, &textDrawOpt)
		} else {
			screen.DrawImage(icon, &DrawOp)
		}
	}
	if loadout.SubLoadOut1[1] != nil {
		icon := loadout.SubLoadOut1[1].GetIcon()
		transformation := ebiten.GeoM{}
		transformation.Scale(2, 2)
		iconBound := icon.Bounds()
		transformation.Translate(float64(Sub1StartIconX+3*iconBound.Dx()*2), float64(LoadOutIconStartY))
		DrawOp := ebiten.DrawImageOptions{
			GeoM: transformation,
		}
		now := time.Now()
		if loadout.SubLoadOut1[1].GetCooldown().After(now) {
			vv := float32(loadout.CurLoadOut[0].GetCooldown().Sub(now))
			timeLeftPercentage := 1 - (vv / float32(loadout.CurLoadOut[0].GetCooldownDuration()))
			// fmt.Println(timeLeftPercentage)
			opts2 := &ebiten.DrawRectShaderOptions{
				GeoM: transformation,
			}
			opts2.Images[0] = icon
			opts2.Uniforms = make(map[string]interface{})
			opts2.Uniforms["Iter"] = float32(timeLeftPercentage)
			screen.DrawRectShader(iconBound.Dx(), iconBound.Dy(), assets.CooldownShader, opts2)
			// dist := loadout.SubLoadOut1[1].GetCooldown().Sub(now)
			// textTranslate := ebiten.GeoM{}
			// textTranslate.Translate(float64(Sub1StartIconX+3*iconBound.Dx()*2+10), float64(LoadOutIconStartY)+float64(iconBound.Dy())*1.5)

			// textDrawOpt := text.DrawOptions{
			// 	LayoutOptions: text.LayoutOptions{
			// 		PrimaryAlign: text.AlignCenter,
			// 	},
			// 	DrawImageOptions: ebiten.DrawImageOptions{
			// 		GeoM: textTranslate,
			// 	},
			// }

			// text.Draw(screen, fmt.Sprintf("%.0fs", dist.Seconds()), MonogramFace, &textDrawOpt)
		} else {
			screen.DrawImage(icon, &DrawOp)
		}
	}
	if loadout.SubLoadOut2[0] != nil {
		icon := loadout.SubLoadOut2[0].GetIcon()
		transformation := ebiten.GeoM{}
		transformation.Scale(2, 2)
		iconBound := icon.Bounds()
		transformation.Translate(float64(Sub1StartIconX+4*iconBound.Dx()*2), float64(LoadOutIconStartY))
		DrawOp := ebiten.DrawImageOptions{
			GeoM: transformation,
		}
		now := time.Now()
		if loadout.SubLoadOut2[0].GetCooldown().After(now) {
			vv := float32(loadout.CurLoadOut[0].GetCooldown().Sub(now))
			timeLeftPercentage := 1 - (vv / float32(loadout.CurLoadOut[0].GetCooldownDuration()))
			// fmt.Println(timeLeftPercentage)

			opts2 := &ebiten.DrawRectShaderOptions{
				GeoM: transformation,
			}
			opts2.Uniforms = make(map[string]interface{})
			opts2.Uniforms["Iter"] = float32(timeLeftPercentage)
			opts2.Images[0] = icon
			screen.DrawRectShader(iconBound.Dx(), iconBound.Dy(), assets.CooldownShader, opts2)
			// dist := loadout.SubLoadOut2[0].GetCooldown().Sub(now)
			// textTranslate := ebiten.GeoM{}
			// textTranslate.Translate(float64(Sub1StartIconX+4*iconBound.Dx()*2+10), float64(LoadOutIconStartY)+float64(iconBound.Dy())*1.5)

			// textDrawOpt := text.DrawOptions{
			// 	LayoutOptions: text.LayoutOptions{
			// 		PrimaryAlign: text.AlignCenter,
			// 	},
			// 	DrawImageOptions: ebiten.DrawImageOptions{
			// 		GeoM: textTranslate,
			// 	},
			// }

			// text.Draw(screen, fmt.Sprintf("%.0fs", dist.Seconds()), MonogramFace, &textDrawOpt)
		} else {
			screen.DrawImage(icon, &DrawOp)
		}
	}
	if loadout.SubLoadOut2[1] != nil {
		icon := loadout.SubLoadOut2[1].GetIcon()
		transformation := ebiten.GeoM{}
		transformation.Scale(2, 2)
		iconBound := icon.Bounds()
		transformation.Translate(float64(Sub1StartIconX+5*iconBound.Dx()*2), float64(LoadOutIconStartY))
		DrawOp := ebiten.DrawImageOptions{
			GeoM: transformation,
		}
		now := time.Now()
		if loadout.SubLoadOut2[1].GetCooldown().After(now) {
			vv := float32(loadout.CurLoadOut[0].GetCooldown().Sub(now))
			timeLeftPercentage := 1 - (vv / float32(loadout.CurLoadOut[0].GetCooldownDuration()))
			opts2 := &ebiten.DrawRectShaderOptions{
				GeoM: transformation,
			}
			opts2.Images[0] = icon
			opts2.Uniforms = make(map[string]interface{})
			opts2.Uniforms["Iter"] = float32(timeLeftPercentage)
			screen.DrawRectShader(iconBound.Dx(), iconBound.Dy(), assets.CooldownShader, opts2)
			// dist := loadout.SubLoadOut2[1].GetCooldown().Sub(now)
			// textTranslate := ebiten.GeoM{}
			// textTranslate.Translate(float64(Sub1StartIconX+5*iconBound.Dx()*2+10), float64(LoadOutIconStartY)+float64(iconBound.Dy())*1.5)

			// textDrawOpt := text.DrawOptions{
			// 	LayoutOptions: text.LayoutOptions{
			// 		PrimaryAlign: text.AlignCenter,
			// 	},
			// 	DrawImageOptions: ebiten.DrawImageOptions{
			// 		GeoM: textTranslate,
			// 	},
			// }

			// text.Draw(screen, fmt.Sprintf("%.0fs", dist.Seconds()), MonogramFace, &textDrawOpt)
		} else {
			screen.DrawImage(icon, &DrawOp)
		}
	}
}
