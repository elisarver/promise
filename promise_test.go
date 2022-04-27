package promise_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	. "promise"
)

type testtype[T any] struct {
	name    string
	fn      func() (T, error)
	want    T
	wantErr bool
}

func TestPromiseGetString(t *testing.T) {
	tests := []testtype[string]{
		{
			name: "value",
			fn: func() (string, error) {
				return "hello", nil
			},
			want:    "hello",
			wantErr: false,
		},
		{
			name: "error",
			fn: func() (string, error) {
				return "", errors.New("error")
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.testPromiseGet)
	}
}

func TestPromiseGetStringPointer(t *testing.T) {
	msg := "hello"
	tests := []testtype[*string]{
		{
			name: "value",
			fn: func() (*string, error) {
				return &msg, nil
			},
			want:    &msg,
			wantErr: false,
		},
		{
			name: "error",
			fn: func() (*string, error) {
				return nil, errors.New("error")
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
	p := Promise(test.fn)
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
