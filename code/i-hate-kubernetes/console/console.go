package console

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func Log(arguments ...any) {
	common(arguments...)
}

func Info(arguments ...any) {
	common(arguments...)
}

func Debug(arguments ...any) {
	common(arguments...)
}

func Trace(arguments ...any) {
	common(arguments...)
}

func Error(arguments ...any) {
	common(arguments...)
}

func common(arguments ...any) {
	var builder strings.Builder
	if len(arguments) == 1 {
		examiner(&builder, 0, reflect.ValueOf(arguments[0]), reflect.ValueOf(arguments[0]).Kind() == reflect.String)
	} else {
		examiner(&builder, 0, reflect.ValueOf(arguments), true)
	}
	fmt.Println(builder.String())
}

func merge(arguments ...string) string {
	var builder strings.Builder
	for _, arg := range arguments {
		builder.WriteString(arg)
	}
	return builder.String()
}

func log(arguments ...string) {
	fmt.Println(merge(arguments...))
}

func examiner(str *strings.Builder, depth int, v reflect.Value, isVarargs bool) {
	switch v.Kind() {
	case reflect.Array:
		str.WriteString("[ ")
		for i := range v.Len() {
			examiner(str, depth+1, v.Index(i).Elem(), false)
			if i+1 < v.Len() {
				str.WriteString(", ")
			}
		}
		str.WriteString(" ]")
	case reflect.Slice:
		if !isVarargs || depth != 0 {
			if v.Len() == 0 {
				str.WriteString("[]")
				return
			} else {
				str.WriteString("[ ")
			}
		}
		for i := range v.Len() {
			examiner(str, depth+1, v.Index(i), isVarargs && depth == 0)
			if i+1 < v.Len() && (!isVarargs) {
				str.WriteString(", ")
			} else if isVarargs {
				str.WriteString(" ")
			}
		}
		if !isVarargs || depth != 0 {
			str.WriteString(" ]")
		}
	case reflect.Interface:
		examiner(str, depth, v.Elem(), isVarargs)
	case reflect.Chan:
	case reflect.Map:
		str.WriteRune('{')
		for iter := v.MapRange(); iter.Next(); {
			key := iter.Key()
			value := iter.Value()
			examiner(str, depth+1, key, false)
			str.WriteString(": ")
			examiner(str, depth+1, value, false)
		}
		str.WriteRune('}')
	case reflect.Ptr:
		break
	case reflect.Struct:
		str.WriteString("{ ")
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			str.WriteString(merge(f.Name, ": "))
			examiner(str, depth+1, v.Field(i), false)
			if i+1 < v.NumField() {
				str.WriteString(", ")
			}
		}
		str.WriteString(" }")
	case reflect.Invalid:
		panic("Oh no, invalid type")
	case reflect.String:
		if isVarargs {
			str.WriteString(v.String())
		} else {
			str.WriteString(merge("\"", v.String(), "\""))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str.WriteString(strconv.Itoa(int(v.Int())))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str.WriteString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		str.WriteString(strconv.FormatFloat(v.Float(), 'f', -1, 64))
	case reflect.Bool:
		str.WriteString(strconv.FormatBool(v.Bool()))
	default:
		fmt.Println(v)
		log(v.Type().Kind().String())
		log(v.Kind().String())
		panic("Oh no, invalid type")
	}
}
