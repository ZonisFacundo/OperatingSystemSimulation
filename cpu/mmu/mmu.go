package mmu

import (
	"log"
	"math"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func TraducirDireccion(direccionLogica int, pid int, nroPagina int) []int {

	if globals.ClientConfig.Page_size == 0 { //Tampagina
		log.Fatalf("Error: TamPagina no puede ser 0. Verificá la configuración o la inicialización de la MMU.")
	}

	// Crear un slice para guardar las entradas de las tablas de páginas
	entradas := make([]int, globals.ClientConfig.Niveles)

	// Calcular las entradas de la tabla de páginas para cada nivel
	for x := 1; x <= globals.ClientConfig.Niveles; x++ {
		exp := globals.ClientConfig.Niveles - x
		divisor := int(math.Pow(float64(globals.ClientConfig.Entradas), float64(exp)))

		// Calculamos la entrada en el nivel X
		entradaNivelX := (nroPagina / divisor) % globals.ClientConfig.Entradas
		entradas[x-1] = entradaNivelX
	}

	desplazamiento := direccionLogica % globals.ClientConfig.Page_size

	log.Printf("Desplazamiento?: %d", desplazamiento)
	globals.ID.Desplazamiento = desplazamiento

	resultado := append([]int{pid}, entradas...) // Agrego el pid al principio del slice y concateno las entradas de nivel

	log.Println("que es esto?: ", resultado)

	return resultado

}
