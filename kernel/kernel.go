package main

import (
	"github.com/sisoputnfrba/tp-golang/utils/utilsKernel"
	//"encoding/json"
	"log"
	"net/http"
	"time"
	//"github.com/sisoputnfrba/tp-golang/estructurasKernel"
)

func main() {
	utilsKernel.ConfigurarLogger()

	time.Sleep(3 * time.Second)
	utilsKernel.PeticionClienteKERNELServidorIO("127.0.0.1", 8003)

	time.Sleep(3 * time.Second)

	utilsKernel.PeticionClienteKERNELServidorMemoria(0, 250, "127.0.0.1", 8002)
	http.HandleFunc("/handshake", utilsKernel.RetornoClienteIOServidorKERNEL)
	http.HandleFunc("POST /IO", utilsKernel.RetornoClienteIOServidorKERNEL)
	http.HandleFunc("POST /handshake", utilsKernel.RetornoClienteCPUServidorKERNEL)
	log.Printf("Servidor corriendo.\n")
	http.ListenAndServe(":8001", nil)

}
