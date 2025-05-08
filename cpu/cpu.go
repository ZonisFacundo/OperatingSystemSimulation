package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instruction_cycle"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

// falta hacer la conexion del lado del cpu como servidor hacia el kernel pero no sabia donde hacerlas ni les queria romper el codigo =)
func main() {

	var wg sync.WaitGroup

	if len(os.Args) < 2 {
		log.Fatal("Error: Debe indicar el identificador de la CPU como argumento, por ejemplo: ./cpu cpuX")
	}

	instanceID := os.Args[1]

	// Esto fue para probar la traduccion
	mmU := mmu.MMU{
		Niveles: 2,
		// TamPagina:           256,
		Cant_entradas_tabla: 4,
	}

	dirLogica := 1800
	resultado := mmu.TraducirDireccion(dirLogica, mmU, globals.Instruction.Pid)

	utilsCPU.ConfigurarLogger(instanceID)
	log.Printf("CPU %s inicializada correctamente.\n", instanceID)
	globals.CargarConfig("./cpu/globals/config.json", instanceID)

	utilsCPU.EnvioPortKernel(globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Instance_id, globals.ClientConfig.Port_cpu, globals.ClientConfig.Ip_cpu)

	fmt.Println("Entradas + Desplazamiento:", resultado)

	http.HandleFunc("/KERNELCPU", utilsCPU.RecibirPCyPID)
	log.Printf("Servidor corriendo, esperando PID y PC de Kernel.")
	go http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_cpu), nil)

	instruction_cycle.Fetch(globals.Instruction.Pid, globals.Instruction.Pc, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	instruction_cycle.Decode(globals.ID)
	instruction_cycle.Execute(globals.ID)

	utilsCPU.FinEjecucion(globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.Instruction.Pid, globals.Instruction.Pc, globals.ClientConfig.Instance_id, globals.InstruccionDetalle.Syscall, 0, "")

	wg.Wait()

}
