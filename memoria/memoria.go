package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsMemoria"
)

func main() {
	utilsMemoria.ConfigurarLogger()
	globals.CargarConfig("./memoria/globals/config.json") //decodifica de json a go y guarda los datos en un puntero (variable global) ClientConfig

	//var wg sync.WaitGroup
	//wg.Add(2)
	/*DESARROLLO DE MEMORIA PRINCIPAL DEL SISTEMA
	como aclara el json, vamos a tener 4096 bytes de MP, consecuentemente 12 bits para las direcciones de memoria

	*/

	http.HandleFunc("POST /CPUMEMORIA", utilsMemoria.RetornoClienteCPUServidorMEMORIA)
	http.HandleFunc("POST /KERNELMEMORIA", utilsMemoria.RetornoClienteKernelServidorMEMORIA)
	log.Printf("Servidor corriendo (Memoria) en puerto %d.\n", globals.ClientConfig.Port_memory)
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_memory), nil)

}
