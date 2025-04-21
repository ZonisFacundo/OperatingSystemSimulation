package main

import (
	"github.com/sisoputnfrba/tp-golang/utils/utilsKernel"
	//"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	//"github.com/sisoputnfrba/tp-golang/estructurasKernel"
)

func main() {
	globals.CargarConfiguracion("./kernel/globals/config.json")
	utilsKernel.ConfigurarLogger()

	/*
	   conexion entre IO (Client) con Kernel (Server)
	   recibimos handshake de parte del IO con datos del modulo damos respuesta
	*/
	go func() {
		time.Sleep(2 * time.Second)
		utilsKernel.PeticionClienteKERNELServidorIO("127.0.0.1", 8003)
	}()

	http.HandleFunc("/handshake", utilsKernel.RetornoClienteIOServidorKERNEL)
	log.Printf("servidor corriendo, peticion io")
	http.ListenAndServe(":8001", nil)

	/*
		http.HandleFunc("POST /IO", utilsKernel.RetornoClienteIOServidorKERNEL)
		http.HandleFunc("POST /handshake", utilsKernel.RetornoClienteCPUServidorKERNEL)
		log.Printf("Servidor corriendo.\n")
		http.ListenAndServe(":8001", nil)
	*/

	//utilsKernel.PeticionClienteKERNELServidorMemoria("codigo", 25, "127.0.0.1", 8002)
	//va andar cuando implementemos hilos

}
