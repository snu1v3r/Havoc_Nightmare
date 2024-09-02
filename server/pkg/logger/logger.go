package logger

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"Havoc/pkg/colors"
)

func FunctionTrace() (string, int) {
	var (
		frame    runtime.Frame
		frames   *runtime.Frames
		caller   = make([]uintptr, 15)
		callNums int
	)

	callNums = runtime.Callers(2, caller)

	frames = runtime.CallersFrames(caller[:callNums])

	frame, _ = frames.Next()
	frame, _ = frames.Next()
	frame, _ = frames.Next()

	return frame.Function, frame.Line
}

type Logger struct {
	log   *log.Logger
	debug bool
}

func (*Logger) time() string {
	return colors.BoldWhite(time.Now().Format("2-Jan-2006 15:04:05"))
}

func (logger *Logger) Info(fmt string, args ...any) {
	logger.log.SetPrefix("[" + logger.time() + " " + colors.BoldGreen("INFO") + "] ")
	logger.log.Printf(fmt+"\n", args...)
}

func (logger *Logger) Debug(fmt string, args ...any) {
	var Trace, Line = FunctionTrace()
	var Functions = strings.Split(Trace, "/")
	if logger.debug {
		logger.log.SetPrefix("[" + logger.time() + " " + colors.BoldYellow("DBUG") + "] (" + colors.BlueUnderline(Functions[len(Functions)-1]+":"+strconv.Itoa(Line)) + ") ")
		logger.log.Printf(fmt+"\n", args...)
	}
}

func (logger *Logger) DebugError(fmt string, args ...any) {
	var Trace, Line = FunctionTrace()
	var Functions = strings.Split(Trace, "/")
	if logger.debug {
		logger.log.SetPrefix("[" + logger.time() + " " + colors.BoldRed("DBER") + "] (" + colors.BlueUnderline(Functions[len(Functions)-1]+":"+strconv.Itoa(Line)) + ") ")
		logger.log.Printf(fmt+"\n", args...)
	}
}

func (logger *Logger) Warn(fmt string, args ...any) {
	logger.log.SetPrefix("[" + logger.time() + " " + colors.BoldYellow("WARN") + "] ")
	logger.log.Printf(fmt+"\n", args...)
}

func (logger *Logger) Error(fmt string, args ...any) {
	logger.log.SetPrefix("[" + logger.time() + " " + colors.BoldRed("ERRO") + "] ")
	logger.log.Printf(fmt+"\n", args...)
}

func (logger *Logger) Fatal(fmt string, args ...any) {
	logger.log.SetPrefix("[" + logger.time() + " " + colors.BoldRed("FATA") + "] ")
	logger.log.Printf(fmt+"\n", args...)
	os.Exit(1)
}

func (logger *Logger) Panic(fmt string, args ...any) {
	logger.log.SetPrefix("[" + logger.time() + " " + colors.BoldRed("PANIC") + "] ")
	logger.log.Printf(fmt+"\n", args...)
	panic(args)
}

func (logger *Logger) SetDebug(enable bool) {
	logger.debug = enable
}
