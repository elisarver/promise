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

func (p *promise[T]) Wait() (*T, error) {
	select {
	case <-p.ctx.Done():
		return nil, p.ctx.Err()
	case res := <-p.valCh:
		return res, nil
	case err := <-p.errCh:
		return nil, err
	}
}

var ErrNoImmediateResults = errors.New("no immediate results")

// Any returns any of value or error, IF IT IS IMMEDIATELY AVAILABLE,
// or else returns ErrNoImmediateResults when exit conditions have
// not been met.
func (p *promise[T]) Any() (*T, error) {
	select {
	case <-p.ctx.Done():
		return nil, p.ctx.Err()
	case res := <-p.valCh:
		return res, nil
	case err := <-p.errCh:
		return nil, err
	default:
		return nil, ErrNoImmediateResults
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
