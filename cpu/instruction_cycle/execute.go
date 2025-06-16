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
func Execute(detalle globals.Instruccion) {

	log.Printf("type: %s", globals.ID.InstructionType)
	log.Printf("value: %d", globals.ID.Valor)

	switch detalle.InstructionType {

	case "NOOP": //?
		/*if detalle.Tiempo != 0 {
			tiempoEjecucion := NOOP(detalle.Tiempo)
			detalle.ProcessValues.Pc = detalle.ProcessValues.Pc + 1
			fmt.Printf("NOOP ejecutado con tiempo:%d , y actualizado el PC:%d.\n", tiempoEjecucion, detalle.ProcessValues.Pc)
			log.Printf("## PID: %d - Ejecutando -> TYPE: %s ", detalle.ProcessValues.Pid, detalle.InstructionType)

		} else {
			fmt.Println("Tiempo no especificado u acción incorrecta.")
			detalle.Syscall = "Tiempo no especificado u acción incorrecta."
		}*/
		log.Println("Se ejecuta noop")

	case "WRITE":

		if globals.ID.DireccionFis != 0 { //Ésta habria que imprimir
			Write(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, globals.ID.DireccionFis, globals.ID.Datos)

			log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - DATOS: %s - DIRECCION: %d", detalle.ProcessValues.Pid, detalle.InstructionType, globals.ID.Datos, globals.ID.DireccionFis)
		} else {
			fmt.Println("WRITE inválido.")
			
			detalle.Syscall = "WRITE inválido."
		}

	case "READ":

		if globals.ID.DireccionFis != 0 {

			Read(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, globals.ID.DireccionFis)
			log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - SIZE: %d - DIRECCION: %d", detalle.ProcessValues.Pid, detalle.InstructionType, globals.ID.Tamaño, globals.ID.DireccionFis)

		} else {
			fmt.Sprintln("READ inválido.")

			detalle.Syscall = "READ inválido."
		}

	case "GOTO":

		pcInstrNew := GOTO(detalle.ProcessValues.Pc, detalle.Valor)

		fmt.Println("PC actualizado en: ", pcInstrNew)

		detalle.Syscall = fmt.Sprintf("PC actualizado en: %d ", pcInstrNew)
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - VALUE: %d", detalle.ProcessValues.Pid, detalle.InstructionType, pcInstrNew)

	// SYSCALLS.

	case "IO": //IO(Dispositivo y tiempo)
		log.Printf("Se envia syscall: %s a Kernel", detalle.InstructionType)
		utilsCPU.FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.Instruction.Pid,
			globals.Instruction.Pc,
			globals.ClientConfig.Instance_id,
			globals.InstruccionDetalle.InstructionType,
			globals.InstruccionDetalle.Parametro1,
			globals.InstruccionDetalle.Parametro2)

	case "INIT_PROC": //INIT_PROC (Archivo de instrucciones, Tamaño)
		log.Printf("Se envia syscall: %s a Kernel", detalle.InstructionType)
		utilsCPU.FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.Instruction.Pid, globals.Instruction.Pc,
			globals.ClientConfig.Instance_id,
			globals.InstruccionDetalle.Syscall,
			globals.InstruccionDetalle.Parametro1,
			globals.InstruccionDetalle.Parametro2)

	case "DUMP_MEMORY": //
		log.Printf("Se envia syscall: %s a Kernel", detalle.InstructionType)
		utilsCPU.FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.Instruction.Pid,
			globals.Instruction.Pc,
			globals.ClientConfig.Instance_id,
			globals.InstruccionDetalle.Syscall,
			0,
			"")

	case "EXIT":
		log.Printf("Se envia syscall: %s a Kernel", detalle.InstructionType)
		utilsCPU.FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.Instruction.Pid,
			globals.Instruction.Pc,
			globals.ClientConfig.Instance_id,
			globals.InstruccionDetalle.Syscall,
			0,
			"")

	default:
		fmt.Println("Instrucción inválida.")
	}
}

func Write(ip string, port int, direccion int, contenido string) {

	var paquete utilsCPU.WriteStruct

	paquete.Contenido = contenido
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

	log.Printf("Respuesta de Memoria: %d\n", respuesta.Mensaje) // Nos devuelve memoria el mensaje de escritura.

	globals.ID.DireccionLog = respuesta.Mensaje
}

func Read(ip string, port int, direccion int) {

	var paquete utilsCPU.ReadStruct

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

	log.Printf("Valor en memoria: %s\n", respuesta.Info) // Nos devuelve memoria el mensaje de escritura.

}
