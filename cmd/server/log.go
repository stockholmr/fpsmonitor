package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLog(file string) *log.Logger {

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

	return log.New(logRoller, prefix, 0)
}
