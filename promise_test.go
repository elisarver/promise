package promise_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	. "promise"
)

type testtype[T any] struct {
	name        string
	contextFunc func() context.Context
	fn          func(ctx context.Context) (*T, error)
	want        *T
	wantErr     bool
}

func TestPromiseString(t *testing.T) {
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
				return &msg, nil
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.testPromiseWait)
	}
}

func (test testtype[T]) testPromiseWait(t *testing.T) {
	t.Helper()
	t.Run(test.name, func(t *testing.T) {
		ctx := test.contextFunc()

		p := Promise(ctx, test.fn)
		got, err := p.Wait()

		if (err != nil) != test.wantErr {
			t.Errorf("Wait() error = %v, wantErr %v", err, test.wantErr)
			return
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Wait() got = %v, want %v", got, test.want)
		}
	})
}
