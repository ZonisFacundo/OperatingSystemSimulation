package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instruction_cycle"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
	//"github.com/sisoputnfrba/tp-golang/cpu/mmu"
)

// falta hacer la conexion del lado del cpu como servidor hacia el kernel pero no sabia donde hacerlas ni les queria romper el codigo =)
func main() {

	if len(os.Args) < 2 {
		log.Fatal("Error: Debe indicar el identificador de la CPU como argumento, por ejemplo: ./cpu cpuX")
	}

	instanceID := os.Args[1]

	/* Esto fue para probar la traduccion
	mmU := mmu.MMU{
		Niveles:             2,
		TamPagina:           256,
		Cant_entradas_tabla: 4,
	}

	dirLogica := 1800
	resultado := mmu.TraducirDireccion(dirLogica, mmU)
	*/

	utilsCPU.ConfigurarLogger(instanceID)
	log.Printf("CPU %s inicializada correctamente.\n", instanceID)
	globals.CargarConfig("./cpu/globals/config.json", instanceID)

	utilsCPU.EnvioPortKernel(globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Instance_id)

	//fmt.Println("Entradas + Desplazamiento:", resultado)

	http.HandleFunc("/KERNELCPU", utilsCPU.RecibirPCyPID)
	log.Printf("Servidor corriendo, esperando PID y PC de Kernel.")
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.Port_cpu), nil)

	instruction_cycle.Fetch(globals.Instruction.Pid, globals.Instruction.Pc, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	instruction_cycle.Execute(globals.ID)

	utilsCPU.FinEjecucion(globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.Instruction.Pid, globals.Instruction.Pc, globals.ClientConfig.Instance_id, globals.InstruccionDetalle.Contexto)

}
