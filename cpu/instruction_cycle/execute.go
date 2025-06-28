package instruction_cycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

// switch para ver que hace dependiendo la instruccion:
func Execute(detalle globals.Instruccion) bool {

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

		if globals.ID.DireccionFis >= 0 {
			if globals.ClientConfig.Cache_entries > 0 { // caché habilitada
				if mmu.EstaEnCache(globals.ID.NroPag) {
					mmu.WriteEnCache(globals.ID.Datos)
				} else {
					// leer en memoria y traer la página a la caché
					Write(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, globals.ID.DireccionFis, globals.ID.Datos)
					AgregarEnTLB(globals.ID.NroPag, globals.ID.DireccionFis)
					AgregarEnCache(globals.ID.NroPag, globals.ID.DireccionFis)
				}
			} else {
				// caché deshabilitada
				Write(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, globals.ID.DireccionFis, globals.ID.Datos)
				AgregarEnTLB(globals.ID.NroPag, globals.ID.DireccionFis)
			}

		} else {
			fmt.Println("WRITE inválido.")
			detalle.Syscall = "WRITE inválido."
		}

		return false

	case "READ":

		if globals.ID.DireccionFis >= 0 {

			if globals.ClientConfig.Cache_entries > 0 { // caché habilitada
				if mmu.EstaEnCache(globals.ID.NroPag) {
					mmu.ReadEnCache()
				} else {
					Read(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, globals.ID.DireccionFis, globals.ID.Tamaño)
					AgregarEnTLB(globals.ID.NroPag, globals.ID.DireccionFis)
					AgregarEnCache(globals.ID.NroPag, globals.ID.DireccionFis)
				}
			} else {
				// caché deshabilitada
				Read(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, globals.ID.DireccionFis, globals.ID.Tamaño)
				AgregarEnTLB(globals.ID.NroPag, globals.ID.DireccionFis)
			}

			log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - SIZE: %d - DIRECCION: %d",
				detalle.ProcessValues.Pid, detalle.InstructionType, globals.ID.Tamaño, globals.ID.DireccionFis)

		} else {
			fmt.Sprintln("READ inválido.")
			detalle.Syscall = "READ inválido."
		}
		return false

	case "GOTO":

		pcInstrNew := GOTO(detalle.ProcessValues.Pc, detalle.Valor)

		fmt.Println("PC actualizado en: ", pcInstrNew)

		detalle.Syscall = fmt.Sprintf("PC actualizado en: %d ", pcInstrNew)
		globals.ID.ProcessValues.Pc = pcInstrNew
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - VALUE: %d", detalle.ProcessValues.Pid, detalle.InstructionType, globals.ID.ProcessValues.Pc)
		return false

	// SYSCALLS.

	case "IO": //IO(Dispositivo y tiempo)
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - DISPOSITIVO: %s - TIME: %d", detalle.ProcessValues.Pid, detalle.InstructionType, detalle.Dispositivo, detalle.Tiempo)
		FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.Instruction.Pid,
			globals.Instruction.Pc,
			globals.ClientConfig.Instance_id,
			detalle.InstructionType,
			globals.InstruccionDetalle.Tiempo,
			globals.InstruccionDetalle.Dispositivo)

		return true

	case "INIT_PROC": //INIT_PROC (Archivo de instrucciones, Tamaño)
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - TAM: %d - ARCHIVO: %s", detalle.ProcessValues.Pid, detalle.InstructionType, detalle.Tamaño, detalle.ArchiInstr)

		FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.Instruction.Pid, globals.Instruction.Pc,
			globals.ClientConfig.Instance_id,
			detalle.InstructionType,
			globals.ID.Tamaño,
			globals.ID.ArchiInstr)

		return true

	case "DUMP_MEMORY": //
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s", detalle.ProcessValues.Pid, detalle.InstructionType)
		FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.Instruction.Pid,
			globals.Instruction.Pc,
			globals.ClientConfig.Instance_id,
			detalle.InstructionType,
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
			detalle.InstructionType,
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

	log.Printf("syscall enviada: %s", paquete.Syscall)

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

// agregar en cache segun algortimo

func AgregarEnCache(nroPagina int, direccionFisica int) {
	entrada := globals.EntradaCacheDePaginas{
		NroPag:          nroPagina,
		Contenido:       globals.ID.Datos,
		DireccionFisica: direccionFisica,
		Modificada:      false,
		BitUso:          true,
	}

	if len(globals.CachePaginas.Entradas) < globals.CachePaginas.Tamanio {
		globals.CachePaginas.Entradas = append(globals.CachePaginas.Entradas, entrada)
		return
	}
	switch globals.AlgoritmoReemplazo {
	case "CLOCK":
		ReemplazarConCLOCK(entrada)
	case "CLOCK-M":
		ReemplazarConCLOCKM(entrada)
	default:
		log.Printf("Algoritmo incorrecto.\n")
	}

}

func AgregarEnTLB(nroPagina int, direccion int) {
	entrada := globals.Entrada{
		NroPagina: nroPagina,
		Direccion: direccion,
	}

	if len(globals.Tlb.Entradas) < globals.Tlb.Tamanio {
		globals.Tlb.Entradas = append(globals.Tlb.Entradas, entrada)
		return
	} else {
		switch globals.AlgoritmoReemplazoTLB {
		case "FIFO":
			ReemplazarTLB_FIFO(entrada)
		case "LRU":
			ReemplazarTLB_LRU(entrada)
		default:
			log.Printf("Algoritmo incorrecto.\n")
		}
	}
}
