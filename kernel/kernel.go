package main

import (
	"log"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsKernel"

	//"encoding/json"
	"net/http"

	//"github.com/sisoputnfrba/tp-golang/estructurasKernel"
	"fmt"
	"os"
	"strconv"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	globals.CargarConfig("./kernel/globals/config.json")
	utilsKernel.ConfigurarLogger()
	archivo := os.Args[1]
	tamanio, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Tamaño inválido: %v", err)
	}
	utilsKernel.InicializarSemaforos()
	go func() {
		utilsKernel.IniciarPlanifcador(tamanio, archivo)
		utilsKernel.PlanificadorLargoPlazo()
	}()

	go utilsKernel.PlanificadorCortoPlazo()

	http.HandleFunc("POST /IO", utilsKernel.RecibirDatosIO)
	http.HandleFunc("POST /handshake", utilsKernel.RecibirDatosCPU)
	http.HandleFunc("POST /PCB", utilsKernel.RecibirProceso)
	http.HandleFunc("POST /finIO", utilsKernel.FinalizarIO)
	log.Printf("Servidor corriendo.\n")
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_kernel), nil)
	wg.Wait()
}
