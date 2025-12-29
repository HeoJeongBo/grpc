package result

import "fmt"

type Result[T any, E error] struct {
	v   T
	err E
	ok  bool
}

func Err[T any, E error](err E) Result[T, E] {
	return Result[T, E]{
		err: err,
		ok:  false,
	}
}

func Ok[T any, E error](v T) Result[T, E] {
	return Result[T, E]{
		v:  v,
		ok: true,
	}
}

func (r Result[T, E]) Unwrap() (T, E) {
	return r.v, r.err
}

func (r Result[T, E]) UnwrapOr(defaultValue T) T {
	if r.ok {
		return r.v
	}
	return defaultValue
}

func (r Result[T, E]) UnwrapOrElse(f func() T) T {
	if r.ok {
		return r.v
	}
	return f()
}

func (r Result[T, E]) IsOk() bool {
	return r.ok
}

func (r Result[T, E]) IsErr() bool {
	return !r.ok
}

func (r Result[T, E]) Expect(msg string) (T, error) {
	if !r.ok {
		return r.v, fmt.Errorf("%s: %w", msg, r.err)
	}
	return r.v, nil
}
