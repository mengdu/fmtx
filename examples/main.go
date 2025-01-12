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
type MyFloat float32
type MyString string
type MyBool bool

type Foo struct {
	Key  int
	Loop *Foo
}

type Data struct {
	A      int
	B      float32
	C      string
	D      bool
	E      []int
	F      map[string]int
	G      *map[string]int
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
	Loop   *Data
	Foo
	AnonymousField struct {
		One            string
		AnonymousField struct {
			Two string
		}
	}
}

func (d *Data) Method1() (int, error) {
	return 0, nil
}

func (d *Data) Method2() (string, error) {
	return "", nil
}

func (d Data) Method3() (string, error) {
	return "", nil
}

func (Data) Method4(a int, b float32, c string, d map[string]any, e []int, f func(a int) (int, error), g ...[]int) (int, error) {
	return 0, nil
}

func (d Data) Method5() {}

func genData() *Data {
	data := &Data{
		A: 1,
		B: 3.14,
		C: "Hello \nworld",
		D: true,
		E: []int{1, 2, 3},
		F: map[string]int{
			"a": 1,
			"b": 2,
		},
		G:    &map[string]int{},
		Nil1: nil,
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
	data.anyMap = map[any]any{
		1: 1,
	}
	for i := 0; i < 12; i++ {
		data.anyMap[i] = i
		data.anyMap[fmt.Sprintf("key-%d", i)] = i
	}
	foo := Foo{
		Key: 123,
	}
	foo.Loop = &foo
	data.Foo = foo
	data.AnonymousField = struct {
		One            string
		AnonymousField struct{ Two string }
	}{
		One: "one",
		AnonymousField: struct{ Two string }{
			Two: "two",
		},
	}
	return data
}

func main() {
	// fmtx.SetEnableColor(true)
	// fmtx.Options.MaxPropertyBreakLine = 3
	// fmtx.Options.ShowStructMethod = false
	var initMap map[string]int
	var initArr [2]int
	var initSlice []int
	var initChan chan int
	var initFn func() error
	var initAny any
	var initErr error

	fmtx.Println("Hello, \nworld", 'a', ptr("Hi \"Tom\""))
	fmtx.Println(true, false)
	fmtx.Println(123, 3.14, int32(123), ptr(124))
	fmtx.Println(123, 3.14, complex64(1+2i), complex(3.14, -2.71))
	fmtx.Println(MyInt(123), MyFloat(3.14), MyString("string"), MyBool(true))
	fmtx.Println(ptr(MyInt(123)), ptr(MyFloat(3.14)), ptr(MyString("string")), ptr(MyBool(true)))
	fmtx.Println(initMap, initArr, initSlice, initChan, initFn, initAny, initErr)
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
	fmtx.Println(genData())
	foo := Foo{
		Key: 123,
	}
	foo.Loop = &foo
	fmtx.Println(foo)
	fmtx.Println(fmtx.Default)
	fn := func(a int, b float32, c string, d map[string]any, e []int, f ...[]int) error {
		return nil
	}
	fmtx.Println(ptr[int], fmt.Println, func(a any) {}, fn)
	var sendOnlyCh chan<- int
	var readOnlyCh <-chan int
	ch1 := make(chan int)
	ch2 := make(chan int, 4)
	ch2 <- 1
	ch2 <- 2
	fmtx.Println(sendOnlyCh, readOnlyCh, ch1, ch2)
	fmtx.Println(errors.New("some error"))
	// fmtx.Println(fmtx.Color(" Hello ", "40", "49"))
	myOpt := fmtx.Default
	myOpt.MaxDepth = 2
	myOpt.ColorMap.Bool = [2]string{"35", "39"}
	dump := fmtx.New(&myOpt)
	fmt.Println(dump(true), dump(false))
}
