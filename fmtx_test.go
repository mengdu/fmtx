package fmtx

import (
	"errors"
	"fmt"
	"testing"
)

type MyInt int

type Data struct {
	A      int
	B      float32
	C      string
	D      bool
	E      []int
	F      map[string]int
	fn     func() error
	err    error
	any    any
	anyMap map[any]any
	loop   *Data
}

func (d Data) SayHello() {
	d.hi()
}

func (d Data) hi() {}

func genData() Data {
	data := Data{
		A: 1,
		B: 3.14,
		C: "Hello \nworld",
		D: true,
		E: []int{1, 2, 3},
		F: map[string]int{
			"a": 1,
			"b": 2,
		},
		fn: func() error {
			return nil
		},
		err: errors.New("some error"),
		any: []any{genData},
	}
	data.loop = &data
	data.anyMap = map[any]any{
		1: data,
		2: []any{data},
	}
	return data
}

func BenchmarkString(b *testing.B) {
	EnableColor = true
	data := genData()
	for i := 0; i < b.N; i++ {
		String(data)
	}
}

func BenchmarkStringDisableColor(b *testing.B) {
	EnableColor = false
	data := genData()
	for i := 0; i < b.N; i++ {
		String(data)
	}
}

func BenchmarkSprintf(b *testing.B) {
	data := genData()
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%#v", data)
	}
}
