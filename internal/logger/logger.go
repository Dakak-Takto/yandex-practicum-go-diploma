package logger

import (
	"fmt"
	"path"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var once sync.Once

func GetLoggerInstance() *logrus.Logger {

	once.Do(func() {
		log = logrus.StandardLogger()
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "15:05:05",
			FullTimestamp:   true,
			ForceColors:     true,
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {

				return "", fmt.Sprintf(" %s:%d", path.Base(f.File), f.Line)
			},
		})
		log.SetReportCaller(true)
		log.SetLevel(logrus.DebugLevel)
		log.Debug("init logger")
	})

	return log
}
