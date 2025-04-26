package main

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instruction_cycle"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

// falta hacer la conexion del lado del cpu como servidor hacia el kernel pero no sabia donde hacerlas ni les queria romper el codigo =)
func main() {
	utilsCPU.ConfigurarLogger()
	globals.CargarConfig("./cpu/globals/config.json")

	/*log.Println("Ingrese el nombre de la instancia a ejecutar: ")

	fmt.Scanln(&globals.ClientConfig.Instance_id)*/

	instruction_cycle.SolicitudPIDyPCAKernel(globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Instance_id)
	instruction_cycle.Fetch(globals.Instruction.Pid, globals.Instruction.Pc, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	instruction_cycle.Execute(globals.ID)
}
