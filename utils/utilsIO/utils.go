package utilsIO

import (
	"io"
	"log"
	"os"
)

func ConfigurarLogger() {
	logFile, err := os.OpenFile("io.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
