package promise

import (
	"context"
)

type promise[T any] struct {
	done chan struct{}

	val *T
	err error
}

func (p *promise[T]) Get() (*T, error) {
	<-p.done

	return p.val, p.err
}

func Promise[T any](ctx context.Context, fn func(ctx context.Context) (value *T, err error)) *promise[T] {
	p := promise[T]{
		done: make(chan struct{}, 1),
	}

	go func() {
		defer close(p.done)

		p.val, p.err = fn(ctx)
	}()

	return &p
}
