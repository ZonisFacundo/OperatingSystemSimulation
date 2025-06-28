package mmu

import (
	"log"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func EstaEnCache(nroPagina int) bool {
	
	if len(globals.CachePaginas.Entradas) == 0 {
		return false
	}
	for i := 0; i < len(globals.CachePaginas.Entradas); i++ {
		if globals.CachePaginas.Entradas[i].NroPag == nroPagina {
			globals.ID.PosicionPag = i
			return true
		}
	}
	return false
}

func WriteEnCache(datos string) {
	pos := globals.ID.PosicionPag
	if pos < 0 || pos >= len(globals.CachePaginas.Entradas) {
		log.Printf("WriteEnCache: índice de página inválido %d", pos)
		return
	}
	globals.CachePaginas.Entradas[globals.ID.PosicionPag].Contenido = datos
	globals.CachePaginas.Entradas[globals.ID.PosicionPag].Modificada = true
	globals.CachePaginas.Entradas[globals.ID.PosicionPag].Modificada = true
}

func ReadEnCache() {
	pos := globals.ID.PosicionPag
	if pos < 0 || pos >= len(globals.CachePaginas.Entradas) {
		log.Printf("ReadEnCache: índice de página inválido %d", pos)
		return
	}
	contenidoCompleto := globals.CachePaginas.Entradas[globals.ID.PosicionPag].Contenido
	desplazamiento := globals.ID.Desplazamiento

	lectura := contenidoCompleto[desplazamiento : desplazamiento+globals.ID.Tamaño]
	log.Printf("READ en cache: %s", lectura)
}


// investigar cache y como convinar con tlb y traducir
