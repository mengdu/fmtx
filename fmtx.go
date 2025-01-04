package fmtx

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var Options options = options{
	MaxDepth: 3,
	MaxArray: 50,
	ColorMap: colorMap{
		Int:      [2]string{"34", "39;49"},
		Float:    [2]string{"36", "39;49"},
		String:   [2]string{"32", "39;49"},
		Bool:     [2]string{"33", "39;49"},
		Ptr:      [2]string{"33", "39;49"},
		Property: [2]string{"90", "39;49"},
		Func:     [2]string{"3;94", "39;49;23"},
		Chan:     [2]string{"31", "39;49"},
		Nil:      [2]string{"93", "39;49"},
	},
}

type options struct {
	MaxDepth uint
	MaxArray int
	ColorMap colorMap
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
	return stringify(v, Options, false, true, 0)
}

func stringify(v reflect.Value, opt options, escapeString bool, showAliasName bool, level uint) string {
	colors := opt.ColorMap
	switch v.Kind() {
	case reflect.Invalid:
		// var initAny any or var initInter error
		return fmt.Sprintf("<%s>", color("nil", colors.Nil[0], colors.Nil[1]))
	case reflect.Ptr:
		val := ""
		if v.IsNil() {
			val = fmt.Sprintf("%s<%s>", getType(v.Type()), color("nil", colors.Nil[0], colors.Nil[1]))
		} else {
			val = stringify(v.Elem(), opt, true, showAliasName, level)
		}
		return color("&", colors.Ptr[0], colors.Ptr[1]) + val
	case reflect.String:
		if escapeString {
			str := fmt.Sprintf("%q", v.String())
			return color(str, colors.String[0], colors.String[1])
		}
		return color(v.String(), colors.String[0], colors.String[1])
	case reflect.Bool:
		if v.Bool() {
			return color("true", colors.Bool[0], colors.Bool[1])
		} else {
			return color("false", colors.Bool[0], colors.Bool[1])
		}
	case reflect.Complex64, reflect.Complex128:
		return color(fmt.Sprintf("%v", v), colors.Float[0], colors.Float[1])
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		str := ""
		if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
			str = color(fmt.Sprintf("%v", v), colors.Float[0], colors.Float[1])
		} else {
			str = color(fmt.Sprintf("%v", v), colors.Int[0], colors.Int[1])
		}
		if showAliasName && v.Type().Name() != v.Type().Kind().String() {
			return fmt.Sprintf("%s(%s)", v.Type().String(), str)
		}
		return str
	case reflect.Interface:
		return stringify(v.Elem(), opt, escapeString, showAliasName, level)
	case reflect.Slice, reflect.Array:
		typ := getType(v.Type())
		size := ""
		if v.Kind() == reflect.Slice {
			len := color(fmt.Sprintf("%d", v.Len()), colors.Int[0], colors.Int[1])
			cap := color(fmt.Sprintf("%d", v.Cap()), colors.Int[0], colors.Int[1])
			size = fmt.Sprintf("%s/%s", len, cap)
		}
		typ = fmt.Sprintf("[%s%s", size, strings.TrimPrefix(typ, "["))
		if v.Kind() == reflect.Slice && v.IsNil() {
			return fmt.Sprintf("%s<%s>", typ, color("nil", colors.Nil[0], colors.Nil[1]))
		}
		if level >= opt.MaxDepth {
			return fmt.Sprintf("%s{…}", typ)
		}
		n := v.Len()
		hasMore := v.Len() > opt.MaxArray
		if n > opt.MaxArray {
			n = opt.MaxArray
		}
		arr := make([]string, n)
		for i := 0; i < n; i++ {
			arr[i] = stringify(v.Index(i), opt, true, false, level+1)
		}
		if hasMore {
			arr = append(arr, "…")
		}
		val := fmt.Sprintf("{%s}", strings.Join(arr, ", "))
		return fmt.Sprintf("%s%s", typ, val)
	case reflect.Map:
		typ := getType(v.Type())
		if v.IsNil() {
			return fmt.Sprintf("%s<%s>", typ, color("nil", colors.Nil[0], colors.Nil[1]))
		}
		if level >= opt.MaxDepth {
			return fmt.Sprintf("%s{…}", typ)
		}
		fields := make([]string, v.Len())
		for i, k := range v.MapKeys() {
			key := stringify(k, opt, true, false, level+1)
			fields[i] = fmt.Sprintf("%s: %s", color(key, colors.Property[0], colors.Property[1]), stringify(v.MapIndex(k), opt, true, false, level+1))
		}
		body := fmt.Sprintf("{%s}", strings.Join(fields, ", "))
		if len(fields) > 20 || len(strip(body)) > 100 {
			indent := strings.Repeat(" ", int((level+1)*2))
			body = fmt.Sprintf("{\n%s\n%s}", indent+strings.Join(fields, ",\n"+indent), strings.Repeat(" ", int(level*2)))
		}
		return fmt.Sprintf("%s%s", typ, body)
	case reflect.Struct:
		typ := getType(v.Type())
		if level >= opt.MaxDepth {
			return fmt.Sprintf("%s{…}", typ)
		}

		fields := make([]string, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			val := stringify(v.Field(i), opt, true, false, level+1)
			fields[i] = fmt.Sprintf("%s: %s", color(f.Name, colors.Property[0], colors.Property[1]), val)
		}
		for i := 0; i < v.NumMethod(); i++ {
			m := v.Method(i)
			fname := v.Type().Method(i).Name
			fields = append(fields, fmt.Sprintf("%s: %s", color(fname, colors.Property[0], colors.Property[1]), stringify(m, opt, true, false, level+1)))
		}
		body := fmt.Sprintf("{%s}", strings.Join(fields, ", "))
		if len(fields) > 20 || len(strip(body)) > 100 {
			indent := strings.Repeat(" ", int((level+1)*2))
			body = fmt.Sprintf("{\n%s\n%s}", indent+strings.Join(fields, ",\n"+indent), strings.Repeat(" ", int(level*2)))
		}
		return fmt.Sprintf("%s.%s%s", v.Type().PkgPath(), v.Type().Name(), body)
	case reflect.Chan:
		typ := getType(v.Type())
		if v.Cap() > 0 {
			typ = strings.TrimSuffix(typ, ")")
			len := color(fmt.Sprintf("%d", v.Len()), colors.Int[0], colors.Int[1])
			cap := color(fmt.Sprintf("%d", v.Cap()), colors.Int[0], colors.Int[1])
			return fmt.Sprintf("%s,%s/%s)", typ, len, cap)
		}
		if v.IsNil() {
			return fmt.Sprintf("%s<%s>", typ, color("nil", colors.Nil[0], colors.Nil[1]))
		}
		return typ
	case reflect.Func:
		typ := getType(v.Type())
		if v.IsNil() {
			return fmt.Sprintf("%s<%s>", typ, color("nil", colors.Nil[0], colors.Nil[1]))
		}
		return color(typ, colors.Func[0], colors.Func[1])
	default:
		return fmt.Sprintf("%v", v)
	}
}

func getType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr:
		return "*" + getType(t.Elem())
	case reflect.Map:
		key := getType(t.Key())
		val := getType(t.Elem())
		// return fmt.Sprintf("map[%s]%s", key, val)
		return fmt.Sprintf("map<%s,%s>", key, val)
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), getType(t.Elem()))
	case reflect.Slice:
		return fmt.Sprintf("[]%s", getType(t.Elem()))
	case reflect.Struct:
		return t.String()
	case reflect.Chan:
		dr := ""
		if t.ChanDir() == reflect.RecvDir {
			dr = "->"
		} else if t.ChanDir() == reflect.SendDir {
			dr = "<-"
		}
		return fmt.Sprintf("chan%s(%s)", dr, getType(t.Elem()))
	case reflect.Func:
		ni := t.NumIn()
		args := make([]string, ni)
		for i := 0; i < ni; i++ {
			typ := getType(t.In(i))
			if i == ni-1 && t.IsVariadic() {
				args[i] = fmt.Sprintf("...%s", strings.TrimPrefix(typ, "[]"))
			} else {
				args[i] = typ
			}
		}
		no := t.NumOut()
		outs := make([]string, no)
		for i := 0; i < no; i++ {
			outs[i] = getType(t.Out(i))
		}
		out := ""
		if len(outs) > 0 {
			out = fmt.Sprintf("(%s)", strings.Join(outs, ", "))
			out = " " + out
		}
		return fmt.Sprintf("func(%s)%s", strings.Join(args, ", "), out)
	case reflect.Interface:
		if t.Name() != "" {
			if t.PkgPath() != "" {
				return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
			}
			return t.Name()
		}
		return "any"
	default:
		if t.Name() != t.Kind().String() {
			return t.Name()
		}
		return t.String()
	}
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// removes all ANSI escape sequences from a string.
func strip(input string) string {
	return ansiRegex.ReplaceAllString(input, "")
}
