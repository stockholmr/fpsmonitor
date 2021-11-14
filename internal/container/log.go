package container

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func (c *container) Fatal(v ...interface{}) {
	c.Log().Fatal(append([]interface{}{"FATAL"}, v...))
}

func (c *container) Error(v ...interface{}) {
	c.Log().Print(append([]interface{}{"ERROR"}, v...))
}

func (c *container) Warning(v ...interface{}) {
	c.Log().Print(append([]interface{}{"WARN"}, v...))
}

func (c *container) Info(v ...interface{}) {
	c.Log().Print(append([]interface{}{"INFO"}, v...))
}

func (c *container) Debug(v ...interface{}) {
	c.Log().Print(append([]interface{}{"DEBUG"}, v...))
}

func (c *container) SetLog(log *log.Logger) {
	c.log = log
}

func (c *container) InitFileLog(file string) {

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

	c.log = log.New(logRoller, prefix, 0)
}

func (c *container) Log() *log.Logger {
	if c.log == nil {
		c.log = log.Default()
	}

	return c.log
}
