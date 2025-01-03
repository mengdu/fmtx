# fmtmx

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
