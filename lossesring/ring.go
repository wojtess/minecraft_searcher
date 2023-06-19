package lossesring

type LossesRing[T any] struct {
	data  []T
	index int
}

func New[T any](size int) LossesRing[T] {
	return LossesRing[T]{
		data:  make([]T, size),
		index: 0,
	}
}

func (r *LossesRing[T]) Push(value T) {
	if r.index >= len(r.data) {
		r.index = 0
	}
	r.data[r.index] = value
	r.index++
}

func (r *LossesRing[T]) GetArray() []T {
	return append(r.data[r.index:], r.data[:r.index]...)
}
