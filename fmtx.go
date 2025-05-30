package fmtx

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Options struct {
	// print max depth of struct, map, slice, array
	MaxDepth uint
	// print max array length
	MaxArray int
	// max property triggering line breaks
	MaxPropertyBreakLine int
	// show struct method
	ShowStructMethod bool
	// show type name
	ShowTypeName bool
	// default color map
	ColorMap ColorMap
	// coustom color function
	Color func(s string, start string, end string) string
	// custom prints
	Prints []func(v any, opt Options) (string, bool)
}

type ColorMap struct {
	Int      [2]string
	Float    [2]string
	Complex  [2]string
	Bool     [2]string
	String   [2]string
	Ptr      [2]string
	Property [2]string
	FuncTag  [2]string
	Func     [2]string
	Chan     [2]string
	Nil      [2]string
	Tip      [2]string
}

var Default Options = Options{
	MaxDepth:             3,
	MaxArray:             25,
	MaxPropertyBreakLine: 10,
	ColorMap: ColorMap{
		Int:      [2]string{"34", "39"},
		Float:    [2]string{"36", "39"},
		Complex:  [2]string{"35", "39"},
		String:   [2]string{"32", "39"},
		Bool:     [2]string{"33", "39"},
		Ptr:      [2]string{"33", "39"},
		Property: [2]string{"39", "39"},
		FuncTag:  [2]string{"3;33", "39;23"},
		Func:     [2]string{"3;36", "39;23"},
		Chan:     [2]string{"31", "39"},
		Nil:      [2]string{"93", "39"},
		Tip:      [2]string{"2", "22"},
	},
	Color: Color,
}

var DefaultWriter = os.Stdout

func Println(a ...any) (n int, err error) {
	if DefaultWriter == nil {
		DefaultWriter = os.Stdout
	}
	l := len(a)
	for i, v := range a {
		n, err = DefaultWriter.WriteString(String(v))
		if err != nil {
			return
		}
		if i < l-1 {
			n, err = DefaultWriter.WriteString(" ")
			if err != nil {
				return
			}
		}
	}
	n, err = DefaultWriter.WriteString("\n")
	return
}

func String(o any) string {
	v := reflect.ValueOf(o)
	p := getPP()
	defer p.free()
	stringify(p, v, &Default, false, 0, nil)
	return string(p.buf)
}

// create a custom options function
func New(opt *Options) func(v any) string {
	return func(o any) string {
		v := reflect.ValueOf(o)
		p := getPP()
		defer p.free()
		stringify(p, v, opt, false, 0, nil)
		return string(p.buf)
	}
}

func stringify(p *pp, v reflect.Value, opt *Options, escapeString bool, level uint, parent *reflect.Value) {
	color := opt.Color
	colors := opt.ColorMap
	kind := v.Kind()

	if len(opt.Prints) > 0 && kind != reflect.Invalid && v.CanInterface() {
		for _, printer := range opt.Prints {
			if s, ok := printer(v.Interface(), *opt); ok {
				p.buf.WriteString(s)
				return
			}
		}
	}

	switch kind {
	case reflect.Invalid:
		// var initAny any or var initInter error
		p.buf.WriteString(nilVal(opt))
		if opt.ShowTypeName {
			p.buf.WriteChar('.')
			p.buf.WriteString(color("(<invalid>)", colors.Tip[0], colors.Tip[1]))
		}
		return
	case reflect.Ptr:
		if v.IsNil() {
			p.buf.WriteString(nilVal(opt))
			if opt.ShowTypeName {
				p.buf.WriteChar('.')
				p.buf.WriteString(color("", colors.Tip[0], ""))
				p.buf.WriteChar('(')
				getType(p, v.Type())
				p.buf.WriteChar(')')
				p.buf.WriteString(color("", "", colors.Tip[1]))
			}
			return
		}
		if opt.ShowTypeName {
			p.buf.WriteString(color("&", colors.Ptr[0], colors.Ptr[1]))
		}
		stringify(p, v.Elem(), opt, level > 0, level, &v)
		return
	case reflect.String:
		t := v.Type()
		val := v.String()
		showType := opt.ShowTypeName && t.Name() != t.Kind().String()
		if escapeString || showType {
			val = color(strconv.Quote(val), colors.String[0], colors.String[1])
		}
		if showType {
			getType(p, t)
			p.buf.WriteChar('(')
			p.buf.WriteString(val)
			p.buf.WriteChar(')')
			return
		}
		p.buf.WriteString(val)
	case reflect.Bool:
		t := v.Type()
		val := ""
		if v.Bool() {
			val = color("true", colors.Bool[0], colors.Bool[1])
		} else {
			val = color("false", colors.Bool[0], colors.Bool[1])
		}
		if opt.ShowTypeName && t.Name() != t.Kind().String() {
			getType(p, t)
			p.buf.WriteChar('(')
			p.buf.WriteString(val)
			p.buf.WriteChar(')')
			return
		}
		p.buf.WriteString(val)
	case reflect.Complex64, reflect.Complex128:
		val := v.Complex()
		p.buf.WriteString(color("", colors.Complex[0], ""))
		p.buf.WriteString(strconv.FormatFloat(real(val), 'f', -1, 64))
		p.buf.WriteString(strconv.FormatFloat(imag(val), 'f', -1, 64))
		p.buf.WriteString(color("", "", colors.Complex[1]))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		val := ""
		if kind == reflect.Float32 {
			val = color(strconv.FormatFloat(v.Float(), 'g', -1, 32), colors.Float[0], colors.Float[1])
		} else if kind == reflect.Float64 {
			val = color(strconv.FormatFloat(v.Float(), 'g', -1, 64), colors.Float[0], colors.Float[1])
		} else {
			if kind == reflect.Uint || kind == reflect.Uint8 ||
				kind == reflect.Uint16 || kind == reflect.Uint32 ||
				kind == reflect.Uint64 {
				val = color(strconv.FormatUint(v.Uint(), 10), colors.Int[0], colors.Int[1])
			} else {
				val = color(strconv.FormatInt(v.Int(), 10), colors.Int[0], colors.Int[1])
			}
		}
		t := v.Type()
		if opt.ShowTypeName && t.Name() != t.Kind().String() {
			getType(p, t)
			p.buf.WriteChar('(')
			p.buf.WriteString(val)
			p.buf.WriteChar(')')
			return
		}
		p.buf.WriteString(val)
	case reflect.Interface:
		stringify(p, v.Elem(), opt, escapeString, level, nil)
	case reflect.Slice, reflect.Array:
		if kind == reflect.Slice && v.IsNil() {
			p.buf.WriteString(nilVal(opt))
			if opt.ShowTypeName {
				p.buf.WriteString(color("", colors.Tip[0], ""))
				p.buf.WriteString(".(")
				getType(p, v.Type())
				p.buf.WriteString(")")
				p.buf.WriteString(color("", "", colors.Tip[1]))
			}
			return
		}

		if opt.ShowTypeName {
			i := len(p.buf)
			getType(p, v.Type())
			if kind == reflect.Slice {
				l := strconv.FormatInt(int64(v.Len()), 10)
				p.buf.Splice(i+1, l)
				i = len(l) + i
				p.buf.Splice(i+1, "/")
				p.buf.Splice(i+2, strconv.FormatInt(int64(v.Cap()), 10))
			}
		}

		if level >= opt.MaxDepth {
			p.buf.WriteString("[…]")
			return
		}
		n := v.Len()
		hasMore := v.Len() > opt.MaxArray
		if n > opt.MaxArray {
			n = opt.MaxArray
		}

		p.buf.WriteChar('[')
		for i := 0; i < n; i++ {
			if i > 0 {
				p.buf.WriteString(", ")
			}
			stringify(p, v.Index(i), opt, true, level+1, nil)
		}
		if hasMore {
			p.buf.WriteString(", …")
		}
		p.buf.WriteChar(']')
	case reflect.Map:
		if v.IsNil() {
			p.buf.WriteString(nilVal(opt))
			if opt.ShowTypeName {
				p.buf.WriteChar('.')
				p.buf.WriteString(color("", colors.Tip[0], ""))
				p.buf.WriteChar('(')
				getType(p, v.Type())
				p.buf.WriteChar(')')
				p.buf.WriteString(color("", "", colors.Tip[1]))
			}
			return
		}
		if opt.ShowTypeName {
			getType(p, v.Type())
		}
		if level >= opt.MaxDepth {
			p.buf.WriteString("{…}")
			return
		}
		needBreak := v.Len() > opt.MaxPropertyBreakLine
		indent := strings.Repeat(" ", int(level*2))
		p.buf.WriteChar('{')
		if needBreak {
			p.buf.WriteChar('\n')
			p.buf.WriteString(indent)
			p.buf.WriteString("  ")
		}

		for i, k := range v.MapKeys() {
			if i > 0 {
				if needBreak {
					p.buf.WriteString(",\n")
					p.buf.WriteString(indent)
					p.buf.WriteString("  ")
				} else {
					p.buf.WriteString(", ")
				}
			}
			stringify(p, k, opt, true, level+1, nil)
			p.buf.WriteString(": ")
			stringify(p, v.MapIndex(k), opt, true, level+1, nil)
		}
		if needBreak {
			p.buf.WriteChar('\n')
			p.buf.WriteString(indent)
		}
		p.buf.WriteChar('}')
	case reflect.Struct:
		if opt.ShowTypeName {
			getType(p, v.Type())
		}
		if level >= opt.MaxDepth {
			p.buf.WriteString("{…}")
			return
		}

		val := v
		numMethod := 0
		if opt.ShowStructMethod {
			if v.CanAddr() {
				val = val.Addr()
			} else if parent != nil {
				val = *parent
			}
			numMethod = val.NumMethod()
		}
		numField := v.NumField()
		needBreak := (numField + numMethod) > opt.MaxPropertyBreakLine

		indent := strings.Repeat(" ", int(level*2))
		p.buf.WriteChar('{')
		if needBreak {
			p.buf.WriteChar('\n')
			p.buf.WriteString(indent)
			p.buf.WriteString("  ")
		}

		for i := 0; i < numField; i++ {
			f := v.Type().Field(i)
			if f.Type.Kind() == reflect.Func && !opt.ShowStructMethod {
				continue
			}
			if i > 0 {
				if needBreak {
					p.buf.WriteString(",\n")
					p.buf.WriteString(indent)
					p.buf.WriteString("  ")
				} else {
					p.buf.WriteString(", ")
				}
			}
			p.buf.WriteString(color(f.Name, colors.Property[0], colors.Property[1]))
			p.buf.WriteString(": ")
			stringify(p, v.Field(i), opt, true, level+1, nil)
		}

		if opt.ShowStructMethod {
			typ := val.Type()
			for i := 0; i < numMethod; i++ {
				m := val.Method(i)
				fname := typ.Method(i).Name
				if i > 0 || numField > 0 {
					if needBreak {
						p.buf.WriteString(",\n")
						p.buf.WriteString(indent)
						p.buf.WriteString("  ")
					} else {
						p.buf.WriteString(", ")
					}
				}

				p.buf.WriteString(color(fname, colors.Property[0], colors.Property[1]))
				p.buf.WriteString(": ")
				stringify(p, m, opt, true, level+1, nil)
			}
			if needBreak {
				p.buf.WriteChar('\n')
				p.buf.WriteString(indent)
			}
		}
		p.buf.WriteChar('}')
	case reflect.Chan:
		if v.IsNil() {
			p.buf.WriteString(nilVal(opt))
			if opt.ShowTypeName {
				p.buf.WriteChar('.')
				p.buf.WriteString(color("", colors.Tip[0], ""))
				p.buf.WriteChar('(')
				getType(p, v.Type())
				p.buf.WriteChar(')')
				p.buf.WriteString(color("", "", colors.Tip[1]))
			}
			return
		}

		p.buf.WriteString(color("", colors.Chan[0], ""))
		getType(p, v.Type())
		if v.Cap() > 0 {
			i := len(p.buf)
			l := strconv.FormatInt(int64(v.Len()), 10)
			p.buf.Splice(i-1, ",")
			p.buf.Splice(i, l)
			i = len(l) + i
			p.buf.Splice(i, "/")
			p.buf.Splice(i+1, strconv.FormatInt(int64(v.Cap()), 10))
		}
		p.buf.WriteString(color("", "", colors.Chan[1]))
	case reflect.Func:
		if v.IsNil() {
			p.buf.WriteString(nilVal(opt))
			if opt.ShowTypeName {
				p.buf.WriteChar('.')
				p.buf.WriteString(color("", colors.Tip[0], ""))
				p.buf.WriteChar('(')
				getType(p, v.Type())
				p.buf.WriteChar(')')
				p.buf.WriteString(color("", "", colors.Tip[1]))
			}
			return
		}
		fname := ""
		fn := runtime.FuncForPC(v.Pointer())
		if fn != nil {
			fname = fn.Name()
			if fname == "reflect.methodValueCall" {
				fname = ""
			}

			if fname != "" && !opt.ShowTypeName {
				if strings.HasSuffix(fname, "[...]") {
					arr := strings.Split(fname[0:len(fname)-5], ".")
					fname = arr[len(arr)-1] + "[...]"
				} else if strings.HasSuffix(fname, "-fm") {
					arr := strings.Split(fname[0:len(fname)-3], ".")
					fname = arr[len(arr)-1]
				} else {
					arr := strings.Split(fname, ".")
					if len(arr) > 1 {
						fname = arr[len(arr)-1]
					}
				}
			}
		}
		t := v.Type()
		p.buf.WriteString(color("f ", colors.FuncTag[0], colors.FuncTag[1]))
		p.buf.WriteString(color("", colors.Func[0], ""))
		if fname != "" {
			p.buf.WriteString(fname)
		}
		p.buf.WriteString("(")
		inLen := t.NumIn()
		for i := 0; i < inLen; i++ {
			if i == inLen-1 && t.IsVariadic() {
				l := len(p.buf)
				p.buf.WriteString(", ...")
				getType(p, t.In(i))
				p.buf.Remove(l, 2) // remove "[]"
			} else {
				if i > 0 {
					p.buf.WriteString(", ")
				}
				getType(p, t.In(i))
			}
		}
		p.buf.WriteChar(')')

		outLen := t.NumOut()
		if outLen > 0 {
			p.buf.WriteChar(' ')
			if outLen > 1 {
				p.buf.WriteChar('(')
			}
			for i := 0; i < outLen; i++ {
				if i > 0 {
					p.buf.WriteString(", ")
				}
				getType(p, t.Out(i))
			}
			if outLen > 1 {
				p.buf.WriteChar(')')
			}
			p.buf.WriteChar(' ')
		}
		p.buf.WriteString("{}")
		p.buf.WriteString(color("", "", colors.Func[1]))
		return
	default:
		p.buf.WriteString(fmt.Sprintf("%v", v))
	}
}

func getType(p *pp, t reflect.Type) {
	switch t.Kind() {
	case reflect.Ptr:
		p.buf.WriteChar('*')
		getType(p, t.Elem())
	case reflect.Map:
		p.buf.WriteString("map<")
		getType(p, t.Key())
		p.buf.WriteString(", ")
		getType(p, t.Elem())
		p.buf.WriteChar('>')
	case reflect.Array:
		p.buf.WriteChar('[')
		p.buf.WriteString(strconv.FormatInt(int64(t.Len()), 10))
		p.buf.WriteChar(']')
		getType(p, t.Elem())
	case reflect.Slice:
		p.buf.WriteString("[]")
		getType(p, t.Elem())
	case reflect.Struct:
		// Anonymous struct
		if t.Name() == "" {
			return
		}
		p.buf.WriteString(t.String())
	case reflect.Chan:
		p.buf.WriteString("chan")
		if t.ChanDir() == reflect.RecvDir {
			p.buf.WriteString("->")
		} else if t.ChanDir() == reflect.SendDir {
			p.buf.WriteString("<-")
		}
		p.buf.WriteChar('(')
		getType(p, t.Elem())
		p.buf.WriteChar(')')
	case reflect.Func:
		p.buf.WriteString("func (")
		inLen := t.NumIn()
		for i := 0; i < inLen; i++ {
			if i == inLen-1 && t.IsVariadic() {
				l := len(p.buf)
				p.buf.WriteString(", ...")
				getType(p, t.In(i))
				p.buf.Remove(l, 2) // remove "[]"
			} else {
				if i > 0 {
					p.buf.WriteString(", ")
				}
				getType(p, t.In(i))
			}
		}
		p.buf.WriteChar(')')

		outLen := t.NumOut()
		if outLen > 0 {
			p.buf.WriteChar(' ')
			if outLen > 1 {
				p.buf.WriteChar('(')
			}
			for i := 0; i < outLen; i++ {
				if i > 0 {
					p.buf.WriteString(", ")
				}
				getType(p, t.Out(i))
			}
			if outLen > 1 {
				p.buf.WriteChar(')')
			}
			p.buf.WriteChar(' ')
		}
		p.buf.WriteString("{}")
	case reflect.Interface:
		if t.Name() != "" {
			if t.PkgPath() != "" {
				p.buf.WriteString(t.PkgPath())
				p.buf.WriteChar('.')
				p.buf.WriteString(t.Name())
				return
			}
			p.buf.WriteString(t.Name())
			return
		}
		p.buf.WriteString("any")
		return
	default:
		p.buf.WriteString(t.String())
	}
}

func nilVal(opt *Options) string {
	return opt.Color("nil", opt.ColorMap.Nil[0], opt.ColorMap.Nil[1])
}
