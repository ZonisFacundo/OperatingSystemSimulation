package mmu

import "github.com/sisoputnfrba/tp-golang/cpu/globals"

func EstaTraducida(nroPagina int) (bool) {

	for i := 0; i < globals.Tlb.Tamanio; i++ {

		if globals.Tlb.Entradas[i].NroPagina == nroPagina {
			
			globals.ID.DireccionFis = globals.Tlb.Entradas[i].Direccion
			return true
		}
	}

	return false
}

/*
Acá faltaría:
1. Ingresar las direcciones ya traducidas para guardarlas. Basicamente guardarlas en slice
2. Sacar las direcciones ya traducidas en caso de algoritmo de reemplazo. Preparar algoritmos
*/