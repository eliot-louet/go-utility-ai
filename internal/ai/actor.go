package ai

type Actor interface {
	SetValue(key string, value any)
	GetValue(key string) (any, bool)
	OffsetValue(key string, amount float64)
}
