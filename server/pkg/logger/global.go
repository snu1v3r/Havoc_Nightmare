package logger

import (
	"io"
	"log"
	"os"
)

var Log Logger

func init() {
	Log = Logger{
		debug: false,
		log:   log.New(os.Stdout, "", 0),
	}
}

func Info(fmt string, args ...any) {
	Log.Info(fmt, args...)
}

func Debug(fmt string, args ...any) {
	Log.Debug(fmt, args...)
}

func DebugError(fmt string, args ...any) {
	Log.DebugError(fmt, args...)
}

func Warn(fmt string, args ...any) {
	Log.Warn(fmt, args...)
}

func Error(fmt string, args ...any) {
	Log.Error(fmt, args...)
}

func Fatal(fmt string, args ...any) {
	Log.Fatal(fmt, args...)
}

func Panic(fmt string, args ...any) {
	Log.Panic(fmt, args...)
}

func SetDebug(enable bool) {
	Log.SetDebug(enable)
}

func SetStdOut(w io.Writer) {
	Log.log.SetOutput(w)
}
