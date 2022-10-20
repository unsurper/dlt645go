package dlt

import (
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	SetLogLevel(TraceLevel)
	log.StandardLogger().Formatter = &prefixed.TextFormatter{
		ForceColors:      true,
		ForceFormatting:  true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		DisableSorting:   true,
		TimestampFormat:  "2006-01-02 15:04:05",
	}
}

// 日志等级
type Level string

var (
	PanicLevel Level = "panic"
	FatalLevel Level = "fatal"
	ErrorLevel Level = "error"
	WarnLevel  Level = "warn"
	InfoLevel  Level = "info"
	DebugLevel Level = "debug"
	TraceLevel Level = "trace"
)

// 设置日志级别
func SetLogLevel(level Level) error {
	lv, err := log.ParseLevel(string(level))
	if err != nil {
		return err
	}
	log.StandardLogger().SetLevel(lv)
	return nil
}
