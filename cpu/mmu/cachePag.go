package mmu

import (
	"log"
	"time"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func EstaEnCache(nroPagina int) bool {
	for i, entrada := range globals.CachePaginas.Entradas {
		if entrada.PID == globals.ID.ProcessValues.Pid && entrada.NroPag == nroPagina {
			globals.ID.PosicionPag = i
			globals.CachePaginas.Entradas[i].BitUso = true
			log.Printf(">> CACHE HIT: PID=%d, Pagina=%d", entrada.PID, entrada.NroPag)
			return true
		}
	}
	log.Printf(">> CACHE MISS: PID=%d, Pagina=%d", globals.ID.ProcessValues.Pid, nroPagina)
	return false
}

func WriteEnCache(datos string) {

	time.Sleep(time.Duration(globals.ClientConfig.Cache_delay) * time.Millisecond)

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

	time.Sleep(time.Duration(globals.ClientConfig.Cache_delay) * time.Millisecond)

	pos := globals.ID.PosicionPag
	if pos < 0 || pos >= len(globals.CachePaginas.Entradas) {
		log.Printf("ReadEnCache: índice de página inválido %d", pos)
		return
	}

	contenidoCompleto := globals.CachePaginas.Entradas[pos].Contenido
	desplazamiento := globals.ID.Desplazamiento
	tamanio := globals.ID.Tamaño

	if desplazamiento < 0 || desplazamiento+tamanio > len(contenidoCompleto) {
		log.Printf("ReadEnCache: rango inválido para lectura (desplazamiento %d, tamaño %d, longitud contenido %d)", desplazamiento, tamanio, len(contenidoCompleto))
		return
	}

	lectura := contenidoCompleto[desplazamiento : desplazamiento+tamanio]
	log.Printf("READ en cache: %s", lectura)

	// Marcar uso:
	globals.CachePaginas.Entradas[pos].BitUso = true
}