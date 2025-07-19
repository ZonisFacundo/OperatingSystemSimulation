package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
	configSuffix := instanceID

	if strings.HasPrefix(instanceID, "cpu") {
		configSuffix = instanceID[3:] // Quita "cpu"
	}

	configPath := fmt.Sprintf("./cpu/globals/config%s.json", configSuffix)

	log.Printf("Usando config: %s para instancia: %s\n", configPath, instanceID)

	globals.ProcesoNuevo = make(chan struct{}, 1)
	var mutexInterrupcion sync.Mutex

	utilsCPU.ConfigurarLogger(instanceID)
	log.Printf("CPU %s inicializada correctamente.\n", instanceID)
	globals.CargarConfig(configPath, instanceID)

	globals.AlgoritmoReemplazo = globals.ClientConfig.Cache_replacement
	globals.AlgoritmoReemplazoTLB = globals.ClientConfig.Tlb_replacement
	globals.InitTlb()
	globals.InitCache()

	instruction_cycle.RecibirDatosMMU(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)

	utilsCPU.EnvioPortKernel(
		globals.ClientConfig.Ip_kernel,
		globals.ClientConfig.Port_kernel,
		globals.ClientConfig.Instance_id,
		globals.ClientConfig.Port_cpu,
		globals.ClientConfig.Ip_cpu,
	)
	http.HandleFunc("/KERNELCPU", instruction_cycle.RecibirPCyPID)
	/*http.HandleFunc("/KERNELCPU", func(w http.ResponseWriter, r *http.Request) {
		globals.MutexNecesario.Lock()
		instruction_cycle.RecibirPCyPID(w, r)
		globals.MutexNecesario.Unlock()

		log.Printf("Proceso recibido - PID: %d, PC: %d", globals.ID.ProcessValues.Pid, globals.ID.ProcessValues.Pc)

		select {
		case procesoNuevo <- struct{}{}:
			log.Println("## Se notifica a CPU de un proceso entrante.")
		default:
			log.Println("## CPU ya ejecutando. No se notifica de nuevo proceso.")
		}
	})*/

	go func() {
		http.HandleFunc("/INTERRUPCIONCPU", func(w http.ResponseWriter, r *http.Request) {
			utilsCPU.DevolverPidYPCInterrupcion(w, r, globals.ID.ProcessValues.Pc, globals.ID.ProcessValues.Pid)
			mutexInterrupcion.Lock()
			globals.Interruption = true
			mutexInterrupcion.Unlock()
			log.Println("## Llega interrupción al puerto Interrupt.")
		})

		log.Printf("Servidor HTTP activo en puerto %d.", globals.ClientConfig.Port_cpu)
		http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_cpu), nil)
	}()

	for {
		log.Println("## Esperando ingreso de un nuevo proceso.")

		<-globals.ProcesoNuevo

		log.Printf("## Ejecutando proceso (PID: %d)", globals.ID.ProcessValues.Pid)

	ejecucion:
		for {
			mutexInterrupcion.Lock()

			interrumpido := globals.Interruption

			if interrumpido {
				globals.Interruption = false
			}
			mutexInterrupcion.Unlock()

			if interrumpido {
				log.Printf("## Interrupcion recibida -> Deteniendo proceso con PID: %d", globals.ID.ProcessValues.Pid)
				instruction_cycle.VaciarCache(globals.ID.ProcessValues.Pid)
				log.Printf("CHUPETE EN EL ORTO INSIDE")
				break ejecucion
			}

			globals.MutexNecesario.Lock()
			pid := globals.ID.ProcessValues.Pid
			pc := globals.ID.ProcessValues.Pc
			globals.MutexNecesario.Unlock()

			log.Printf("## Ejecutando -> PID: %d, PC: %d", pid, pc)

			instruction_cycle.Fetch(globals.ID.ProcessValues.Pid, globals.ID.ProcessValues.Pc, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
			instruction_cycle.Decode(globals.ID)

			instruction_cycle.Execute(globals.ID)
			globals.ID.ProcessValues.Pc++
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
