package main

import (
	log "logger"
	"os"
)

const (
	logPath = "log/"
	logFile = "Server.log"
)

var fileHandle *os.File

/*OpenLoggerFile is for to put the logs*/
func OpenLoggerFile() error {
	err := os.MkdirAll(logPath, 0777)
	if err != nil {
		return err
	}
	fileHandle, err = os.OpenFile(logPath+logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	log.Init(fileHandle, fileHandle, fileHandle, fileHandle)
	return nil
}

/*GetFilehandle is for to put the logs*/
func GetFilehandle() *os.File {
	return fileHandle

}
