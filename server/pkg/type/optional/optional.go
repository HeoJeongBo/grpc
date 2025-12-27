package optional

type Optional[T any] struct {
	value T
	ok    bool
}

func Some[T any](v T) Optional[T] {
	return Optional[T]{value: v, ok: true}
}

func None[T any]() Optional[T] {
	return Optional[T]{ok: false}
}

func (o Optional[T]) IsSome() bool {
	return o.ok
}

func (o Optional[T]) IsNone() bool {
	return !o.ok
}

func (o Optional[T]) Unwrap() (T, bool) {
	return o.value, o.ok
}

func (o Optional[T]) UnwrapOr(defaultValue T) T {
	if o.ok {
		return o.value
	}
	return defaultValue
}

func (o Optional[T]) UnwrapOrElse(f func() T) T {
	if o.ok {
		return o.value
	}
	return f()
}

func (o Optional[T]) Map(f func(T) T) Optional[T] {
	if !o.ok {
		return None[T]()
	}
	return Some(f(o.value))
}

func (o Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if o.ok && predicate(o.value) {
		return o
	}
	return None[T]()
}

func (o Optional[T]) OrElse(alternative Optional[T]) Optional[T] {
	if o.ok {
		return o
	}
	return alternative
}
