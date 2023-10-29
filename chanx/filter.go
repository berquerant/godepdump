package chanx

type Filter[T any] func(T) bool

func (f Filter[T]) Call(t T) bool {
	if f == nil {
		return true
	}
	return f(t)
}

func (f Filter[T]) And(g Filter[T]) Filter[T] {
	return func(t T) bool {
		return f.Call(t) && g.Call(t)
	}
}

func DebugFilter[T any](f func(T)) Filter[T] {
	return func(t T) bool {
		f(t)
		return true
	}
}
