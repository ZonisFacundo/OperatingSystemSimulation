package main

import (
	"github.com/sisoputnfrba/tp-golang/utils/utilsIO"
)

func main() {

	utilsIO.ConfigurarLogger()

	/*
	   conexion entre IO (Client) con Kernel (Server)
	   enviamos handshake con datos del modulo y esperamos respuesta
	*/
	utilsIO.HandshakeAKernel("pepe", "127.0.0.1", 8003)

}
