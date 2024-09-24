package server

import "Havoc/pkg/logger"

func (*Teamserver) LogInfo(fmt string, args ...any) {
	logger.Info(fmt, args...)
}

func (*Teamserver) LogError(fmt string, args ...any) {
	logger.Error(fmt, args...)
}

func (*Teamserver) LogWarn(fmt string, args ...any) {
	logger.Warn(fmt, args...)
}

func (*Teamserver) LogFatal(fmt string, args ...any) {
	logger.Fatal(fmt, args...)
}

func (*Teamserver) LogPanic(fmt string, args ...any) {
	logger.Panic(fmt, args...)
}

func (*Teamserver) LogDebug(fmt string, args ...any) {
	logger.Log.Debug(fmt, args...)
}

func (*Teamserver) LogDbgError(fmt string, args ...any) {
	logger.Log.DebugError(fmt, args...)
}
