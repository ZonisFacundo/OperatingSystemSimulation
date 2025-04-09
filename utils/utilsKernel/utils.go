package utilsKernel

import (
	"io"
	"log"
	"os"

	//"encoding/json"
	"net/http"
)

func ConfigurarLogger() {
	logFile, err := os.OpenFile("kernel.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func ConexionRecibida(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)

}
