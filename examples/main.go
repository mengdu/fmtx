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
	Int            int
	Int8           int8
	Int32          int32
	Int64          int64
	Uint           uint
	Uint32         uint32
	Uint64         uint64
	Float32        float32
	Float64        float64
	True           bool
	False          bool
	Char           rune
	String         string
	Byte           []byte
	Ptr1           *int
	Ptr2           *string
	Complex64      complex64
	Complex128     complex128
	Array1         [5]int
	Array2         [4]string
	Array3         [2][3]int
	Array4         [1000]int
	Array5         [1000]string
	Array6         [4]*int
	Array7         []*string
	Slice1         []int
	Slice2         []string
	Slice3         []any
	Slice4         []*int
	Map1           map[string]string
	Map2           map[string]any
	Map3           map[any]any
	Map4           map[any]*int
	Map5           map[*int]*int
	Any1           any
	Any2           any
	Nil1           *int
	Nil2           *[]int
	Nil3           *map[string]any
	Nil4           *MyInt
	Nil5           *Data
	Nil6           chan int
	Nil7           func() any
	Chan1          chan int
	Chan2          chan int
	SendOnlyChan   chan<- int
	ReadOnlyChan   <-chan int
	Foo            Foo
	Loop           *Data
	Fn             func() error
	AnonymousField struct {
		One            string
		AnonymousField struct {
			Two string
		}
	}
}

func (d *Data) M1() (int, error) {
	return 0, nil
}

func (d *Data) M2() (string, error) {
	return "", nil
}

func (d Data) M3() (string, error) {
	return "", nil
}

func (Data) M4(a int, b float32, c string, d map[string]any, e []int, f func(a int) (int, error), g ...[]int) (int, error) {
	return 0, nil
}

func (d Data) M5() {
	d.inner1()
	d.inner2()
}

func (d *Data) inner1() {}
func (d Data) inner2()  {}

func genData() *Data {
	data := &Data{
		Int:        123,
		Int8:       123,
		Int32:      123,
		Int64:      123,
		Uint:       123,
		Uint32:     123,
		Uint64:     123,
		Float32:    3.14,
		Float64:    3.14,
		True:       true,
		False:      false,
		Char:       'a',
		String:     "Hello, \nworld.",
		Byte:       []byte("Hello, \nworld."),
		Ptr1:       ptr(123),
		Ptr2:       ptr("Hello"),
		Complex64:  complex64(1 + 2i),
		Complex128: complex(3.14, -2.71),
		Array1:     [5]int{1, 2, 3, 4, 5},
		Array2:     [4]string{"a", "b", "c", "d"},
		Array3:     [2][3]int{},
		Array4:     [1000]int{},
		Slice1:     make([]int, 4, 8),
		Slice2:     make([]string, 2, 5),
		Slice3:     []any{true, false, 123, 3.145, 'A', "Hello, \nworld.", []int{1, 2, 3}, map[string]any{"a": 1, "b": true}},
		Map1: map[string]string{
			"a": "a",
			"b": "b",
			"c": "c",
		},
		Map2: map[string]any{
			"a": 1,
			"b": true,
			"c": 3.14,
		},
		Map3: map[any]any{},
		Map4: map[any]*int{
			1:      ptr(1),
			3.14:   ptr(2),
			true:   ptr(3),
			ptr(1): ptr(4),
		},
		Map5: map[*int]*int{
			ptr(1): ptr(111),
			ptr(2): ptr(222),
		},
		Chan1: make(chan int),
		Chan2: make(chan int, 6),
	}
	data.Any1 = 123
	for i := 0; i < 12; i++ {
		data.Map3[i] = i
		data.Map3[fmt.Sprintf("key-%d", i)] = i
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
