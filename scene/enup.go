package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kharism/grimoiregunner/scene/assets"
)

// add max HP on pickup
type ENUp struct {
}

func (h *ENUp) GetIcon() *ebiten.Image {
	return assets.ENUpIcon
}
func (h *ENUp) GetDescription() string {
	return "Add 1 max EN"
}
func (h *ENUp) GetName() string {
	return "EN UP"
}
func (h *ENUp) OnAquireDo(data *SceneData) {
	data.PlayerMaxHP += 100
	data.PlayerHP += 100
}
