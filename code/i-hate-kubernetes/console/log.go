package console

import (
	"bytes"
	"io"
	"os"
	"sync"
	"time"
)

type LogFlag int
type LogFormat int

const (
	JSON LogFormat = iota
	TEXT
)
const (
	Ltime LogFlag = 1 << iota
	Ldate
	LlongFile
)

type Logger struct {
	destinations []*LogDestination
}

type LogLevel int

const (
	NOT_SPECIFIED LogLevel = iota
	TRACE
	DEBUG
	INFO
	WARNING
	ERROR
)

type LogDestination struct {
	output          io.Writer
	format          LogFormat
	flags           LogFlag
	minimumLogLevel LogLevel //The minimum loglevel required to log to this destination, inclusive
	maximumLogLevel LogLevel //The maximum loglevel to log to this destination, not inclusive
	outMutex        sync.Mutex
}

func StdLog() *Logger {
	return &Logger{
		destinations: []*LogDestination{
			{
				output: os.Stdout,
				format: TEXT,
				flags:  0,
			},
		},
	}
}

func openLogFile(path string) io.Writer {
	logDir := CreateLoggingDirectory()
	path = logDir + path
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		Fatal("Cannot open logfile: ", err)
	}
	return logFile
}

func NewLog() *Logger {
	return &Logger{
		destinations: []*LogDestination{},
	}
}

func (log *Logger) AddDestination(dst *LogDestination) *Logger {
	log.destinations = append(log.destinations, dst)
	return log
}

func (log *Logger) AddFile(path string, logDestination *LogDestination) *Logger {
	logFile := openLogFile(path)

	destination := LogDestination{
		output:          logFile,
		format:          TEXT,
		flags:           Ldate | Ltime,
		minimumLogLevel: INFO,
	}
	if logDestination != nil {
		if logDestination.format != 0 {
			destination.format = logDestination.format
		}
		if logDestination.flags != 0 {
			destination.flags = logDestination.flags
		}
		if logDestination.minimumLogLevel != 0 {
			destination.minimumLogLevel = logDestination.minimumLogLevel
		}
		if logDestination.maximumLogLevel != 0 {
			destination.maximumLogLevel = logDestination.maximumLogLevel
		}
	}
	return log.AddDestination(&destination)
}

func (log *Logger) AddStd(logDestination *LogDestination) *Logger {
	destination := LogDestination{
		output:          os.Stdout,
		format:          TEXT,
		flags:           Ldate | Ltime,
		minimumLogLevel: INFO,
	}
	if logDestination != nil {
		if logDestination.format != 0 {
			destination.format = logDestination.format
		}
		if logDestination.flags != 0 {
			destination.flags = logDestination.flags
		}
		if logDestination.minimumLogLevel != 0 {
			destination.minimumLogLevel = logDestination.minimumLogLevel
		}
		if logDestination.maximumLogLevel != 0 {
			destination.maximumLogLevel = logDestination.maximumLogLevel
		}
	}
	return log.AddDestination(&destination)
}

func (dst *LogDestination) SetFlags(flags LogFlag) *LogDestination {
	dst.flags = flags
	return dst
}

func (dst *LogDestination) SetFormat(format LogFormat) *LogDestination {
	dst.format = format
	return dst
}

func (dst *LogDestination) Write(logString []byte) *LogDestination {
	dst.output.Write(logString)
	return dst
}

func (logger *Logger) Write(level LogLevel, logString []byte) *Logger {
	for _, dst := range logger.destinations {
		if ShouldLog(level, dst.minimumLogLevel) && MaximumLogLevel(level, dst.maximumLogLevel) {
			buf := make([]byte, 0)
			flag := dst.flags
			t := time.Now() // get this early.
			if flag&(Ltime|Ldate) != 0 {
				if flag&Ldate != 0 {
					year, month, day := t.Date()
					itoa(&buf, year, 4)
					buf = append(buf, '/')
					itoa(&buf, int(month), 2)
					buf = append(buf, '/')
					itoa(&buf, day, 2)
					buf = append(buf, ' ')
				}
				if flag&(Ltime) != 0 {
					hour, min, sec := t.Clock()
					itoa(&buf, hour, 2)
					buf = append(buf, ':')
					itoa(&buf, min, 2)
					buf = append(buf, ':')
					itoa(&buf, sec, 2)
					buf = append(buf, ' ')
				}
			}
			buf = append(buf, logString...)
			dst.outMutex.Lock()
			_, err := dst.output.Write(buf)
			dst.outMutex.Unlock()
			if err != nil {
				panic(err)
			}
		}
	}
	return logger
}

func (logger *Logger) LogCopy(level LogLevel, reader io.Reader) *Logger {
	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)
	for _, dst := range logger.destinations {
		replayReader := io.MultiReader(&buf, tee)
		if ShouldLog(level, dst.minimumLogLevel) {
			io.Copy(os.Stdout, replayReader)
		}
	}
	return logger
}

var defaultLogger = StdLog()

func SetLogLevel(level LogLevel) {
	for _, dst := range defaultLogger.destinations {
		dst.minimumLogLevel = level
	}
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
