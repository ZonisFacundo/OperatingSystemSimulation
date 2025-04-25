package instruction_cycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

// switch para ver que hace dependiendo la instruccion:
func InstruccionDetalle(detalle globals.Instruccion) {

	switch detalle.InstructionType {
	case "NOOP":
		if detalle.Tiempo != nil {
			tiempoEjecucion := Noop(*detalle.Tiempo)
			detalle.ProcessValues.Pc = detalle.ProcessValues.Pc + 1
			fmt.Printf("NOOP ejecutado con tiempo:%d , y actualizado el PC:%d.\n", tiempoEjecucion, detalle.ProcessValues.Pc)

			//acá voy a tener que actualizar el PC, ¿cómo? ni idea.
		} else {
			fmt.Println("Tiempo no especificado u acción incorrecta.")
		}

	case "WRITE":
		if detalle.Direccion != 0 || detalle.Datos != nil {
			datosACopiar := detalle.Datos
			direccionObtenida := detalle.Direccion

			/*TraducirDireccion(direccionObtenida)*/
			Write(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, direccionObtenida, *datosACopiar)
		} else {
			fmt.Println("WRITE inválido.")
		}

	case "READ":
		if detalle.Direccion != 0 || detalle.Tamaño != nil {
			/*tamañoDet := detalle.Tamaño
			direccionObtenida := detalle.Direccion

			TraducirDireccion(direccionObtenida)*/
			Read(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, detalle.Direccion, *detalle.Tamaño)

		} else {
			fmt.Println("READ inválido.")
		}

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

func Write(ip string, port int, direccion int, datos string) {

	var paquete utilsCPU.WriteStruct

	paquete.Datos = datos
	paquete.Direccion = direccion

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.\n")
		return
	}
	cliente := http.Client{}

	// log.Printf("Process ID: %d - FETCH - Program Counter: %d.\n", pid, pc) // log mínimo y obligatorio

	url := fmt.Sprintf("http://%s:%d/WRITE", ip, port)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson))

	if err != nil {
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return

	}

	if respuestaJSON.StatusCode != http.StatusOK {

		log.Printf("Status de respuesta el server no fue la esperada.\n")
		return
	}
	defer respuestaJSON.Body.Close()

	fmt.Printf("Conexion establecida con exito.\n")
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta PaqueteRecibidoWRITE

	err = json.Unmarshal(body, &(respuesta))
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}

	log.Printf("Respuesta de Memoria: %s\n", respuesta.Mensaje) // Nos devuelve memoria el mensaje de escritura.

}

func Read(ip string, port int, direccion int, tamaño int) {

	var paquete utilsCPU.ReadStruct

	paquete.Tamaño = tamaño
	paquete.Direccion = direccion

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.\n")
		return
	}
	cliente := http.Client{}

	// log.Printf("Process ID: %d - FETCH - Program Counter: %d.\n", pid, pc) // log mínimo y obligatorio

	url := fmt.Sprintf("http://%s:%d/READ", ip, port)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson))

	if err != nil {
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return

	}

	if respuestaJSON.StatusCode != http.StatusOK {

		log.Printf("Status de respuesta el server no fue la esperada.\n")
		return
	}
	defer respuestaJSON.Body.Close()

	fmt.Printf("Conexion establecida con exito.\n")
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta PaqueteRecibidoREAD

	err = json.Unmarshal(body, &(respuesta))
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}

	log.Printf("Valor en memoria: %s\n", respuesta.ValorInstruccion) // Nos devuelve memoria el mensaje de escritura.

}
