package main

import (
	"log"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsKernel"

	//"encoding/json"
	"net/http"
	"time"
	//"github.com/sisoputnfrba/tp-golang/estructurasKernel"
)

func main() {
	globals.CargarConfig("./kernel/globals/config.json")
	utilsKernel.ConfigurarLogger()

	time.Sleep(3 * time.Second)
	utilsKernel.PeticionClienteKERNELServidorIO("127.0.0.1", 8003)

	time.Sleep(3 * time.Second)

	//utilsKernel.PeticionClienteKERNELServidorMemoria(0, 250, "127.0.0.1", 8002) debe recibir un pcb despues hago uno de prueba
	http.HandleFunc("/handshake", utilsKernel.RetornoClienteIOServidorKERNEL)
	http.HandleFunc("POST /IO", utilsKernel.RetornoClienteIOServidorKERNEL)
	http.HandleFunc("POST /handshake", utilsKernel.RetornoClienteCPUServidorKERNEL)
	log.Printf("Servidor corriendo.\n")
	http.ListenAndServe(":8001", nil)

}
