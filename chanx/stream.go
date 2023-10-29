package chanx

type Stream[T any] interface {
	C() <-chan T
	Filter(filter Filter[T]) Stream[T]
	IntoSlice() []T
}

func NewStream[T any](c <-chan T) Stream[T] {
	return &stream[T]{
		c: c,
	}
}

type stream[T any] struct {
	c      <-chan T
	filter Filter[T]
}

func (s *stream[T]) C() <-chan T {
	resultC := make(chan T, 100)
	go func() {
		defer close(resultC)
		for x := range s.c {
			if s.filter.Call(x) {
				resultC <- x
			}
		}
	}()
	return resultC
}

func (s *stream[T]) Filter(filter Filter[T]) Stream[T] {
	s.filter = s.filter.And(filter)
	return s
}

func NewStreamFromSlice[T any](list []T) Stream[T] {
	resultC := make(chan T, 100)
	go func() {
		defer close(resultC)
		for _, x := range list {
			resultC <- x
		}
	}()
	return NewStream(resultC)
}

func (s *stream[T]) IntoSlice() []T {
	list := []T{}
	for x := range s.C() {
		list = append(list, x)
	}
	return list
}
