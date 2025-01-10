# fmtx

Better Formatting and Printing in Golang

```sh
go get github.com/mengdu/fmtx
```

```go
func main() {
  // fmtx.EnableColor = false
  // fmtx.Options.MaxDepth = 5

  fmtx.Println(123, 3.14, ptr(124), MyInt(123), ptr(MyInt(123)))
  fmtx.Println([]int{1, 2, 3, -1, -2})
  fmtx.Println([]string{"a", "b", "c", "Hello \n \"world\"."})
}
```

![](preview.png)

## Benchmark

```sh
go test -bench=. -run=^$ -benchmem -benchtime=5s
```

```
goos: darwin
goarch: amd64
pkg: github.com/mengdu/fmtx
cpu: Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz
BenchmarkString-16                        212820             27961 ns/op            9429 B/op        278 allocs/op
BenchmarkStringDisableColor-16            357429             16721 ns/op            4998 B/op        151 allocs/op
BenchmarkSprintf-16                      1760349              3470 ns/op            1200 B/op         34 allocs/op
PASS
ok      github.com/mengdu/fmtx  22.277s
```
