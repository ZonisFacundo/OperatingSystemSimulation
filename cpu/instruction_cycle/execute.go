package instruction_cycle

import (
	"fmt"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

// switch para ver que hace dependiendo la instruccion:
func InstruccionDetalle(detalle globals.Instruccion) {

	switch detalle.InstructionType {
	case "NOOP\n":
		if detalle.Tiempo != nil {
			tiempoEjecucion := Noop(*detalle.Tiempo)
			detalle.ProcessValues.Pc = detalle.ProcessValues.Pc + 1
			fmt.Printf("NOOP ejecutado con tiempo:%d , y actualizado el PC:%d.\n", tiempoEjecucion, detalle.ProcessValues.Pc)

			//acá voy a tener que actualizar el PC, ¿cómo? ni idea.
		} else {
			fmt.Println("Tiempo no especificado u acción incorrecta.")
		}
	case "WRITE":
	case "READ":
	case "GOTO":
		if detalle.Valor != nil {
			pcInstrNew := GOTO(detalle.ProcessValues.Pc, *detalle.Valor)
			fmt.Println("PC actualizado en: ", pcInstrNew)
		} else {
			fmt.Println("Valor no modificado.")
		}
		// LLamada a Kernel, debido a que son parte principalmente de interrupciones.
	case "IO":
	case "INIT_PROC":
	case "DUMP_MEMORY":
	case "EXIT":
		fmt.Println("Nada que hacer.")
		return
	default:
		fmt.Println("Instrucción inválida.")
	}

}

func Noop(Tiempo int) int {
	return Tiempo
}

func Write(direccion int, datos string) {

}

func Read(direccion int, tamaño *int) {
	if direccion == 0 || tamaño == nil {
		fmt.Println("READ mal formada")
	}

}

func GOTO(pcInstr int, valor int) int {
	return pcInstr + valor
}

/*

Write(direccion, datos){
	escribe los datos en la direccion especifica, primero voy a tener que traducir la dir. lógica y luego,
	voy a tener que acceder a esa dirección (No se como) y voy a tener que escribir esos datos
	en esa dirección// datos string
	+1 PC
}

Read(direccion, tamaño){
	printf(direccion,direccion.tamaño) //Lee la dirección , e imprime en pantalla el tamaño de esa dirección con log obligatorio
	+1 PC
}
    case "READ":
        if instr.Direccion == 0 || instr.Tamaño == nil {
            fmt.Println("READ mal formada")
            return
        }
        dirFisica := cpu.TraducirDireccion(instr.Direccion)
        datos := cpu.Memoria.Leer(dirFisica, *instr.Tamaño)
        fmt.Println("READ:", datos)
        kernel.LoggearLectura(cpu.PID, datos)


GOTO(valor){
	pc = pc + valor //Actualiza el valor del PC sumandole el valor indicado.
}


// instrucciones que realiza kernel, la cpu no puede ejecutarlaS
IO(tiempo){
	... ¿interrupcion?
	pc ++
}
INIT_PROC(archivoInstr, tamaño){
	... "la hace kernel"
	pc ++
}
DUMP_MEMORY(){
	retornarPIDAKernel(detalle.pid)
	pc ++
}
EXIT(){
	...
	pc ++
}

Las siguientes instrucciones se considerarán Syscalls, ya que las mismas no pueden ser resueltas por la CPU y
depende de la acción del Kernel para su realización, a diferencia de la vida real donde la llamada es a una única instrucción,
para simplificar la comprensión de los scripts, vamos a utilizar un nombre diferente para cada Syscall.



*/
