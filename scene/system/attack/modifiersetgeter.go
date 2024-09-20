package attack

import "github.com/yohamta/donburi"

type ModifierGetSetter interface {
	GetModifierEntry() *donburi.Entry
	SetModifier(*donburi.Entry)
}
