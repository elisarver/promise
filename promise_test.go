package promise_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	. "promise"
)

type testtype[T any] struct {
	name        string
	contextFunc func() context.Context
	fn          func(ctx context.Context) (*T, error)
	want        *T
	wantErr     bool
}

func TestPromiseGetString(t *testing.T) {
	msg := "hello"
	tests := []testtype[string]{
		{
			name: "string",
			contextFunc: func() context.Context {
				return context.Background()
			},
			fn: func(ctx context.Context) (*string, error) {
				return &msg, nil
			},
			want:    &msg,
			wantErr: false,
		},
		{
			name: "error",
			contextFunc: func() context.Context {
				return context.Background()
			},
			fn: func(ctx context.Context) (*string, error) {
				return nil, errors.New("error")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "canceled",
			contextFunc: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				return ctx
			},
			fn: func(ctx context.Context) (*string, error) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.testPromiseGet)
	}
}

func (test testtype[T]) testPromiseGet(t *testing.T) {
	ctx := test.contextFunc()

	p := Promise(ctx, test.fn)
	time.Sleep(time.Second / 2)
	got, err := p.Get()

	if (err != nil) != test.wantErr {
		t.Errorf("Wait() error = %v, wantErr %v", err, test.wantErr)
		return
	}

	if !reflect.DeepEqual(got, test.want) {
		t.Errorf("Wait() got = %v, want %v", got, test.want)
	}
}
