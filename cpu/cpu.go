package main

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

func main() {
	utilsCPU.ConfigurarLogger()
	globals.CargarConfig("./cpu/globals/config.json")

	/*
	   conexion entre CPU (Client) con Kernel (Server)
	   enviamos handshake con datos del modulo y esperamos respuesta
	*/

	utilsCPU.PeticionClienteCPUServidorKERNEL(globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel)
	//utilsCPU.PeticionCLienteCPUServidorMEMORIA("NOOP", globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)

}
