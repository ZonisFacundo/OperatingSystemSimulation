package mmu

import (
	"math"

	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
	/*"github.com/sisoputnfrba/tp-golang/cpu/instruction_cycle"
	"github.com/sisoputnfrba/tp-golang/io/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"*/)

/*
1) Definir las estructuras de datos para la tabla de paginas
2) Traducir direccion logica a fisica
3)Gestionar TLB "se implementará una TLB para agilizar la traducción de las direcciones lógicas a direcciones físicas"
La TLB contará con la siguiente estructura base: [ página | marco ]*/

/*func TraducirDireccion(direccionLogica int){

	var memoryManagement MMU

	nro_pagina := (math.Floor(float64(direccionLogica) / float64(memoryManagement.TamPagina)))
	entrada_nivel_X := math.Floor(float64(nro_pagina)/float64(memoryManagement.cant_entradas_tabla ^ (memoryManagement.Niveles))) % float64(memoryManagement.cant_entradas_tabla) //X??
	desplazamiento := direccionLogica % memoryManagement.TamPagina

	//direction[] -> Memoria -> Tabla de Paginas ->
}*/

type MMU struct {
	ProcesoActual       utilsCPU.Proceso
	Niveles             int
	TamPagina           int
	Cant_entradas_tabla int
	TablasPaginas       map[int]int
}

func TraducirDireccion(direccionLogica int, memoryManagement MMU, pid int) []int {
	
	// Calcular el número de página
	memoryManagement.TamPagina = 1000
	nroPagina := direccionLogica / memoryManagement.TamPagina

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

	//desplazamiento := direccionLogica % memoryManagement.TamPagina
	resultado := append([]int{pid}, entradas...) // Agrego el pid al principio del slice y concateno las entradas de nivel

	// Retorno el array con las entradas de nivel + desplazamiento

	return resultado

}
