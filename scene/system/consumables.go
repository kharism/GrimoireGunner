package system

type Consumables interface {
	GetCharge() int
	SetCharge(int)
}
