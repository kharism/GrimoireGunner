package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
)

// add max HP on pickup
type Medkit struct {
}

func (h *Medkit) GetIcon() *ebiten.Image {
	return assets.MedkitIcon
}
func (h *Medkit) GetDescription() string {
	return "Recover 200HP"
}
func (h *Medkit) GetName() string {
	return "Medkit"
}
func (h *Medkit) OnAquireDo(data *SceneData) {
	data.PlayerHP += 200
	if data.PlayerHP > data.PlayerMaxHP {
		data.PlayerHP = data.PlayerMaxHP
	}
}
