package mmu

import (
	"math"

	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
	/*"github.com/sisoputnfrba/tp-golang/cpu/instruction_cycle"
	"github.com/sisoputnfrba/tp-golang/io/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"*/)

type MMU struct {
	ProcesoActual       utilsCPU.Proceso
	Niveles             int
	TamPagina           int
	cant_entradas_tabla int
	TablasPaginas       map[int]int
}

/*
1) Definir las estructuras de datos para la tabla de paginas
2) Traducir direccion logica a fisica
3)Gestionar TLB "se implementará una TLB para agilizar la traducción de las direcciones lógicas a direcciones físicas"
La TLB contará con la siguiente estructura base: [ página | marco ]*/

func TraducirDireccion(direccionLogica int){

	var memoryManagement MMU

	nro_pagina := (math.Floor(float64(direccionLogica) / float64(memoryManagement.TamPagina)))
	entrada_nivel_X := math.Floor(float64(nro_pagina)/float64(memoryManagement.cant_entradas_tabla ^ (memoryManagement.Niveles))) % float64(memoryManagement.cant_entradas_tabla) //X??
	desplazamiento := direccionLogica % memoryManagement.TamPagina

	//direction[] -> Memoria -> Tabla de Paginas -> 
}
