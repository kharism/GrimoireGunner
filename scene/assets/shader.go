package assets

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/component"
)

//go:embed shader/dakka.kage
var dakkaShader []byte

//go:embed shader/darker.kage
var darkerShader []byte

//go:embed shader/icy.kage
var icyShader []byte

//go:embed shader/cooldown.kage
var cooldownShader []byte

//go:embed shader/shocky.kage
var shockyShader []byte

//go:embed shader/shocky2.kage
var shocky2Shader []byte

//go:embed shader/woody.kage
var woodyShader []byte

var DakkaShader *ebiten.Shader
var DarkerShader *ebiten.Shader
var IcyShader *ebiten.Shader
var CooldownShader *ebiten.Shader
var WoodyShader *ebiten.Shader
var ShockyShader *ebiten.Shader
var Shocky2Shader *ebiten.Shader

func Element2Shader(el component.Elemental) *ebiten.Shader {
	switch el {
	case component.ELEC:
		return ShockyShader
	case component.FIRE:
		return DakkaShader
	case component.WATER:
		return IcyShader
	case component.WOOD:
		return WoodyShader
	}
	return nil
}
func init() {
	if DakkaShader == nil {
		DakkaShader, _ = ebiten.NewShader(dakkaShader)
	}
	if DarkerShader == nil {
		DarkerShader, _ = ebiten.NewShader(darkerShader)
	}
	if IcyShader == nil {
		IcyShader, _ = ebiten.NewShader(icyShader)
	}
	if WoodyShader == nil {
		WoodyShader, _ = ebiten.NewShader(woodyShader)
	}
	if ShockyShader == nil {
		ShockyShader, _ = ebiten.NewShader(shockyShader)
	}
	if Shocky2Shader == nil {
		Shocky2Shader, _ = ebiten.NewShader(shocky2Shader)
	}
	if CooldownShader == nil {
		var err error
		CooldownShader, err = ebiten.NewShader(cooldownShader)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}
}
