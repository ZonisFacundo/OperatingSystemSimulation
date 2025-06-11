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

// falta hacer la conexion del lado del cpu como servidor hacia el kernel pero no sabia donde hacerlas ni les queria romper el codigo =)
func main() {

	var wg sync.WaitGroup
	procesoListo := make(chan bool)

	if len(os.Args) < 2 {
		log.Fatal("Error: Debe indicar el identificador de la CPU como argumento, por ejemplo: ./cpu cpuX")
	}

	instanceID := os.Args[1]

	utilsCPU.ConfigurarLogger(instanceID)
	log.Printf("CPU %s inicializada correctamente.\n", instanceID)
	globals.CargarConfig("./cpu/globals/config.json", instanceID)

	utilsCPU.EnvioPortKernel(globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Instance_id, globals.ClientConfig.Port_cpu, globals.ClientConfig.Ip_cpu)

	http.HandleFunc("/KERNELCPU", func(w http.ResponseWriter, r *http.Request) {
		utilsCPU.RecibirPCyPID(w, r)
		procesoListo <- true
	})

	go func() {
		log.Printf("Servidor corriendo, esperando PID y PC de Kernel.")
		err := http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_cpu), nil)
		if err != nil {
			log.Fatalf("Error al iniciar servidor HTTP: %v", err)
		}
	}()

	// Goroutine que espera el proceso y ejecuta el ciclo de instrucción
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			<-procesoListo // espera a que llegue un procesos (habria que hacer algo para que no espere eternamente)
			log.Printf("Iniciando ciclo de instrucción para PID: %d, PC: %d", globals.Instruction.Pid, globals.Instruction.Pc)
			instruction_cycle.Fetch(globals.Instruction.Pid, globals.Instruction.Pc, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
			instruction_cycle.Decode(globals.ID)
			instruction_cycle.Execute(globals.ID)
		}
	}()

	wg.Wait()
}