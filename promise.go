package promise

import (
	"context"
	"errors"
)

type promise[T any] struct {
	fn  func(context.Context) (*T, error)
	ctx context.Context

	valCh chan *T
	errCh chan error
}

func (p *promise[T]) Fulfill() (*T, error) {
	select {
	case <-p.ctx.Done():
		return nil, errors.New("context canceled")
	case res := <-p.valCh:
		return res, nil
	case err := <-p.errCh:
		return nil, err
	}
}

func Promise[T any](ctx context.Context, fn func(ctx context.Context) (value *T, err error)) *promise[T] {
	p := promise[T]{
		fn:    fn,
		ctx:   ctx,
		valCh: make(chan *T, 1),
		errCh: make(chan error, 1),
	}

	go p.run(ctx)

	return &p
}

func (p *promise[T]) run(ctx context.Context) {
	value, err := p.fn(ctx)
	if err != nil {
		p.errCh <- err

		return
	}

	p.valCh <- value
}
