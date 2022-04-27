package promise

type promise[T any] struct {
	done chan struct{}
	val  T
	err  error
}

func (p *promise[T]) Get() (T, error) {
	<-p.done

	return p.val, p.err
}

func Promise[T any](fn func() (value T, err error)) *promise[T] {
	p := promise[T]{
		done: make(chan struct{}, 1),
	}

	go func() {
		defer close(p.done)

		p.val, p.err = fn()
	}()

	return &p
}
