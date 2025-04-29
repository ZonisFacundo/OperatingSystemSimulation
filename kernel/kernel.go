package main

import (
	"log"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsKernel"

	//"encoding/json"
	"net/http"
	"time"

	//"github.com/sisoputnfrba/tp-golang/estructurasKernel"
	"fmt"
)

func main() {
	globals.CargarConfig("./kernel/globals/config.json")
	utilsKernel.ConfigurarLogger()

	time.Sleep(1 * time.Second)
	utilsKernel.PeticionClienteKERNELServidorIO("127.0.0.1", 8003)

	time.Sleep(2 * time.Second)

	//utilsKernel.PeticionClienteKERNELServidorMemoria(0, 250, "127.0.0.1", 8002) //debe recibir un pcb despues hago uno de prueba
	//utilsKernel.PeticionClienteKERNELServidorMemoria(5, 250, "127.0.0.1", 8002) //con este codigo y la parte comentada de cliente kernel anda, entiendo que quieren pasarle un pcb o algo asi pero lo use para probar

	//utilsKernel.PeticionClienteKERNELServidorCPU() // ??? que seria pcb PCB (al llamar a la función deben ir parametros, cuales?)
	http.HandleFunc("/handshake", utilsKernel.RetornoClienteIOServidorKERNEL)
	http.HandleFunc("POST /IO", utilsKernel.RetornoClienteIOServidorKERNEL)
	http.HandleFunc("POST /handshake", utilsKernel.RetornoClienteCPUServidorKERNEL)
	http.HandleFunc("POST /PCB", utilsKernel.RetornoClienteCPUServidorKERNEL2)
	log.Printf("Servidor corriendo.\n")
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_kernel), nil)

}
