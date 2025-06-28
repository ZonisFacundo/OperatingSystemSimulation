package mmu

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func EstaEnCache(nroPagina int) bool {
	for i := 0 ; i < globals.CachePaginas.Tamanio ; i++{
		if globals.CachePaginas.Entradas[i].NroPag == nroPagina {
			globals.ID.PosicionPag = i
			return true
		}
	}
	return false
}

func WriteEnCache(datos string){
	globals.CachePaginas.Entradas[globals.ID.PosicionPag].Contenido = datos
	globals.CachePaginas.Entradas[globals.ID.PosicionPag].Modificada = true
}
// investigar cache y como convinar con tlb y traducir
