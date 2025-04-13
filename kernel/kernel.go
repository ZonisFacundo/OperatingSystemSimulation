package main

import (
	"github.com/sisoputnfrba/tp-golang/utils/utilsKernel"
	//"encoding/json"
	"net/http"
	"log"
)

func main() {
	utilsKernel.ConfigurarLogger()

	/*
	   conexion entre IO (Client) con Kernel (Server)
	   recibimos handshake de parte del IO con datos del modulo damos respuesta
	
	
	http.HandleFunc("POST /handshake", utilsKernel.ConexionRecibidaIO)
	log.Printf("servidor corriendo, peticion io")
	http.ListenAndServe(":8001", nil)
*/
	http.HandleFunc("POST /handshake", utilsKernel.ConexionRecibidaCPU)
	log.Printf("servidor corriendo, peticion cpu")
	http.ListenAndServe(":8001", nil)

}
