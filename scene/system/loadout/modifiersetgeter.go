package loadout

type ModifierGetSetter interface {
	GetModifierEntry() *CasterModifierData
	SetModifier(*CasterModifierData)
}
