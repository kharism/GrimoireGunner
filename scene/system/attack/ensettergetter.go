package attack

type ENSetGetter interface {
	SetEn(val int)
	GetEn() int
	GetMaxEn() int
}
