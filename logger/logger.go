package logger

import (
	"fmt"
	"os"
	"strings"

	logrus "github.com/sirupsen/logrus"
)

type LogrusFileHook struct {
	file      *os.File
	flag      int
	chmod     os.FileMode
	formatter *logrus.JSONFormatter
}

// Log Reinstantiable logrus to be used globally in the application.
//var Log *logrus.Logger

func init() {
	//Log = logrus.New()
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
}


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

func SetLogLevel(level string) {
	//We want the logrus to be reset whenever it is initialized.
	logLevel := strings.ToLower(level)

	switch logLevel {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func AddLogFile(outputFile string, overwrite bool) {

	if outputFile != "" {

		var osFlag int

		if overwrite {
			osFlag = os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_RDWR
		} else {
			osFlag = os.O_CREATE|os.O_APPEND|os.O_RDWR
		}

		fileHook, err := NewLogrusFileHook(outputFile, osFlag, 0666)
		if err == nil {
			logrus.AddHook(fileHook)
		}
	}
}
