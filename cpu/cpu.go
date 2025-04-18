package main

import (
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

func main() {
	utilsCPU.ConfigurarLogger()

	/*
	   conexion entre CPU (Client) con Kernel (Server)
	   enviamos handshake con datos del modulo y esperamos respuesta
	*/

	utilsCPU.PeticionClienteCPUServidorKERNEL("127.0.0.1", 8001)
	utilsCPU.PeticionCLienteCPUServidorMEMORIA("NOOP", "127.0.0.1", 8002)

}
