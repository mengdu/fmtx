package fmtx

import (
	"sync"
)

type buf []byte

func (b *buf) Write(s []byte) {
	*b = append(*b, s...)
}

func (b *buf) WriteChar(c byte) {
	*b = append(*b, c)
}

func (b *buf) WriteString(s string) {
	*b = append(*b, s...)
}

func (b *buf) Remove(i int, c int) {
	*b = append((*b)[:i], (*b)[i+c:]...)
}

func (b *buf) Splice(i int, s string) {
	*b = append((*b)[:i], append([]byte(s), (*b)[i:]...)...)
}

type pp struct {
	buf buf
}

func (p *pp) free() {
	// See https://golang.org/issue/23199.
	if cap(p.buf) > 64*1024 {
		p.buf = nil
	} else {
		p.buf = p.buf[:0]
	}
	pool.Put(p)
}

var pool = sync.Pool{
	New: func() any {
		return new(pp)
	},
}

func getPP() *pp {
	return pool.Get().(*pp)
}
