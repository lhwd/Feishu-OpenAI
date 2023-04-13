package utils

import (
	"fmt"
	"log"
	"os"
)

func LogInit() {
	logFile, err := os.OpenFile("./info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetPrefix("larkOpenAI")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
}
