package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instruction_cycle"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Error: Debe indicar el identificador de la CPU como argumento, por ejemplo: ./cpu cpuX")
	}

	instanceID := os.Args[1]
	procesoNuevo := make(chan struct{})
	var mutexInterrupcion sync.Mutex

	utilsCPU.ConfigurarLogger(instanceID)
	log.Printf("CPU %s inicializada correctamente.\n", instanceID)
	globals.CargarConfig("./cpu/globals/config.json", instanceID)

	utilsCPU.EnvioPortKernel(
		globals.ClientConfig.Ip_kernel,
		globals.ClientConfig.Port_kernel,
		globals.ClientConfig.Instance_id,
		globals.ClientConfig.Port_cpu,
		globals.ClientConfig.Ip_cpu,
	)
/* esto lo use para probar que tlb y cache funcionan bien 
	globals.Tlb.Tamanio = 4
	globals.Tlb.Entradas = make([]globals.Entrada, globals.Tlb.Tamanio)

	globals.CachePaginas.Tamanio = 4
	globals.CachePaginas.Entradas = make([]globals.EntradaCacheDePaginas, globals.CachePaginas.Tamanio)

	instruction_cycle.AgregarEnTLB(4, 4000)
	if !mmu.EstaTraducida(4) {
		log.Fatal("ERROR: la página 4 debería estar en la TLB")
	} else {
		log.Println("PRUEBA TLB OK")
	}

	instruction_cycle.AgregarEnCache(3, 3000)
	if !mmu.EstaEnCache(3) {
		log.Fatal("ERROR: la página 3 debería estar en la cache de páginas")
	} else {
		log.Println("PRUEBA CACHE DE PÁGINAS OK")
	}
*/
	go func() {
		http.HandleFunc("/KERNELCPU", func(w http.ResponseWriter, r *http.Request) {
			utilsCPU.RecibirPCyPID(w, r)
			log.Printf("Proceso recibido - PID: %d, PC: %d", globals.Instruction.Pid, globals.Instruction.Pc)
			select {
			case procesoNuevo <- struct{}{}:
				log.Println("Notificando CPU de un nuevo proceso entrante.")
			default:
				log.Println("CPU ya ejecutando. No se notifica de nuevo proceso")
			}
		})

		http.HandleFunc("/INTERRUPCIONCPU", func(w http.ResponseWriter, r *http.Request) {
			utilsCPU.DevolverPidYPCInterrupcion(w, r, globals.Instruction.Pc, globals.Instruction.Pid)
			mutexInterrupcion.Lock()
			globals.Interruption = true
			mutexInterrupcion.Unlock()
			log.Println("## Llega interrupción al puerto Interrupt.")
		})

		log.Printf("Servidor HTTP activo en puerto %d.", globals.ClientConfig.Port_cpu)
		http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_cpu), nil)
	}()

	for {
		log.Println("Esperando nuevo proceso...")

		<-procesoNuevo

		log.Printf(" Ejecutando proceso (PID: %d)", globals.Instruction.Pid)

	ejecucion:
		for {
			mutexInterrupcion.Lock()

			interrumpido := globals.Interruption
			if interrumpido {
				globals.Interruption = false
			}
			mutexInterrupcion.Unlock()

			if interrumpido {
				log.Printf("Interrupción. Deteniendo proceso PID %d", globals.Instruction.Pid)
				break ejecucion
			}

			log.Printf("Ejecutando: PID=%d, PC=%d", globals.Instruction.Pid, globals.Instruction.Pc)
			instruction_cycle.Fetch(globals.Instruction.Pid, globals.Instruction.Pc, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
			instruction_cycle.Decode(globals.ID)
			instruction_cycle.Execute(globals.ID)
			globals.Instruction.Pc++
		}
	}
}

/*
LOGS FALTANTES POR PONER:

Lectura/Escritura Memoria: “PID: <PID> - Acción: <LEER / ESCRIBIR> - Dirección Física: <DIRECCION_FISICA> - Valor: <VALOR LEIDO / ESCRITO>”.
Obtener Marco: “PID: <PID> - OBTENER MARCO - Página: <NUMERO_PAGINA> - Marco: <NUMERO_MARCO>”.
TLB Hit: “PID: <PID> - TLB HIT - Pagina: <NUMERO_PAGINA>”
TLB Miss: “PID: <PID> - TLB MISS - Pagina: <NUMERO_PAGINA>”
Página encontrada en Caché: “PID: <PID> - Cache Hit - Pagina: <NUMERO_PAGINA>”
Página faltante en Caché: “PID: <PID> - Cache Miss - Pagina: <NUMERO_PAGINA>”
Página ingresada en Caché: “PID: <PID> - Cache Add - Pagina: <NUMERO_PAGINA>”
Página Actualizada de Caché a Memoria: “PID: <PID> - Memory Update - Página: <NUMERO_PAGINA> - Frame: <FRAME_EN_MEMORIA_PRINCIPAL>”
*/
