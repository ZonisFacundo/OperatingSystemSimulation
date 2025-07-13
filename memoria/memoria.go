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
	utilsMemoria.InicializarSwap()
	//utilsMemoria.EscanearMemoria() DEPRECADO

	//go auxiliares.MostrarPaginasDisponiblesCada10segundos()
	/*
		go func() {
			for {
				auxiliares.Mostrarmemoria()
				time.Sleep(20 * time.Second)
			}
		}()
	*/
	auxiliares.MostrarMemoriaKernel()
	http.HandleFunc("POST /FinProceso", utilsMemoria.RetornoClienteKernelServidorMemoriaFinProceso)
	http.HandleFunc("POST /HANDSHAKE", utilsMemoria.HandshakeACpu)
	http.HandleFunc("POST /SWAPAMEMORIA", utilsMemoria.RetornoClienteKernelServidorMemoriaSwapAMemoria)
	http.HandleFunc("POST /SWAPADISCO", utilsMemoria.RetornoClienteKernelServidorMemoriaSwapADisco)
	http.HandleFunc("POST /KERNELMEMORIADUMP", utilsMemoria.RetornoClienteKernelServidorMemoriaDumpDelProceso)
	http.HandleFunc("POST /READ", utilsMemoria.RetornoClienteCPUServidorMEMORIARead)
	http.HandleFunc("POST /WRITE", utilsMemoria.RetornoClienteCPUServidorMEMORIAWrite)
	http.HandleFunc("POST /TRADUCCIONLOGICAAFISICA", utilsMemoria.RetornoClienteCPUServidorMEMORIATraduccionLogicaAFisica)
	http.HandleFunc("GET /INSTRUCCIONES", utilsMemoria.RetornoClienteCPUServidorMEMORIA)
	http.HandleFunc("POST /KERNELMEMORIA", utilsMemoria.RetornoClienteKernelServidorMEMORIA)
	log.Printf("Servidor corriendo (Memoria) en puerto %d.\n", globals.ClientConfig.Port_memory)

	addr := fmt.Sprintf("%s:%d",
		globals.ClientConfig.Ip_memory,   // la IP que leíste de config.json
		globals.ClientConfig.Port_memory, // el puerto
	)

	log.Printf("Servidor Memoria escuchando en %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error al arrancar servidor en %s: %v\n", addr, err)
	}

}
