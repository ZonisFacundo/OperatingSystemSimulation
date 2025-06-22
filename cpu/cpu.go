package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instruction_cycle"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

func main() {

	//interrupcionActiva := false

	if len(os.Args) < 2 {
		log.Fatal("Error: Debe indicar el identificador de la CPU como argumento, por ejemplo: ./cpu cpuX")
	}

	instanceID := os.Args[1]
	procesoListo := make(chan struct{}, 0)

	utilsCPU.ConfigurarLogger(instanceID)
	log.Printf("CPU %s inicializada correctamente.\n", instanceID)
	globals.CargarConfig("./cpu/globals/config.json", instanceID)

	utilsCPU.EnvioPortKernel(globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Instance_id, globals.ClientConfig.Port_cpu, globals.ClientConfig.Ip_cpu)

	go func() {
		http.HandleFunc("/KERNELCPU", func(w http.ResponseWriter, r *http.Request) {
			utilsCPU.RecibirPCyPID(w, r)
			procesoListo <- struct{}{}
		})
		/*http.HandleFunc("/InterrupcionCPU", func(w http.ResponseWriter, r *http.Request) {
			globals.Interruption = true
		})*/
		log.Printf("Servidor corriendo, esperando PID y PC de Kernel.")
		http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_cpu), nil)
	}()

	// Ciclo principal

	for {

		/*if interrupcionActiva {

			log.Println("Interrupción ejecutada, a la espera de nuevo proceso.")
			globals.Interruption = false
			<-procesoListo
			log.Println("Nuevo proceso recibido, se reinicia el ciclo.")

		}*/
		<-procesoListo
		instruction_cycle.Fetch(globals.Instruction.Pid, globals.Instruction.Pc, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
		instruction_cycle.Decode(globals.ID)
		instruction_cycle.Execute(globals.ID)

		globals.Instruction.Pc++
		//interrupcionActiva = globals.Interruption
	}
}
