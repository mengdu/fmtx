package fmtx

import (
	"fmt"
	"testing"
	"time"
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

func BenchmarkFmtxString(b *testing.B) {
	SetEnableColor(true)
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = String(data.Slice3)
	}
}

func BenchmarkFmtxStringDisableColor(b *testing.B) {
	SetEnableColor(false)
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = String(data.Slice3)
	}
}

func BenchmarkSprintf1(b *testing.B) {
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%v", data.Slice3)
	}
}

func BenchmarkSprintf2(b *testing.B) {
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%+v", data.Slice3)
	}
}

func BenchmarkSprintf3(b *testing.B) {
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%#v", data.Slice3)
	}
}

func BenchmarkFmtxStringBig(b *testing.B) {
	SetEnableColor(true)
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = String(data)
	}
}

func BenchmarkFmtxStringDisableColorBig(b *testing.B) {
	SetEnableColor(false)
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = String(data)
	}
}

func BenchmarkSprintfBig1(b *testing.B) {
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%v", data)
	}
}

func BenchmarkSprintfBig2(b *testing.B) {
	data := genData()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%+v", data)
	}
}

func BenchmarkSprintfBig3(b *testing.B) {
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
		_ = Color("%#v\n", "35", "39")
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
