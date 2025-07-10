package mmu

import (
	"log"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func EstaEnCache(nroPagina int) bool {
	for i, entrada := range globals.CachePaginas.Entradas {
		if entrada.PID == globals.ID.Pid && entrada.NroPag == nroPagina {
			globals.ID.PosicionPag = i
			globals.CachePaginas.Entradas[i].BitUso = true
			log.Printf(">> CACHE HIT: PID=%d, Pagina=%d", entrada.PID, entrada.NroPag)
			return true
		}
	}
	log.Printf(">> CACHE MISS: PID=%d, Pagina=%d", globals.ID.Pid, nroPagina)
	return false
}

func WriteEnCache(datos string) {
	pos := globals.ID.PosicionPag
	if pos < 0 || pos >= len(globals.CachePaginas.Entradas) {
		log.Printf("WriteEnCache: índice de página inválido %d", pos)
		return
	}

	globals.CachePaginas.Entradas[pos].Contenido = datos
	globals.CachePaginas.Entradas[pos].Modificada = true
	globals.CachePaginas.Entradas[pos].BitUso = true
}

func ReadEnCache() {
	pos := globals.ID.PosicionPag
	if pos < 0 || pos >= len(globals.CachePaginas.Entradas) {
		log.Printf("ReadEnCache: índice de página inválido %d", pos)
		return
	}
	
	contenidoCompleto := globals.CachePaginas.Entradas[pos].Contenido
	desplazamiento := globals.ID.Desplazamiento

	lectura := contenidoCompleto[desplazamiento : desplazamiento+globals.ID.Tamaño]
	log.Printf("READ en cache: %s", lectura)

	// Marcar uso:
	globals.CachePaginas.Entradas[pos].BitUso = true
}