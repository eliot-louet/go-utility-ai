package ai

type Actor interface {
	SetTrait(key string, value float64)
	GetTrait(key string) (float64, bool)
	OffsetTrait(key string, amount float64)
}
