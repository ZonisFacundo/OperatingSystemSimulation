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

	go func() {
		http.HandleFunc("/INTERRUPCIONCPU", func(w http.ResponseWriter, r *http.Request) {
			utilsCPU.DevolverPidYPCInterrupcion(w, r, globals.ID.ProcessValues.Pc, globals.ID.ProcessValues.Pid)
			mutexInterrupcion.Lock()
			globals.Interruption = true
			mutexInterrupcion.Unlock()
			log.Println("## Llega interrupción al puerto Interrupt.") //OBLIGATORIO
		})

		log.Printf("Servidor HTTP activo en puerto %d.", globals.ClientConfig.Port_cpu)
		http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_cpu), nil)
	}()

	for {

		log.Println("## Esperando ingreso de un nuevo proceso.")

		<-globals.ProcesoNuevo

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
