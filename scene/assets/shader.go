package assets

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shader/dakka.kage
var dakkaShader []byte

//go:embed shader/darker.kage
var darkerShader []byte

//go:embed shader/icy.kage
var icyShader []byte

//go:embed shader/cooldown.kage
var cooldownShader []byte

var DakkaShader *ebiten.Shader
var DarkerShader *ebiten.Shader
var IcyShader *ebiten.Shader
var CooldownShader *ebiten.Shader

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
	if CooldownShader == nil {
		var err error
		CooldownShader, err = ebiten.NewShader(cooldownShader)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}
}
