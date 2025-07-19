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
			log.Printf("PID: %d - Cache Hit - Pagina: %d", entrada.PID, entrada.NroPag)
			return true
		}
	}
	log.Printf("PID: %d - Cache MISS - Pagina: %d", globals.ID.ProcessValues.Pid, nroPagina)
	return false
}

func WriteEnCache(pid int, nroPag int, despl int, datos []byte) {

	time.Sleep(time.Duration(globals.ClientConfig.Cache_delay) * time.Millisecond)

	/*
		pos := globals.ID.PosicionPag
		if pos < 0 || pos >= len(globals.CachePaginas.Entradas) {
			log.Printf("## WriteEnCache: índice de página inválido %d", pos)
			return
		}

		//globals.CachePaginas.Entradas[pos].Contenido = datos
		//globals.CachePaginas.Entradas[pos].Modificada = true
		//globals.CachePaginas.Entradas[pos].BitUso = true

		entrada := &globals.CachePaginas.Entradas[pos]

		if entrada.Desplazamiento < 0 || entrada.Desplazamiento+len(datos) > len(entrada.PaginaCompleta) {
			log.Printf("## WriteEnCache: rango inválido para escritura (Desplazamiento: %d, Tamaño: %d, Longitud de la página: %d)", entrada.Desplazamiento, len(datos), len(entrada.PaginaCompleta))
			return
		}

		copy(entrada.PaginaCompleta[entrada.Desplazamiento:], datos)

		entrada.Modificada = true
		entrada.BitUso = true*/

		for i := range globals.CachePaginas.Entradas {
			entrada := &globals.CachePaginas.Entradas[i]
		
			if entrada.PID == pid && entrada.NroPag == nroPag {
	
				copy(entrada.PaginaCompleta[despl:], datos)
		
				entrada.Modificada = true
				entrada.BitUso = true

				log.Printf("## PID: %d - Accion: ESCRIBIR - Direccion Física: %d - Valor: %v",
					globals.ID.ProcessValues.Pid, globals.ID.DireccionFis, globals.ID.Datos)
				return
			}
		}
		
	}


func ReadEnCache() {

	time.Sleep(time.Duration(globals.ClientConfig.Cache_delay) * time.Millisecond)

	pos := globals.ID.PosicionPag
	if pos < 0 || pos >= len(globals.CachePaginas.Entradas) {
		log.Printf("## ReadEnCache: índice de página inválido %d", pos)
		return
	}

	pagCompleta := globals.CachePaginas.Entradas[pos].PaginaCompleta
	desplazamiento := globals.ID.Desplazamiento
	tamanio := globals.ID.Tamaño

	if desplazamiento < 0 || desplazamiento+tamanio > len(pagCompleta) {
		log.Printf("## ReadEnCache: rango inválido para lectura (Desplazamiento: %d, Tamaño: %d, Longitud del contenido: %d)", desplazamiento, tamanio, len(pagCompleta))
		return
	}

	globals.ID.LecturaCache = pagCompleta[desplazamiento : desplazamiento+tamanio]

	//log.Printf("## READ en cache: %v", globals.ID.LecturaCache)
	log.Printf("## PID: %d - Accion: LEER - Direccion Física: %d - Valor: %v",
		globals.ID.ProcessValues.Pid, globals.ID.DireccionFis, globals.ID.LecturaCache)

	//lectura := contenidoCompleto[desplazamiento : desplazamiento+tamanio]
	//log.Printf("## READ en cache: %s", lectura)

	// Marcar uso:
	globals.CachePaginas.Entradas[pos].BitUso = true
}
