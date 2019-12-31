package tools

import (
	"log"
	"os"
)

const logFileName string = ".dep2go/logs.log"

func WriteLog (message string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(dir + "/" +logFileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Println(message)
	logger.Println("----------------------------------------------------------")
}