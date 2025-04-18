package main

import (
	"github.com/sisoputnfrba/tp-golang/utils/utilsKernel"
	//"encoding/json"
	"log"
	"net/http"
	"time"

	//"github.com/sisoputnfrba/tp-golang/estructurasKernel"
	"sync"
)

func main() {
	utilsKernel.ConfigurarLogger()
	var wg sync.WaitGroup
	wg.Add(3)
	/*
	   conexion entre IO (Client) con Kernel (Server)
	   recibimos handshake de parte del IO con datos del modulo damos respuesta
	*/

	go func() {
		time.Sleep(3 * time.Second)
		utilsKernel.PeticionClienteKERNELServidorIO("127.0.0.1", 8003)
	}()

	http.HandleFunc("/handshake", utilsKernel.RetornoClienteIOServidorKERNEL)
	log.Printf("servidor corriendo, peticion io")
	go http.ListenAndServe(":8001", nil)

	time.Sleep(3 * time.Second)

	go utilsKernel.PeticionClienteKERNELServidorMemoria("codigo", 25, "127.0.0.1", 8002)

	http.HandleFunc("POST /IO", utilsKernel.RetornoClienteIOServidorKERNEL)
	http.HandleFunc("POST /handshake", utilsKernel.RetornoClienteCPUServidorKERNEL)
	log.Printf("Servidor corriendo.\n")
	http.ListenAndServe(":8001", nil)

	wg.Wait()
	//va andar cuando implementemos hilos

}
