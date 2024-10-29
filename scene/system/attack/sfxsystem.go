package attack

type SfxQueuer interface {
	QueueSFX(sfx []byte)
}

var AtkSfxQueue SfxQueuer
