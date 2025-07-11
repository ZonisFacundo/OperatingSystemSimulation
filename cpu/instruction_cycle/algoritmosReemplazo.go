package instruction_cycle

import(

	"github.com/sisoputnfrba/tp-golang/cpu/globals"

	"log"
	"time"
)

func ReemplazarTLB_FIFO(entrada globals.Entrada) {
	tlb := &globals.Tlb
	if len(tlb.Entradas) < tlb.Tamanio {
		tlb.Entradas = append(tlb.Entradas, entrada)
		return
	}

	tlb.Entradas[tlb.PosDeReemplazo] = entrada
	tlb.PosDeReemplazo = (tlb.PosDeReemplazo + 1) % tlb.Tamanio
}

func ReemplazarTLB_LRU(entrada globals.Entrada) {
	tlb := &globals.Tlb
	now := time.Now().UnixNano()

	if len(tlb.Entradas) < tlb.Tamanio {
		entrada.UltimoAcceso = now
		tlb.Entradas = append(tlb.Entradas, entrada)
		return
	}

	posVictima := 0
	minAcceso := tlb.Entradas[0].UltimoAcceso

	for i := 1; i < len(tlb.Entradas); i++ {
		if tlb.Entradas[i].UltimoAcceso < minAcceso {
			posVictima = i
			minAcceso = tlb.Entradas[i].UltimoAcceso
		}
	}

	entrada.UltimoAcceso = now
	tlb.Entradas[posVictima] = entrada
}

func ReemplazarConCLOCK(entradaNueva globals.EntradaCacheDePaginas) {
	cache := &globals.CachePaginas
	tamanio := cache.Tamanio
	for {
		pos := cache.PosReemplazo
		candidato := &cache.Entradas[pos]

		if !candidato.BitUso {

			if candidato.Modificada {
				frameBase := candidato.NroPag * globals.ClientConfig.Page_size
				Write(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, frameBase, candidato.Contenido)
				log.Printf("PID: %d - Memory Update - Página: %d - Frame: %d", candidato.PID, candidato.NroPag, frameBase)
			}

			*candidato = entradaNueva
			candidato.BitUso = true  // porque la acabamos de traer y usar

			cache.PosReemplazo = (cache.PosReemplazo + 1) % tamanio
			return
		} else {
			candidato.BitUso = false
			cache.PosReemplazo = (cache.PosReemplazo + 1) % tamanio
		}
	}
}

func ReemplazarConCLOCKM(entradaNueva globals.EntradaCacheDePaginas) {
	cache := &globals.CachePaginas
	tamanio := cache.Tamanio

	for {
		for i := 0; i < tamanio; i++ {
			pos := (cache.PosReemplazo + i) % tamanio
			candidato := &cache.Entradas[pos]

			if !candidato.BitUso && !candidato.Modificada {
				*candidato = entradaNueva
				candidato.BitUso = true
				cache.PosReemplazo = (pos + 1) % tamanio
				return
			}
		}
		// si no encuentro Modificada=0 y BitUso=0 entonces:
		for i := 0; i < tamanio; i++ {
			pos := (cache.PosReemplazo + i) % tamanio
			candidato := &cache.Entradas[pos]

			if !candidato.BitUso && candidato.Modificada {
				frameBase := candidato.NroPag * globals.ClientConfig.Page_size
				Write(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, frameBase, candidato.Contenido)
				log.Printf("PID: %d - Memory Update - Página: %d - Frame: %d", candidato.PID, candidato.NroPag, frameBase)

				*candidato = entradaNueva
				candidato.BitUso = true
				cache.PosReemplazo = (pos + 1) % tamanio
				return
			}
		}
		for i := 0; i < tamanio; i++ {
			cache.Entradas[i].BitUso = false
		}
	}


}