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
func Execute(detalle globals.Instruccion) bool {

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

		return false

	case "READ":

		if globals.ID.DireccionFis != 0 {

			Read(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, globals.ID.DireccionFis, globals.ID.Parametro1)
			log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - SIZE: %d - DIRECCION: %d", detalle.ProcessValues.Pid, detalle.InstructionType, globals.ID.Tamaño, globals.ID.DireccionFis)

		} else {
			fmt.Sprintln("READ inválido.")

			detalle.Syscall = "READ inválido."
		}

		return false

	case "GOTO":

		pcInstrNew := GOTO(detalle.ProcessValues.Pc, detalle.Valor)

		fmt.Println("PC actualizado en: ", pcInstrNew)

		detalle.Syscall = fmt.Sprintf("PC actualizado en: %d ", pcInstrNew)
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - VALUE: %d", detalle.ProcessValues.Pid, detalle.InstructionType, pcInstrNew)

		return false

	// SYSCALLS.

	case "IO": //IO(Dispositivo y tiempo)
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - DISPOSITIVO: %d - TIME: %s", detalle.ProcessValues.Pid, detalle.InstructionType, detalle.Parametro1, detalle.Parametro2)
		FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.Instruction.Pid,
			globals.Instruction.Pc,
			globals.ClientConfig.Instance_id,
			globals.InstruccionDetalle.InstructionType,
			globals.InstruccionDetalle.Parametro1,
			globals.InstruccionDetalle.Parametro2)

		return true

	case "INIT_PROC": //INIT_PROC (Archivo de instrucciones, Tamaño)
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - TAM: %d - ARCHIVO: %s", detalle.ProcessValues.Pid, detalle.InstructionType, detalle.Parametro1, detalle.Parametro2)
		FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.Instruction.Pid, globals.Instruction.Pc,
			globals.ClientConfig.Instance_id,
			globals.InstruccionDetalle.Syscall,
			globals.InstruccionDetalle.Parametro1,
			globals.InstruccionDetalle.Parametro2)

		return true

	case "DUMP_MEMORY": //
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s", detalle.ProcessValues.Pid, detalle.InstructionType)
		FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.Instruction.Pid,
			globals.Instruction.Pc,
			globals.ClientConfig.Instance_id,
			globals.InstruccionDetalle.Syscall,
			0,
			"")
		return true

	case "EXIT":
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s", detalle.ProcessValues.Pid, detalle.InstructionType)
		FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.Instruction.Pid,
			globals.Instruction.Pc,
			globals.ClientConfig.Instance_id,
			globals.InstruccionDetalle.Syscall,
			0,
			"")
		return true

	default:
		fmt.Println("Instrucción inválida.")
		return false
	}

	return false
}

func Write(ip string, port int, direccion int, contenido string) {

	var paquete utilsCPU.WriteStruct

	paquete.Contenido = contenido
	paquete.Direccion = direccion

	log.Printf("dir: %d, cont: %s", paquete.Direccion, paquete.Contenido)

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

	paquete.Direccion = direccion
	paquete.Tamaño = tamaño

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

	informacion := string(respuesta.Info)

	log.Printf("Valor en memoria: [%s]", informacion) // Nos devuelve memoria el mensaje de escritura.

}

func FinEjecucion(ip string, puerto int, pid int, pc int, instancia string, syscall string, parametro1 int, parametro2 string) { // si no reciben parametros que sean  0 y "" que nosostros ahi no los usamos
	var paquete utilsCPU.PackageFinEjecucion

	paquete.Pid = pid
	paquete.Pc = pc
	paquete.Syscall = syscall
	paquete.InstanciaCPU = instancia
	paquete.Parametro1 = parametro1
	paquete.Parametro2 = parametro2

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.\n")
		return
	}

	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/PCB", ip, puerto) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

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

	log.Printf("Se envio syscall a Kernel.\n")
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta utilsCPU.RespuestaKernel

	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
	}

	if respuesta.Mensaje != "interrupcion" {
		globals.Interruption = true

	} else {
		log.Printf("El Kernel recibió correctamente el PID y el PC.\n")
	}
}
