package enemies

type SfxQueuer interface {
	QueueSFX(sfx []byte)
}

var EnemySfxQueue SfxQueuer
