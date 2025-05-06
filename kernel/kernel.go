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

	go utilsKernel.IniciarPlanifcador(tamanio, archivo)
	/*
		go func() {
			time.Sleep(4 * time.Second)
			utilsKernel.UtilizarIO("127.0.0.1", 8003, 0, 8, "impresora")
		}()
	*/
	//http.HandleFunc("/handshake", utilsKernel.RecibirDatosIO) no se porque esta esto por las dudas no lo borro
	http.HandleFunc("POST /IO", utilsKernel.RecibirDatosIO)
	http.HandleFunc("POST /handshake", utilsKernel.RecibirDatosCPU)
	http.HandleFunc("POST /PCB", utilsKernel.RecibirProceso)
	log.Printf("Servidor corriendo.\n")
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_kernel), nil)

}
