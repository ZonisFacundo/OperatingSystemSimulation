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
	var mu sync.Mutex

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

		http.HandleFunc("/InterrupcionCPU", func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			globals.Interruption = true
			mu.Unlock()
			log.Println("Interrupción recibida desde Kernel.")
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
			mu.Lock()
			interrumpido := globals.Interruption
			if interrumpido {
				globals.Interruption = false
			}
			mu.Unlock()

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
