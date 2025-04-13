package main

import (
	"github.com/sisoputnfrba/tp-golang/utils/utilsMemoria"
	"net/http"
	"log"
)

func main() {
	utilsMemoria.ConfigurarLogger()

	http.HandleFunc("POST /handshake", utilsMemoria.ConexionRecibidaCPU)
	log.Printf("servidor corriendo, peticion CPU")
	http.ListenAndServe(":8002", nil)
}
