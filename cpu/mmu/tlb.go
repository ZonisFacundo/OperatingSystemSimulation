package mmu

import (
	"log"
	"time"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func EstaTraducida(nroPagina int) bool {
	if globals.Tlb.Tamanio > 0 {
		log.Printf("Buscando en TLB -> PID: %d, Pag: %d", globals.ID.ProcessValues.Pid, nroPagina)
		now := time.Now().UnixNano()

		for i, entrada := range globals.Tlb.Entradas {
			if entrada.PID == globals.ID.ProcessValues.Pid && entrada.NroPagina == nroPagina {
				globals.ID.DireccionFis = entrada.Direccion
				globals.ID.PosicionPag = i // si lo necesitás para LRU

				globals.Tlb.Entradas[i].UltimoAcceso = now

				log.Printf(">> TLB HIT -> PID: %d, Pagina: %d -> DirFis: %d", entrada.PID, entrada.NroPagina, entrada.Direccion)
				return true
			}
		}
		log.Printf(">> TLB MISS -> PID: %d, Pagina: %d", globals.ID.ProcessValues.Pid, nroPagina)
		return false
	}
	return false
}
