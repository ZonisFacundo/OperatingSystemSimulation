package auxiliares

import (
	"fmt"
	"time"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	//	"log"
	//	"os"
)

func Mostrarmemoria() {

	for i := 0; i < globals.ClientConfig.Memory_size; i++ {

		fmt.Printf("%v ", globals.MemoriaPrincipal[i])
		time.Sleep(2 * time.Millisecond)
	}
}

func MostrarPaginasDisponibles() {

	fmt.Printf("SI VALE 1 ESTA OCUPADO, SI VALE 0 ESTA DESOCUPADO\n\n")

	for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ {

		fmt.Printf("MARCO NUMERO %d: \t%d\n", i, globals.PaginasDisponibles[i])
	}
}

func MostrarPaginasDisponiblesCada10segundos() {

	for { //equivalente a while true
		fmt.Printf("SI VALE 1 ESTA OCUPADO, SI VALE 0 ESTA DESOCUPADO\n\n")

		for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ {

			fmt.Printf("MARCO NUMERO %d: \t%d\n", i, globals.PaginasDisponibles[i])
		}

		time.Sleep(10 * time.Second)
	}

}
func MostrarMemoriaKernel() {

	for i := 0; i < len(globals.MemoriaKernel); i++ {
		for j := 0; j < len(globals.MemoriaKernel[i].Instrucciones); j++ {
			fmt.Printf("%s", globals.MemoriaKernel[i].Instrucciones[j])
		}
	}
}

func ActualizarInstrucciones(x globals.ProcesoEnMemoria, pid int) { //hay que hacer esto porque no te deja actualizarle solo un miembro del struct directamente al de globals por algun motivo

	auxi := globals.MemoriaKernel[pid]
	auxi.Instrucciones = x.Instrucciones
	globals.MemoriaKernel[pid] = auxi

}

/*
recibe posicion en donde se comenzo a guardar el proceso en swap y tamao de este, se lo asigna a la variable global porque sino no te deja hacerlo directamente :( no se porque
*/

func ActualizarSwapInfo(posicion int64, tam int, pid int) { //hay que hacer esto porque no te deja actualizarle solo un miembro del struct directamente al de globals por algun motivo

	auxi := globals.MemoriaKernel[pid]
	auxi.SwapStart = posicion
	auxi.SwapTam = tam
	globals.MemoriaKernel[pid] = auxi

}

func ActualizarTablaSimple(x globals.ProcesoEnMemoria, pid int) { ///hay que hacer esto porque no te deja actualizarle solo un miembro del struct directamente al de globals por algun motivo

	auxi := globals.MemoriaKernel[pid]
	auxi.TablaSimple = x.TablaSimple
	globals.MemoriaKernel[pid] = auxi

}

func MostrarProceso(pid int) {

	fmt.Printf("pid: %d \n", pid)
	fmt.Printf("paginas ocupadas: \n")

	for i := 0; i < len(globals.MemoriaKernel[pid].TablaSimple); i++ {
		fmt.Printf("%d \n", globals.MemoriaKernel[pid].TablaSimple[i])

	}
}
