package fmtx

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type MyInt int

type Data struct {
	A      int
	B      float32
	C      string
	D      bool
	E      []int
	F      map[string]int
	Nil1   *int
	Nil2   *[]int
	Nil3   *map[string]any
	Nil4   *MyInt
	Nil5   *Data
	Nil6   chan int
	fn     func() error
	ch1    chan<- int
	ch2    <-chan int
	ch3    chan int
	ch4    chan int
	err    error
	any    any
	anyMap map[any]any
	loop   *Data
}

func (d Data) SayHello() {
	d.hi()
}

func (Data) Fn(a int, b float32, c string, d map[string]any, e []int, f ...[]int) (int, error) {
	return 0, nil
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
		Nil1: nil,
		fn: func() error {
			return nil
		},
		ch1: make(chan<- int),
		ch2: make(<-chan int),
		ch3: make(chan int),
		ch4: make(chan int, 4),
		err: errors.New("some error"),
		any: []any{},
	}
	data.loop = &data
	data.anyMap = map[any]any{
		1: 1,
	}
	return data
}

func BenchmarkString(b *testing.B) {
	SetEnableColor(true)
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = String(data)
	}
}

func BenchmarkStringDisableColor(b *testing.B) {
	SetEnableColor(false)
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = String(data)
	}
}

func BenchmarkSprintf(b *testing.B) {
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%#v", data)
	}
}

func TestPrint(t *testing.T) {
	data := genData()
	c := 8
	defer func() {
		cacheEnableColor = 0
	}()
	for i := 0; i < c; i++ {
		start := time.Now()
		_ = fmt.Sprintf("%#v\n", data)
		fmt.Println(1, time.Since(start).String())
	}
	fmt.Println("")
	for i := 0; i < c; i++ {
		start := time.Now()
		_ = color("%#v\n", "35", "39")
		fmt.Println(1, time.Since(start).String())
	}
	fmt.Println("")
	cacheEnableColor = 0
	for i := 0; i < c; i++ {
		start := time.Now()
		_ = String(data)
		fmt.Println(2, time.Since(start).String())
	}
	fmt.Println("")
	cacheEnableColor = 1
	for i := 0; i < c; i++ {
		start := time.Now()
		_ = String(data)
		fmt.Println(3, time.Since(start).String())
	}
}
