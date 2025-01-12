# fmtx

More intuitive printing values for Golang.

```sh
go get github.com/mengdu/fmtx
```

```go
type Foo struct {
  Key  int
  Loop *Foo
}

func main() {
  fmtx.Println([]int{1, 2, 3, 4}) // [4/4]int[1, 2, 3, 4]
  fmtx.Println(map[any]any{"a": 1, "b": 3.14, 1: true, "c": "Hello, \nworld."}) // map<any,any>{"a": 1, "b": 3.14, 1: true, "c": "Hello, \nworld."}
  fmtx.Println(Foo{}) // main.Foo{Key: 0, Loop: nil.(*main.Foo)}
}
```

![](preview.png)

## Benchmark

```sh
go test -bench=. -run=^$ -benchmem -benchtime=5s
```

```txt
goos: darwin
goarch: amd64
pkg: github.com/mengdu/fmtx
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkFmtxString-12                           1951065              3022 ns/op            1112 B/op         58 allocs/op
BenchmarkFmtxStringDisableColor-12               3332950              1812 ns/op             360 B/op         18 allocs/op
BenchmarkSprintf1-12                             3172249              1800 ns/op             272 B/op         10 allocs/op
BenchmarkSprintf2-12                             3256015              1987 ns/op             272 B/op         10 allocs/op
BenchmarkSprintf3-12                             2901669              2073 ns/op             336 B/op         10 allocs/op
BenchmarkFmtxStringBig-12                          97717             66711 ns/op           24240 B/op       1308 allocs/op
BenchmarkFmtxStringDisableColorBig-12             179902             30867 ns/op            7506 B/op        326 allocs/op
BenchmarkSprintfBig1-12                            31928            212421 ns/op           83876 B/op       2149 allocs/op
BenchmarkSprintfBig2-12                            27352            197811 ns/op           85117 B/op       2205 allocs/op
BenchmarkSprintfBig3-12                            27866            206434 ns/op           89727 B/op       2205 allocs/op
PASS
ok      github.com/mengdu/fmtx  77.373s
```
