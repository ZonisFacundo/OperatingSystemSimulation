package utilsCPU

import (
	"io"
	"log"
	"os"
	"fmt"
)

/*
func ConfigurarLogger() {
	logFile, err := os.OpenFile("cpu.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
*/

func ConfigurarLogger(cpuId string) {
	logFileName := fmt.Sprintf("CPU-%s.log", cpuId)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	//prefija cada línea de log con el cpuId:
	log.SetPrefix(fmt.Sprintf("[CPU-%s] ", cpuId))
}