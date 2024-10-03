package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
)

// add max HP on pickup
type HPUp struct {
}

func (h *HPUp) GetIcon() *ebiten.Image {
	return assets.HPUpIcon
}
func (h *HPUp) GetDescription() string {
	return "Add 100 max HP"
}
func (h *HPUp) GetName() string {
	return "HP UP"
}
func (h *HPUp) OnAquireDo(data *SceneData) {
	data.PlayerMaxHP += 100
	data.PlayerHP += 100
}
