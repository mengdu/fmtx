package main

import (
	"errors"
	"fmt"

	"github.com/mengdu/fmtx"
)

func ptr[T any](a T) *T {
	return &a
}

type MyInt int

type Data struct {
	A      int
	B      float32
	C      string
	D      bool
	E      []int
	F      map[string]int
	Nil    *int
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

func (d Data) hi() {}

func main() {
	// fmtx.EnableColor = false
	// fmtx.Options.MaxDepth = 5
	fmtx.Println("Hello, \nworld", ptr("Hi \"Tom\""))
	fmtx.Println(true, false)
	fmtx.Println(123, 3.14, ptr(124), MyInt(123), ptr(MyInt(123)))
	fmtx.Println([]int{1, 2, 3, -1, -2})
	fmtx.Println([]string{"a", "b", "c", "Hello \n \"world\"."})
	fmtx.Println([]byte("Hello, \nworld"))
	s := make([]any, 0, 10)
	s = append(s, true)
	s = append(s, 123)
	s = append(s, "abc")
	fmtx.Println(s)
	fmtx.Println([]MyInt{1, 2, 3})
	fmtx.Println([]*int{ptr(1), ptr(2), ptr(3)})
	fmtx.Println(ptr([]int{1, 2, 3}))
	fmtx.Println(ptr([]*int{ptr(1), ptr(2), ptr(3)}))
	fmtx.Println([]*MyInt{ptr(MyInt(1)), ptr(MyInt(2)), ptr(MyInt(3))})
	fmtx.Println([][]int{{1}, {2}, {3, 4}})
	fmtx.Println([100]int{})
	fmtx.Println([3]string{})
	fmtx.Println(map[string]string{
		"a": "a",
		"b": "b",
		"c": "hello \n world!",
	})
	fmtx.Println(map[string]any{
		"a":     true,
		"b":     3.14,
		"a\nbc": "hello \n world!",
	})
	fmtx.Println(map[MyInt]MyInt{
		1: 1,
	})
	fmtx.Println(map[MyInt]map[any]any{
		1: {
			2: 2,
		},
	})
	fmtx.Println(map[any]any{
		1:       true,
		true:    3.14,
		"a\nbc": "hello \n world!",
	})
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
		Nil: nil,
		fn: func() error {
			return nil
		},
		ch1: make(chan<- int),
		ch2: make(<-chan int),
		ch3: make(chan int),
		ch4: make(chan int, 4),
		err: errors.New("some error"),
		any: []any{ptr[string]},
	}
	data.loop = &data
	data.anyMap = map[any]any{
		1: data,
		2: []any{data},
	}

	fmtx.Println(data)
	fmtx.Println(ptr[int], fmt.Println, func(a any) {})
	var sendOnlyCh chan<- int
	var readOnlyCh <-chan int
	ch1 := make(chan int)
	ch2 := make(chan int, 4)
	ch2 <- 1
	ch2 <- 2
	fmtx.Println(sendOnlyCh, readOnlyCh, ch1, ch2)
	fmtx.Println(errors.New("some error"))
}
