package main

import (
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/utils/utilsMemoria"
)

func main() {
	utilsMemoria.ConfigurarLogger()

	http.HandleFunc("POST /handshake", utilsMemoria.RetornoClienteCPUServidorMEMORIA)
	http.HandleFunc("POST /KERNELMEMORIA", utilsMemoria.RetornoClienteKernelServidorMEMORIA)
	log.Printf("Servidor corriendo, peticion CPU.\n")
	http.ListenAndServe(":8002", nil)
}
