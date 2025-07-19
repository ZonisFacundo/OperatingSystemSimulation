package mmu

import (
	"log"
	"time"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func EstaTraducida(nroPagina int) bool {
	if globals.Tlb.Tamanio > 0 {
		//log.Printf("Buscando en TLB -> PID: %d, Pag: %d", globals.ID.ProcessValues.Pid, nroPagina)
		now := time.Now().UnixNano()

		for i, entrada := range globals.Tlb.Entradas {
			if entrada.PID == globals.ID.ProcessValues.Pid && entrada.NroPagina == nroPagina {
				globals.ID.DireccionFis = entrada.Direccion
				globals.ID.PosicionPag = i // si lo necesitás para LRU

				globals.Tlb.Entradas[i].UltimoAcceso = now

				log.Printf("PID: %d - TLB HIT - Pagina: %d", entrada.PID, entrada.NroPagina) //OBLIGATORIO
				return true
			}
		}
		log.Printf("PID: %d - TLB MISS - Pagina: %d", globals.ID.ProcessValues.Pid, nroPagina) //OBLIGATORIO
		return false
	}
	return false
}
