package attack

import "github.com/kharism/grimoiregunner/scene/component"

type ModifierGetSetter interface {
	GetModifierEntry() *component.CasterModifierData
	SetModifier(*component.CasterModifierData)
}
