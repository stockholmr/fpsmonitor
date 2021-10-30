package logging

import (
	"log"
	"os"
	"path"

	"github.com/jcelliott/lumber"
)

var logger *lumber.MultiLogger

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

func Init(logFile string, consoleLogLevel int, fileLogLevel int) error {

	logDir := path.Dir(logFile)
	err := os.Mkdir(logDir, 0776)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	consoleLogger := lumber.NewConsoleLogger(consoleLogLevel)

	fileLogger, err := lumber.NewFileLogger(
		logFile,
		fileLogLevel,
		lumber.ROTATE,
		6000,
		5,
		100,
	)
	if err != nil {
		return err
	}

	logger = lumber.NewMultiLogger()
	logger.AddLoggers(consoleLogger)
	logger.AddLoggers(fileLogger)

	return nil
}

func Error(msg interface{}) {
	if logger != nil {
		logger.Error("%s", msg)
	} else {
		log.Print(msg)
	}
}

func Errorf(fmt string, msg ...interface{}) {
	if logger != nil {
		logger.Error(fmt, msg...)
	} else {
		log.Printf(fmt, msg...)
	}
}

func Trace(msg interface{}) {
	if logger != nil {
		logger.Trace("%s", msg)
	} else {
		log.Print(msg)
	}
}

func Tracef(fmt string, msg ...interface{}) {
	if logger != nil {
		logger.Trace(fmt, msg...)
	} else {
		log.Printf(fmt, msg...)
	}
}

func Info(msg interface{}) {
	if logger != nil {
		logger.Info("%s", msg)
	} else {
		log.Print(msg)
	}
}

func Infof(fmt string, msg ...interface{}) {
	if logger != nil {
		logger.Info(fmt, msg...)
	} else {
		log.Printf(fmt, msg...)
	}
}

func Debug(msg interface{}) {
	if logger != nil {
		logger.Debug("%s", msg)
	} else {
		log.Print(msg)
	}
}

func Debugf(fmt string, msg ...interface{}) {
	if logger != nil {
		logger.Debug(fmt, msg...)
	} else {
		log.Printf(fmt, msg...)
	}
}

func Warning(msg interface{}) {
	if logger != nil {
		logger.Warn("%s", msg)
	} else {
		log.Print(msg)
	}
}

func Warningf(fmt string, msg ...interface{}) {
	if logger != nil {
		logger.Warn(fmt, msg...)
	} else {
		log.Printf(fmt, msg...)
	}
}

func Fatal(msg interface{}) {
	if logger != nil {
		logger.Fatal("%s", msg)
	} else {
		log.Fatal(msg)
	}
}

func Fatalf(fmt string, msg ...interface{}) {
	if logger != nil {
		logger.Fatal(fmt, msg...)
	} else {
		log.Fatalf(fmt, msg...)
	}
}
