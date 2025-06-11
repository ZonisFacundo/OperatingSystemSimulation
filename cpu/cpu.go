package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instruction_cycle"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

// falta hacer la conexion del lado del cpu como servidor hacia el kernel pero no sabia donde hacerlas ni les queria romper el codigo =)

func main() {
	newFetch := true
	interrupcionActiva := false

	if len(os.Args) < 2 {
		log.Fatal("Error: Debe indicar el identificador de la CPU como argumento, por ejemplo: ./cpu cpuX")
	}

	instanceID := os.Args[1]
	procesoListo := make(chan struct{}, 1)

	utilsCPU.ConfigurarLogger(instanceID)
	log.Printf("CPU %s inicializada correctamente.\n", instanceID)
	globals.CargarConfig("./cpu/globals/config.json", instanceID)

	utilsCPU.EnvioPortKernel(globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Instance_id, globals.ClientConfig.Port_cpu, globals.ClientConfig.Ip_cpu)

	go func() {
		http.HandleFunc("/KERNELCPU", utilsCPU.RecibirPCyPID)
		procesoListo <- struct{}{}
		log.Printf("Servidor corriendo, esperando PID y PC de Kernel.")
		http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_cpu), nil)
	}()
	// Ciclo principal

	for {
		if interrupcionActiva {

			log.Println("Interrupción ejecutada, a la espera de nuevo proceso.")
			newFetch = true // Ésto se realiza más que nada porque cuando hay una interrupción, se interrumpe la ejecución del proceso y nos van a mandar uno nuevo.
			<-procesoListo
			log.Println("Nuevo proceso recibido, se reinicia el ciclo.")
		}

		if newFetch {
			instruction_cycle.Fetch(globals.Instruction.Pid, globals.Instruction.Pc, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
			newFetch = false
		}

		instruction_cycle.Decode(globals.ID)
		instruction_cycle.Execute(globals.ID)

		interrupcionActiva = instruction_cycle.CheckInterruption()

		time.Sleep(100 * time.Millisecond)
	}
}
