package main

import (
	"github.com/sisoputnfrba/tp-golang/utils/utilsKernel"
	//"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	utilsKernel.ConfigurarLogger()

	/*
	   conexion entre IO (Client) con Kernel (Server)
	   recibimos handshake de parte del IO con datos del modulo damos respuesta
	*/

	http.HandleFunc("POST /handshake", utilsKernel.ConexionRecibida)
	fmt.Printf("servidor corriendo")
	http.ListenAndServe(":8001", nil)

}
