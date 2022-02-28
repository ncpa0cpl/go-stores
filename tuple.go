package stores

type Tuple[T, U any] struct {
	a T
	b U
}

func NewTuple[T, U any](a T, b U) Tuple[T, U] {
	return Tuple[T, U]{a, b}
}

func (t *Tuple[T, U]) Get() (T, U) {
	return t.a, t.b
}
