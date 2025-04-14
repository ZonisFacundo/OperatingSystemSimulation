package main

import (
	"log"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsKernel"

	// "encoding/json"
	"fmt"
	"net/http"
)

func main() {
	utilsKernel.ConfigurarLogger()
	globals.CargarConfig("./kernel/globals/config.json")

	/*
		   conexion entre IO (Client) con Kernel (Server)
		   recibimos handshake de parte del IO con datos del modulo damos respuesta


		http.HandleFunc("POST /handshake", utilsKernel.ConexionRecibidaIO)
		log.Printf("servidor corriendo, peticion io")
		http.ListenAndServe(":8001", nil)
	*/

	//http.HandleFunc("POST /IO", utilsKernel.RetornoClienteIOServidorKERNEL)

	http.HandleFunc("POST /handshake", utilsKernel.RetornoClienteCPUServidorKERNEL)
	log.Printf("Servidor corriendo (kernel).\n")
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_kernel), nil)

	//utilsKernel.PeticionClienteKERNELServidorIO("127.0.0.1", 8003)
	// utilsKernel.PeticionClienteKERNELServidorMemoria("codigo", 25, "127.0.0.1", 8002)
	//va andar cuando implementemos hilos

}
