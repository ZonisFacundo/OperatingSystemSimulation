package main

import (
	"log"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsKernel"

	//"encoding/json"
	"net/http"
	//"time"

	//"github.com/sisoputnfrba/tp-golang/estructurasKernel"
	"fmt"
)

func main() {
	globals.CargarConfig("./kernel/globals/config.json")
	utilsKernel.ConfigurarLogger()
	utilsKernel.CrearPCB(2, "/home/utnso/Desktop/tp-2025-1c-NutriGO/archi.txt") //Esta hardcodeado para probar

	go utilsKernel.IniciarPlanifcador()
	/*
		time.Sleep(1 * time.Second)
		utilsKernel.PeticionClienteKERNELServidorIO("127.0.0.1", 8003)

		time.Sleep(2 * time.Second)
	*/

	//http.HandleFunc("/handshake", utilsKernel.RetornoClienteIOServidorKERNEL)
	//http.HandleFunc("POST /IO", utilsKernel.RetornoClienteIOServidorKERNEL)
	http.HandleFunc("POST /handshake", utilsKernel.RetornoClienteCPUServidorKERNEL)
	http.HandleFunc("POST /PCB", utilsKernel.RetornoClienteCPUServidorKERNEL2)
	log.Printf("Servidor corriendo.\n")
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_kernel), nil)

}
