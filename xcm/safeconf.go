package xcm

type SafeConf[T any] struct {
	conf *T
}

func (sc *SafeConf[T]) GetConf() T {
	conf := sc.conf
	return *conf
}

func (sc *SafeConf[T]) SetConf(conf *T) {
	sc.conf = conf
}

func NewSafeConf[T any](conf *T) *SafeConf[T] {
	return &SafeConf[T]{conf: conf}
}
