package fmtx

import (
	"fmt"
	"reflect"
	"strings"
)

var Options options = options{
	MaxDepth:             3,
	MaxArray:             50,
	MaxPropertyBreakLine: 10,
	ShowStructMethod:     true,
	ColorMap: colorMap{
		Int:      [2]string{"34", "39;49"},
		Float:    [2]string{"36", "39;49"},
		String:   [2]string{"32", "39;49"},
		Bool:     [2]string{"33", "39;49"},
		Ptr:      [2]string{"33", "39;49"},
		Property: [2]string{"90", "39;49"},
		Func:     [2]string{"3;36", "39;49;23"},
		Chan:     [2]string{"31", "39;49"},
		Nil:      [2]string{"93", "39;49"},
		Tip:      [2]string{"2", "22"},
	},
}

type options struct {
	MaxDepth             uint
	MaxArray             int
	MaxPropertyBreakLine int
	ShowStructMethod     bool
	ColorMap             colorMap
}

type colorMap struct {
	Int      [2]string
	Float    [2]string
	Bool     [2]string
	String   [2]string
	Ptr      [2]string
	Property [2]string
	Func     [2]string
	Chan     [2]string
	Nil      [2]string
	Tip      [2]string
}

func Println(a ...any) (n int, err error) {
	arr := make([]any, len(a))
	for i, v := range a {
		arr[i] = String(v)
	}
	return fmt.Println(arr...)
}

func String(o any) string {
	v := reflect.ValueOf(o)
	p := getPP()
	defer p.free()
	stringify(p, v, &Options, false, true, 0, nil)
	return string(p.buf)
}

func stringify(p *pp, v reflect.Value, opt *options, escapeString bool, showAliasName bool, level uint, parent *reflect.Value) {
	colors := opt.ColorMap
	switch v.Kind() {
	case reflect.Invalid:
		// var initAny any or var initInter error
		p.buf.WriteString(nilVal())
		p.buf.WriteChar('.')
		p.buf.WriteString(color("(<invalid>)", colors.Tip[0], colors.Tip[1]))
		return
	case reflect.Ptr:
		if v.IsNil() {
			p.buf.WriteString(nilVal())
			p.buf.WriteChar('.')
			p.buf.WriteString(color("", colors.Tip[0], ""))
			p.buf.WriteChar('(')
			getType(p, v.Type())
			p.buf.WriteChar(')')
			p.buf.WriteString(color("", "", colors.Tip[1]))
			return
		}
		p.buf.WriteString(color("&", colors.Ptr[0], colors.Ptr[1]))
		stringify(p, v.Elem(), opt, true, showAliasName, level, &v)
		return
	case reflect.String:
		t := v.Type()
		val := v.String()
		showType := showAliasName && t.Name() != t.Kind().String()
		if escapeString || showType {
			val = fmt.Sprintf("%q", val)
		}
		if level > 0 || showType || (parent != nil && parent.Kind() == reflect.Ptr) {
			val = color(val, colors.String[0], colors.String[1])
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
		if showAliasName && t.Name() != t.Kind().String() {
			getType(p, t)
			p.buf.WriteChar('(')
			p.buf.WriteString(val)
			p.buf.WriteChar(')')
			return
		}
		p.buf.WriteString(val)
	case reflect.Complex64, reflect.Complex128:
		p.buf.WriteString(color(fmt.Sprintf("%v", v), colors.Float[0], colors.Float[1]))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		val := ""
		if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
			val = color(fmt.Sprintf("%v", v), colors.Float[0], colors.Float[1])
		} else {
			val = color(fmt.Sprintf("%v", v), colors.Int[0], colors.Int[1])
		}
		if showAliasName && v.Type().Name() != v.Type().Kind().String() {
			getType(p, v.Type())
			p.buf.WriteChar('(')
			p.buf.WriteString(val)
			p.buf.WriteChar(')')
			return
		}
		p.buf.WriteString(val)
	case reflect.Interface:
		stringify(p, v.Elem(), opt, escapeString, showAliasName, level, nil)
	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice && v.IsNil() {
			p.buf.WriteString(nilVal())
			p.buf.WriteString(color("", colors.Tip[0], ""))
			p.buf.WriteString(".(")
			getType(p, v.Type())
			p.buf.WriteString(")")
			p.buf.WriteString(color("", "", colors.Tip[1]))
			return
		}
		i := len(p.buf)
		getType(p, v.Type())
		if v.Kind() == reflect.Slice {
			p.buf.Splice(i+1, fmt.Sprintf("%d/%d", v.Len(), v.Cap()))
		}
		if level >= opt.MaxDepth {
			p.buf.WriteString("{…}")
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
			stringify(p, v.Index(i), opt, true, false, level+1, nil)
		}
		if hasMore {
			p.buf.WriteString(", …")
		}
		p.buf.WriteChar(']')
	case reflect.Map:
		if v.IsNil() {
			p.buf.WriteString(nilVal())
			p.buf.WriteChar('.')
			p.buf.WriteString(color("", colors.Tip[0], ""))
			p.buf.WriteChar('(')
			getType(p, v.Type())
			p.buf.WriteChar(')')
			p.buf.WriteString(color("", "", colors.Tip[1]))
			return
		}
		getType(p, v.Type())
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
			stringify(p, k, opt, true, false, level+1, nil)
			p.buf.WriteString(": ")
			stringify(p, v.MapIndex(k), opt, true, false, level+1, nil)
		}
		if needBreak {
			p.buf.WriteChar('\n')
			p.buf.WriteString(indent)
		}
		p.buf.WriteChar('}')
	case reflect.Struct:
		getType(p, v.Type())
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
			stringify(p, v.Field(i), opt, true, false, level+1, nil)
		}

		if opt.ShowStructMethod {
			typ := val.Type()
			for i := 0; i < numMethod; i++ {
				m := val.Method(i)
				fname := typ.Method(i).Name

				if i > 0 || numMethod > 0 {
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
				stringify(p, m, opt, true, false, level+1, nil)
			}
			if needBreak {
				p.buf.WriteChar('\n')
				p.buf.WriteString(indent)
			}
		}
		p.buf.WriteChar('}')
	case reflect.Chan:
		if v.IsNil() {
			p.buf.WriteString(nilVal())
			p.buf.WriteChar('.')
			p.buf.WriteString(color("", colors.Tip[0], ""))
			p.buf.WriteChar('(')
			getType(p, v.Type())
			p.buf.WriteChar(')')
			p.buf.WriteString(color("", "", colors.Tip[1]))
			return
		}

		p.buf.WriteString(color("", colors.Chan[0], ""))
		getType(p, v.Type())
		if v.Cap() > 0 {
			i := len(p.buf)
			p.buf.Remove(i-1, 1)
			p.buf.WriteString(fmt.Sprintf(",%d/%d)", v.Len(), v.Cap()))
		}
		p.buf.WriteString(color("", "", colors.Chan[1]))
	case reflect.Func:
		if v.IsNil() {
			p.buf.WriteString(nilVal())
			p.buf.WriteChar('.')
			p.buf.WriteString(color("", colors.Tip[0], ""))
			p.buf.WriteChar('(')
			p.buf.WriteString("func(){}")
			p.buf.WriteChar(')')
			p.buf.WriteString(color("", "", colors.Tip[1]))
			return
		}
		p.buf.WriteString(color("", colors.Func[0], ""))
		getType(p, v.Type())
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
		p.buf.WriteChar(',')
		getType(p, t.Elem())
		p.buf.WriteChar('>')
	case reflect.Array:
		p.buf.WriteChar('[')
		p.buf.WriteString(fmt.Sprintf("%d", t.Len()))
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
		p.buf.WriteString("[func(")
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
				p.buf.WriteChar(')')
				p.buf.WriteChar(')')
			}
			p.buf.WriteChar(' ')
		}
		p.buf.WriteString("{}]")
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

func nilVal() string {
	return color("nil", Options.ColorMap.Nil[0], Options.ColorMap.Nil[1])
}
