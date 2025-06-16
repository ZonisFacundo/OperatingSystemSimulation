package main

/*
SI DE CAUSALIDAD VES QUE LAS COSAS NO ANDAN CUANDO CAMBIAS
VALORES EN EL CONFIG
ES PORQUE POR AHI EN CPU ESTAN HARDCODEANDO VALORES (TAM PAGINA) PARA PROBARLO
POR SI NOS OLVIDAMOS
*/
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

	http.HandleFunc("POST /SWAPADISCO", utilsMemoria.RetornoClienteKernelServidorMemoriaSwapADisco)
	http.HandleFunc("POST /KERNELMEMORIADUMP", utilsMemoria.RetornoClienteKernelServidorMemoriaDumpDelProceso)
	http.HandleFunc("POST /READ", utilsMemoria.RetornoClienteCPUServidorMEMORIARead)
	http.HandleFunc("POST /WRITE", utilsMemoria.RetornoClienteCPUServidorMEMORIAWrite)
	http.HandleFunc("POST /TRADUCCIONLOGICAAFISICA", utilsMemoria.RetornoClienteCPUServidorMEMORIATraduccionLogicaAFisica)
	http.HandleFunc("GET /INSTRUCCIONES", utilsMemoria.RetornoClienteCPUServidorMEMORIA)
	http.HandleFunc("POST /KERNELMEMORIA", utilsMemoria.RetornoClienteKernelServidorMEMORIA)
	log.Printf("Servidor corriendo (Memoria) en puerto %d.\n", globals.ClientConfig.Port_memory)
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_memory), nil)
}
