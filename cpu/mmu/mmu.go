package mmu

import (
	"log"
	"math"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

type MMU struct {
	Pc                  int
	Pid                 int
	Niveles             int
	TamPagina           int
	Cant_entradas_tabla int
	TablasPaginas       map[int]int
}

func TraducirDireccion(direccionLogica int,memoryManagement MMU, pid int, nroPagina int) []int {

	if memoryManagement.TamPagina == 0 {
		log.Fatalf("Error: TamPagina no puede ser 0. Verificá la configuración o la inicialización de la MMU.")
	}

	// Crear un slice para guardar las entradas de las tablas de páginas
	entradas := make([]int, memoryManagement.Niveles)

	// Calcular las entradas de la tabla de páginas para cada nivel
	for x := 1; x <= memoryManagement.Niveles; x++ {
		exp := memoryManagement.Niveles - x
		divisor := int(math.Pow(float64(memoryManagement.Cant_entradas_tabla), float64(exp)))

		// Calculamos la entrada en el nivel X
		entradaNivelX := (nroPagina / divisor) % memoryManagement.Cant_entradas_tabla
		entradas[x-1] = entradaNivelX
	}

	desplazamiento := direccionLogica % memoryManagement.TamPagina

	globals.ID.Desplazamiento = desplazamiento

	resultado := append([]int{pid}, entradas...) // Agrego el pid al principio del slice y concateno las entradas de nivel

	return resultado

}
