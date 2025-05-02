package main

import (
	"log"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsKernel"

	//"encoding/json"
	"net/http"
	//"time"

	//"github.com/sisoputnfrba/tp-golang/estructurasKernel"
	"fmt"
	"os"
	"strconv"
)

func main() {
	globals.CargarConfig("./kernel/globals/config.json")
	utilsKernel.ConfigurarLogger()
	archivo := os.Args[1]
	tamanio, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Tamaño inválido: %v", err)
	}
	println(archivo)
	println(tamanio)
	utilsKernel.CrearPCB(tamanio, archivo)

	go utilsKernel.IniciarPlanifcador()
	go func() {
		time.Sleep(4 * time.Second)
		utilsKernel.PeticionClienteKERNELServidorIO("127.0.0.1", 8003, 8)
	}()
	http.HandleFunc("/handshake", utilsKernel.RetornoClienteIOServidorKERNEL)
	http.HandleFunc("POST /IO", utilsKernel.RetornoClienteIOServidorKERNEL)
	http.HandleFunc("POST /handshake", utilsKernel.RetornoClienteCPUServidorKERNEL)
	http.HandleFunc("POST /PCB", utilsKernel.RetornoClienteCPUServidorKERNEL2)
	log.Printf("Servidor corriendo.\n")
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_kernel), nil)

}
