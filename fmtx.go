package fmtx

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
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
	case reflect.Ptr:
		return color("&", colors.Ptr[0], colors.Ptr[1]) + stringify(v.Elem(), opt, true, showAliasName, level)
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
		t := v.Type()
		kind := t.Elem().Name()
		if t.Elem().Kind() == reflect.Interface {
			kind = "any"
		} else if t.Elem().Kind() == reflect.Ptr {
			kind = color("*", colors.Ptr[0], colors.Ptr[1]) + t.Elem().Elem().Name()
		}
		size := fmt.Sprintf("%d", v.Len())
		if v.Kind() == reflect.Slice {
			size = fmt.Sprintf("%d/%d", v.Len(), v.Cap())
		}
		typeName := fmt.Sprintf("%s(%s)", kind, size)
		if level >= opt.MaxDepth {
			return fmt.Sprintf("%s[…]", typeName)
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
		val := fmt.Sprintf("[%s]", strings.Join(arr, ", "))
		return fmt.Sprintf("%s%s", typeName, val)
	case reflect.Map:
		kt := v.Type().Key().Kind().String()
		if v.Type().Key().Kind() == reflect.Interface {
			kt = "any"
		} else if showAliasName && v.Type().Key().Name() != kt {
			kt = v.Type().Key().String()
		}
		vt := v.Type().Elem().Kind().String()
		if v.Type().Elem().Kind() == reflect.Interface {
			vt = "any"
		} else if showAliasName && v.Type().Elem().Name() != vt {
			vt = v.Type().Elem().String()
		}
		name := fmt.Sprintf("map<%s,%s>", kt, vt)
		if level >= opt.MaxDepth {
			return fmt.Sprintf("%s{…}", name)
		}
		fields := make([]string, v.Len())
		for i, k := range v.MapKeys() {
			key := stringify(k, opt, true, false, level+1)
			fields[i] = fmt.Sprintf("%s: %s", color(key, colors.Property[0], colors.Property[1]), stringify(v.MapIndex(k), opt, true, false, level+1))
		}
		body := fmt.Sprintf("{%s}", strings.Join(fields, ", "))
		if len(fields) < 20 && len(strip(body)) > 100 {
			indent := strings.Repeat(" ", int((level+1)*2))
			body = fmt.Sprintf("{\n%s\n%s}", indent+strings.Join(fields, ",\n"+indent), strings.Repeat(" ", int(level*2)))
		}
		return fmt.Sprintf("%s%s", name, body)
	case reflect.Struct:
		name := fmt.Sprintf("%s.%s", v.Type().PkgPath(), v.Type().Name())
		if level >= opt.MaxDepth {
			return fmt.Sprintf("%s{…}", name)
		}

		fields := make([]string, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			val := stringify(v.FieldByName(f.Name), opt, true, false, level+1)
			fields[i] = fmt.Sprintf("%s: %s", color(f.Name, colors.Property[0], colors.Property[1]), val)
		}
		for i := 0; i < v.NumMethod(); i++ {
			m := v.Method(i)
			fname := v.Type().Method(i).Name
			fields = append(fields, fmt.Sprintf("%s: %s", color(fname, colors.Func[0], colors.Func[1]), stringify(m, opt, true, false, level+1)))
		}
		body := fmt.Sprintf("{%s}", strings.Join(fields, ", "))
		if len(fields) < 20 && len(strip(body)) > 100 {
			indent := strings.Repeat(" ", int((level+1)*2))
			body = fmt.Sprintf("{\n%s\n%s}", indent+strings.Join(fields, ",\n"+indent), strings.Repeat(" ", int(level*2)))
		}
		return fmt.Sprintf("%s.%s%s", v.Type().PkgPath(), v.Type().Name(), body)
	case reflect.Chan:
		return "chan{}"
	case reflect.Func:
		pc := runtime.FuncForPC(v.Pointer())
		return color(fmt.Sprintf("[Func: %s()]", pc.Name()), colors.Func[0], colors.Func[1])
	default:
		return fmt.Sprintf("%v", v)
	}
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// removes all ANSI escape sequences from a string.
func strip(input string) string {
	return ansiRegex.ReplaceAllString(input, "")
}
