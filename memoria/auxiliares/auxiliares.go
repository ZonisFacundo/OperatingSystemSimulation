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

		fmt.Printf("%v\n", globals.MemoriaPrincipal[i])
		time.Sleep(10 * time.Millisecond)
	}
}
func MostrarPaginasDisponibles() {

	fmt.Printf("SI VALE 1 ESTA OCUPADO, SI VALE 0 ESTA DESOCUPADO\n\n")

	for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ {

		fmt.Printf("MARCO NUMERO %d: \t%d\n", i, globals.PaginasDisponibles[i])
	}
}
