package console

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	moveCursorOneLineUpInTheTerminal = "\033[A"
	moveCursorToStartOfLine          = "\r"
	moveCursorToEndOfLine            = "\033[999C"
	clearConsole                     = "\033[H\033[2J"
	clearLine                        = "\033[2K"
)

var lastPrintWasASpinner = false
var spinnerCount int = 0 //Overflow is fine
// var indicators = []string{"-", "\\", "|", "/"}
var indicators = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

func Spinner(arguments ...any) {
	if !ShouldLog(INFO) {
		return
	}
	spinnerCount++
	controlCharacters := ""
	if lastPrintWasASpinner {
		controlCharacters = moveCursorOneLineUpInTheTerminal + moveCursorToStartOfLine + clearLine
	}
	nextIndicator := indicators[spinnerCount%len(indicators)]

	fmt.Print(controlCharacters)
	fmt.Print(nextIndicator, " ")
	common(commonArguments{Types: false}, arguments...)
	fmt.Print(moveCursorOneLineUpInTheTerminal, moveCursorToEndOfLine, nextIndicator)
	lastPrintWasASpinner = true
	fmt.Print("\r\n")
}

func PrettyMemoryAllocation(memoryInBytes uint64) string {
	if memoryInBytes >= 1<<20 { // 1 MB or more
		return fmt.Sprintf("%.2f MB", float64(memoryInBytes)/(1<<20))
	} else if memoryInBytes >= 1<<10 { // 1 KB or more
		return fmt.Sprintf("%.2f KB", float64(memoryInBytes)/(1<<10))
	} else {
		return fmt.Sprintf("%d bytes", memoryInBytes)
	}
}

type commonArguments struct {
	Types bool
}

func Clear() {
	if ShouldLog(INFO) {
		fmt.Print(clearConsole)
	}
}

type GoIsDumb struct{}

func (GoIsDumb) Log(arguments ...any) {
	lastPrintWasASpinner = false
	common(commonArguments{Types: true}, arguments...)
}

var Types = GoIsDumb{}

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
)

var logLevel = INFO

func SetLogLevel(ll LogLevel) {
	logLevel = ll
}

func ShouldLog(ll LogLevel) bool {
	return ll >= logLevel
}

func Copy(src io.Reader) {
	if ShouldLog(DEBUG) {
		io.Copy(os.Stdout, src)
	} else {
		//TODO: Wtf, images are not pulled without doing this?
		io.Copy(io.Discard, src)
	}
}

func Log(arguments ...any) {
	Info(arguments...)
}

func Info(arguments ...any) {
	if ShouldLog(INFO) {
		lastPrintWasASpinner = false
		common(commonArguments{Types: false}, arguments...)
	}
}

func Debug(arguments ...any) {
	if ShouldLog(DEBUG) {
		lastPrintWasASpinner = false
		common(commonArguments{Types: false}, arguments...)
	}
}

func Trace(arguments ...any) {
	if ShouldLog(TRACE) {
		lastPrintWasASpinner = false
		common(commonArguments{Types: false}, arguments...)
	}
}

func Error(arguments ...any) {
	if ShouldLog(ERROR) {
		lastPrintWasASpinner = false
		common(commonArguments{Types: false}, arguments...)
	}
}

func Fatal(arguments ...any) {
	if ShouldLog(ERROR) {
		lastPrintWasASpinner = false
		common(commonArguments{Types: false}, arguments...)
		panic("Fatal")
	}
}

func common(settings commonArguments, arguments ...any) {
	var builder strings.Builder
	if len(arguments) == 1 {
		examiner(&settings, &builder, 0, reflect.ValueOf(arguments[0]), reflect.ValueOf(arguments[0]).Kind() == reflect.String)
	} else {
		examiner(&settings, &builder, 0, reflect.ValueOf(arguments), true)
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

func examiner(settings *commonArguments, str *strings.Builder, depth int, v reflect.Value, isVarargs bool) {
	switch v.Kind() {
	case reflect.Array:
		str.WriteString("[ ")
		for i := range v.Len() {
			examiner(settings, str, depth+1, v.Index(i).Elem(), false)
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
			examiner(settings, str, depth+1, v.Index(i), isVarargs && depth == 0)
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
		//Go is so butiful, not
		errorInterface := reflect.TypeOf((*error)(nil)).Elem()
		if v.Type().Implements(errorInterface) {
			str.WriteString(fmt.Sprint(v))
		}
		examiner(settings, str, depth, v.Elem(), isVarargs)
	case reflect.Chan:
	case reflect.Map:
		str.WriteRune('{')
		for iter := v.MapRange(); iter.Next(); {
			key := iter.Key()
			value := iter.Value()
			examiner(settings, str, depth+1, key, false)
			str.WriteString(": ")
			examiner(settings, str, depth+1, value, false)
		}
		str.WriteRune('}')
	case reflect.Ptr:
		examiner(settings, str, depth, v.Elem(), false)
		break
	case reflect.Struct:
		str.WriteString("{ ")
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			str.WriteString(merge(f.Name, ": "))
			examiner(settings, str, depth+1, v.Field(i), false)
			if i+1 < v.NumField() {
				str.WriteString(", ")
			}
		}
		str.WriteString(" }")
	case reflect.Invalid:
		break
		//fmt.Println()
		//fmt.Println(v)
		//fmt.Println(v.Type())
		//fmt.Println(v.Elem())
		//panic("Oh no, invalid type")
	case reflect.String:
		if isVarargs {
			str.WriteString(v.String())
		} else {
			str.WriteString(merge("\"", v.String(), "\""))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str.WriteString(strconv.Itoa(int(v.Int())))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
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
