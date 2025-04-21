package main

import (
	"log"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/utils/utilsIO"
)

func main() {

	utilsIO.ConfigurarLogger()

	/*
	   conexion entre IO (Client) con Kernel (Server)
	   enviamos handshake con datos del modulo y esperamos respuesta
	*/
	go func() {
		time.Sleep(4 * time.Second)

		utilsIO.PeticionClienteIOServidorKERNEL("pepe", "127.0.0.1", 8001)
	}()

	http.HandleFunc("/KERNELIO", utilsIO.RetornoClienteKERNELServidorIO)
	log.Printf("Servidor IO corriendo.\n")
	http.ListenAndServe(":8003", nil)

}
