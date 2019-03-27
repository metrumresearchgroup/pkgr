package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogrusFileHook struct {
	file      *os.File
	flag      int
	chmod     os.FileMode
	formatter *logrus.JSONFormatter
}

// Log Reinstantiable log to be used globally in the application.
var Log *logrus.Logger

// NewLogrusFileHook
func NewLogrusFileHook(file string, flag int, chmod os.FileMode) (*LogrusFileHook, error) {

	jsonFormatter := &logrus.JSONFormatter{}
	logFile, err := os.OpenFile(file, flag, chmod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook %v", err)
		return nil, err
	}

	return &LogrusFileHook{logFile, flag, chmod, jsonFormatter}, err
}

// Fire event
func (hook *LogrusFileHook) Fire(entry *logrus.Entry) error {

	jsonformat, err := hook.formatter.Format(entry)
	line := string(jsonformat)
	_, err = hook.file.WriteString(line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook(entry.String)%v", err)
		return err
	}

	return nil
}

func (hook *LogrusFileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}

}

func InitLog(outputFile string, level string, overwrite bool) {

	//We want the log to be reset whenever it is initialized.
	Log = logrus.New()

	logLevel := strings.ToLower(level)

	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	switch logLevel {
	case "trace":
		Log.SetLevel(logrus.TraceLevel)
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		Log.SetLevel(logrus.FatalLevel)
	case "panic":
		Log.SetLevel(logrus.PanicLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	if outputFile != "" {

		var osFlag int

		if overwrite {
			osFlag = os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_RDWR
		} else {
			osFlag = os.O_CREATE|os.O_APPEND|os.O_RDWR
		}

		fileHook, err := NewLogrusFileHook(outputFile, osFlag, 0666)
		if err == nil {
			Log.AddHook(fileHook)
		}
	}
}
