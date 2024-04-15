package xcm

type SafeConfig[T any] struct {
	conf *T
}

func (sc *SafeConfig[T]) GetConfig() T {
	conf := sc.conf
	return *conf
}

func (sc *SafeConfig[T]) SetConfig(conf *T) {
	sc.conf = conf
}

func NewSafeConfig[T any](conf *T) *SafeConfig[T] {
	return &SafeConfig[T]{conf: conf}
}
