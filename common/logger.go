package common

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var usaLog *log.Logger
var lock sync.Mutex

func init() {
	file, err := os.OpenFile("./log/main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	usaLog = logger
}

func Log(method int, v ...any) {
	lock.Lock()
	defer lock.Unlock()
	if method == 0 {
		usaLog.SetPrefix("INFO: ")
		defer usaLog.SetPrefix("")
	} else if method == 1 {
		usaLog.SetPrefix("Warn: ")
		defer usaLog.SetPrefix("")
	} else if method == 2 {
		usaLog.SetPrefix("Error: ")
		defer usaLog.SetPrefix("")
	}
	usaLog.Output(2, fmt.Sprint(v...))
}
func Logf(method int, format string, v ...any) {
	lock.Lock()
	defer lock.Unlock()
	if method == 0 {
		usaLog.SetPrefix("INFO: ")
		defer usaLog.SetPrefix("")
	} else if method == 1 {
		usaLog.SetPrefix("Warn: ")
		defer usaLog.SetPrefix("")
	} else if method == 2 {
		usaLog.SetPrefix("Error: ")
		defer usaLog.SetPrefix("")
	}
	usaLog.Output(2, fmt.Sprintf(format, v...))
}

// Logln 0-info 1-warn 2-error
func Logln(method int, v ...any) {
	lock.Lock()
	defer lock.Unlock()

	if method == 0 {
		usaLog.SetPrefix("INFO: ")
		defer usaLog.SetPrefix("")
	} else if method == 1 {
		usaLog.SetPrefix("Warn: ")
		defer usaLog.SetPrefix("")
	} else if method == 2 {
		usaLog.SetPrefix("Error: ")
		defer usaLog.SetPrefix("")
	}
	usaLog.Output(2, fmt.Sprintln(v...))
}
