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

	if len(os.Args) < 2 {
		log.Fatal("Error: Debe indicar el identificador de la CPU como argumento, por ejemplo: ./cpu cpuX")
	}

	instanceID := os.Args[1]

	utilsCPU.ConfigurarLogger(instanceID)
	log.Printf("CPU %s inicializada correctamente.\n", instanceID)
	globals.CargarConfig("./cpu/globals/config.json", instanceID)

	utilsCPU.EnvioPortKernel(globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Instance_id, globals.ClientConfig.Port_cpu, globals.ClientConfig.Ip_cpu)

	go func() {
		http.HandleFunc("/KERNELCPU", utilsCPU.RecibirPCyPID)
		log.Printf("Servidor corriendo, esperando PID y PC de Kernel.")
		http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_cpu), nil)
	}()

	primeraEjecucion := true

	// Ciclo principal
	for {
		if primeraEjecucion {
			instruction_cycle.Fetch(globals.Instruction.Pid, globals.Instruction.Pc, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
			primeraEjecucion = false
		}

		instruction_cycle.Decode(globals.ID)
		instruction_cycle.Execute(globals.ID)

		time.Sleep(100 * time.Millisecond)
	}
}
