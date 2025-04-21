package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/auxiliares"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsMemoria"
)

func main() {
	globals.CargarConfig("./memoria/globals/config.json") //decodifica de json a go y guarda los datos en un puntero (variable global) ClientConfig
	utilsMemoria.ConfigurarLogger()

	utilsMemoria.InicializarMemoria()
	utilsMemoria.InicializarPaginasDisponibles()
	utilsMemoria.ActualizaPaginasDisponibles()
	auxiliares.MostrarPaginasDisponibles()
	//auxiliares.Mostrarmemoria()

	http.HandleFunc("POST /CPUMEMORIA", utilsMemoria.RetornoClienteCPUServidorMEMORIA)
	http.HandleFunc("POST /KERNELMEMORIA", utilsMemoria.RetornoClienteKernelServidorMEMORIA)
	log.Printf("Servidor corriendo (Memoria) en puerto %d.\n", globals.ClientConfig.Port_memory)
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_memory), nil)

}
