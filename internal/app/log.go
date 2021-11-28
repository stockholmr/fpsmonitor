package app

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func (a *App) Fatal(v ...interface{}) {
	a.Log().Fatal(append([]interface{}{"FATAL"}, v...))
}

func (a *App) Error(v ...interface{}) {
	a.Log().Print(append([]interface{}{"ERROR"}, v...))
}

func (a *App) Warning(v ...interface{}) {
	a.Log().Print(append([]interface{}{"WARN"}, v...))
}

func (a *App) Info(v ...interface{}) {
	a.Log().Print(append([]interface{}{"INFO"}, v...))
}

func (a *App) Debug(v ...interface{}) {
	a.Log().Print(append([]interface{}{"DEBUG"}, v...))
}

func (a *App) SetLog(log *log.Logger) {
	a.log = log
}

func (a *App) InitFileLog(file string) {

	dir := path.Dir(file)
	err := os.Mkdir(dir, 0776)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	prefix := fmt.Sprintf(
		"%s ",
		time.Now().Format("2006-01-02 15:04:05"),
	)

	logRoller := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	}

	a.log = log.New(logRoller, prefix, 0)
}

func (a *App) Log() *log.Logger {
	if a.log == nil {
		a.log = log.Default()
	}

	return a.log
}
