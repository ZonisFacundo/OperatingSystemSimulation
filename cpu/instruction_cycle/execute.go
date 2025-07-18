package instruction_cycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

// switch para ver que hace dependiendo la instruccion:
func Execute(detalle globals.Instruccion) bool {

	switch detalle.InstructionType {

	case "NOOP":
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s ", detalle.ProcessValues.Pid, detalle.InstructionType)
		return false

	case "WRITE":
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - DIRECCION: %d DATOS: %s", detalle.ProcessValues.Pid, detalle.InstructionType, globals.ID.DireccionLog, globals.ID.Datos)
		if globals.ID.DireccionFis >= 0 {
			if globals.ClientConfig.Cache_entries > 0 { //cache esta habilitada (está vacia?)
				if mmu.EstaEnCache(globals.ID.NroPag) {
					mmu.WriteEnCache(globals.ID.ProcessValues.Pid, globals.ID.NroPag, globals.ID.Desplazamiento, []byte(globals.ID.Datos))
					//log.Printf("## WRITE en Cache: PID: %d, Pag: %d, Datos: %s", globals.ID.ProcessValues.Pid, globals.ID.NroPag, globals.ID.Datos)
				} else {

					// leer en memoria y traer la página a la caché
					Write(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, globals.ID.DireccionFis, globals.ID.Datos)
					log.Printf("## PID: %d - Accion: ESCRIBIR - Direccion Física: %d - Valor: %v",
						globals.ID.ProcessValues.Pid, globals.ID.DireccionFis, globals.ID.Datos)
					pagina := make([]byte, globals.ClientConfig.Page_size)
					copy(pagina[globals.ID.Desplazamiento:], []byte(globals.ID.Datos))
					globals.ID.PaginaCompleta = pagina
					AgregarEnTLB(globals.ID.NroPag, globals.ID.DireccionFis)
					AgregarEnCache(globals.ID.NroPag, globals.ID.DireccionFis)
					log.Printf("PID: %d - Cache Add - Pagina: %d", globals.ID.ProcessValues.Pid, globals.ID.NroPag)
				}
			} else {
				Write(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, globals.ID.DireccionFis, globals.ID.Datos)
				AgregarEnTLB(globals.ID.NroPag, globals.ID.DireccionFis)
				log.Printf("## PID: %d - Accion: ESCRIBIR - Direccion Física: %d - Valor: %v",
					globals.ID.ProcessValues.Pid, globals.ID.DireccionFis, globals.ID.Datos)
			}

		} else {
			fmt.Println("## ERROR -> WRITE inválido: Direccion fisica inválida.")
			detalle.Syscall = "WRITE inválido."
		}
		return false

	case "READ":
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - DIRECCION: %d TAMAÑO: %d", detalle.ProcessValues.Pid, detalle.InstructionType, globals.ID.DireccionLog, globals.ID.Tamaño)
		if globals.ID.DireccionFis >= 0 {

			if globals.ClientConfig.Cache_entries > 0 {
				if mmu.EstaEnCache(globals.ID.NroPag) {
					mmu.ReadEnCache()
				} else {
					Read(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, globals.ID.DireccionFis, globals.ID.Tamaño)
					AgregarEnTLB(globals.ID.NroPag, globals.ID.DireccionFis)
					AgregarEnCache(globals.ID.NroPag, globals.ID.DireccionFis)
					log.Printf("PID: %d - Cache Add - Pagina: %d", globals.ID.ProcessValues.Pid, globals.ID.NroPag)
					log.Printf("## PID: %d - Accion: LEER - Direccion Física: %d - Valor: %v",
						globals.ID.ProcessValues.Pid, globals.ID.DireccionFis, globals.ID.ValorLeido)
				}
			} else {

				Read(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, globals.ID.DireccionFis, globals.ID.Tamaño)
				AgregarEnTLB(globals.ID.NroPag, globals.ID.DireccionFis)
				log.Printf("## PID: %d - Accion: LEER - Direccion Física: %d - Valor: %v",
					globals.ID.ProcessValues.Pid, globals.ID.DireccionFis, globals.ID.ValorLeido)
			}

		} else {
			fmt.Sprintln("## ERROR -> READ invalido.")
			detalle.Syscall = "READ invalido."
		}
		return false

	case "GOTO":

		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - VALUE: %d", detalle.ProcessValues.Pid, detalle.InstructionType, globals.ID.ProcessValues.Pc)
		pcInstrNew := detalle.Valor
		//fmt.Println("## EJECUCIÓN -> GOTO - PC actualizado en: ", pcInstrNew)
		detalle.Syscall = fmt.Sprintf("## PC actualizado en: %d ", pcInstrNew)

		globals.ID.ProcessValues.Pc = pcInstrNew

		globals.ID.ProcessValues.Pc--
		return false

	// SYSCALLS.

	case "IO": //IO(Dispositivo y tiempo)
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - DISPOSITIVO: %s - TIME: %d", detalle.ProcessValues.Pid, detalle.InstructionType, detalle.Dispositivo, detalle.Tiempo)
		FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.ID.ProcessValues.Pid,
			globals.ID.ProcessValues.Pc+1,
			globals.ClientConfig.Instance_id,
			detalle.InstructionType,
			globals.ID.Tiempo,
			globals.ID.Dispositivo)

		return true

	case "INIT_PROC": //INIT_PROC (Archivo de instrucciones, Tamaño)
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s - TAM: %d - ARCHIVO: %s", detalle.ProcessValues.Pid, detalle.InstructionType, detalle.Tamaño, detalle.ArchiInstr)

		FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.ID.ProcessValues.Pid, globals.ID.ProcessValues.Pc,
			globals.ClientConfig.Instance_id,
			detalle.InstructionType,
			globals.ID.Tamaño,
			globals.ID.ArchiInstr)

		return true

	case "DUMP_MEMORY": //
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s", detalle.ProcessValues.Pid, detalle.InstructionType)
		FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.ID.ProcessValues.Pid,
			globals.ID.ProcessValues.Pc+1,
			globals.ClientConfig.Instance_id,
			detalle.InstructionType,
			0,
			"")
		return true

	case "EXIT":
		log.Printf("## PID: %d - Ejecutando -> INSTRUCCION: %s", detalle.ProcessValues.Pid, detalle.InstructionType)
		FinEjecucion(globals.ClientConfig.Ip_kernel,
			globals.ClientConfig.Port_kernel,
			globals.ID.ProcessValues.Pid,
			globals.ID.ProcessValues.Pc+1,
			globals.ClientConfig.Instance_id,
			detalle.InstructionType,
			0,
			"")
		return true

	default:
		fmt.Println("## ERROR -> Instruccion inválida.")
		return false
	}
}

func Write(ip string, port int, direccion int, contenido string) {

	var paquete utilsCPU.WriteStruct

	paquete.Contenido = contenido
	paquete.Direccion = direccion
	paquete.Pid = globals.ID.ProcessValues.Pid

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
		log.Printf("## Error al recibir respuesta.\n")
		return

	}

	if respuestaJSON.StatusCode != http.StatusOK {
		log.Printf("## Status de respuesta no fue la esperada.\n")
		return
	}
	defer respuestaJSON.Body.Close()

	//fmt.Printf("Conexion establecida con exito.\n")

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

	log.Printf("## Respuesta de Memoria: %s\n", respuesta.Mensaje) // Nos devuelve memoria el mensaje de escritura.

}

func Read(ip string, port int, direccion int, tamaño int) {

	var paquete utilsCPU.ReadStruct

	paquete.Direccion = direccion
	paquete.Tamaño = tamaño
	paquete.Pid = globals.ID.ProcessValues.Pid

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
		log.Printf("## ERROR -> Error al generar la peticion al server.")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("## ERROR -> Error al recibir respuesta.")
		return

	}

	if respuestaJSON.StatusCode != http.StatusOK {

		log.Printf("## ERROR -> Status de respuesta el server no fue la esperada.")
		return
	}
	defer respuestaJSON.Body.Close()

	// fmt.Printf("Conexion establecida con exito.\n")

	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta PaqueteRecibidoREAD

	err = json.Unmarshal(body, &(respuesta))
	if err != nil {
		log.Printf("## ERROR -> Error al decodificar el JSON.")
		return
	}

	globals.ID.ValorLeido = respuesta.Info               //guarda el valor leido en memoria en el tamanio que se paso por parametro
	globals.ID.PaginaCompleta = respuesta.PaginaCompleta //guarda pagina completa para guardar en cache

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
		log.Printf("## ERROR -> Error al convertir a json.")
		return
	}

	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/PCB", ip, puerto) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

	if err != nil {
		log.Printf("## ERROR -> Error al generar la peticion al server.")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	respuestaJSON, err := cliente.Do(req)

	if err != nil {
		log.Printf("## ERROR ->Error al recibir respuesta.")
		return
	}

	defer respuestaJSON.Body.Close()

	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta utilsCPU.RespuestaKernel

	if respuestaJSON.StatusCode != http.StatusOK {
		globals.Interruption = true
	}

	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("## ERROR -> Error al decodificar el JSON.")
	}

	//log.Printf("## Kernel -> Recibió correctamente el PID: %d y el PC: %d.", respuesta.Pid, respuesta.Pc)

}

// agregar en cache segun algoritmo
func AgregarEnCache(nroPagina int, direccionFisica int) {
	if globals.CachePaginas.Tamanio == 0 {
		return
	}
	var paginaCompleta []byte

	if len(globals.ID.PaginaCompleta) > 0 {
		paginaCompleta = globals.ID.PaginaCompleta
	} else {
		paginaCompleta = make([]byte, globals.ClientConfig.Page_size)
		copy(paginaCompleta, []byte(globals.ID.Datos))
	}

	if globals.ID.InstructionType == "WRITE" {
		copy(paginaCompleta[globals.ID.Desplazamiento:], []byte(globals.ID.Datos))
	}
	var contenido []byte

	if globals.ID.InstructionType == "READ" {
		contenido = globals.ID.ValorLeido
	}

	entrada := globals.EntradaCacheDePaginas{
		PID:             globals.ID.ProcessValues.Pid,
		NroPag:          nroPagina,
		PaginaCompleta:  globals.ID.PaginaCompleta,
		Frame:           globals.ID.Frame,
		Desplazamiento:  globals.ID.Desplazamiento,
		Contenido:       contenido,
		DireccionFisica: direccionFisica,
		Modificada:      globals.ID.InstructionType == "WRITE",
		BitUso:          true,
	}

	if len(globals.CachePaginas.Entradas) < globals.CachePaginas.Tamanio {
		globals.CachePaginas.Entradas = append(globals.CachePaginas.Entradas, entrada)
		return
	}

	switch globals.ClientConfig.Cache_replacement {
	case "CLOCK":
		ReemplazarConCLOCK(entrada)
	case "CLOCK-M":
		ReemplazarConCLOCKM(entrada)
	default:
		log.Printf("## ERROR -> Algoritmo de reemplazo incorrecto.")
	}
}

func AgregarEnTLB(nroPagina int, direccion int) {
	/*if globals.ClientConfig.Tlb_entries <= 0 {
		fmt.Printf("## ERROR -> Entradas de TLB: %d -> No hay TLB.", globals.ClientConfig.Tlb_entries)
		return
	}*/

	tlb := &globals.Tlb
	pid := globals.ID.ProcessValues.Pid

	for i, entrada := range tlb.Entradas {
		if entrada.PID == pid && entrada.NroPagina == nroPagina {
			tlb.Entradas[i].UltimoAcceso = time.Now().UnixNano()
			return
		}
	}

	entrada := globals.Entrada{
		PID:          pid,
		NroPagina:    nroPagina,
		Direccion:    direccion,
		UltimoAcceso: time.Now().UnixNano(),
	}

	if len(tlb.Entradas) < tlb.Tamanio {
		tlb.Entradas = append(tlb.Entradas, entrada)
	} else {
		switch globals.ClientConfig.Tlb_replacement {
		case "FIFO":
			ReemplazarTLB_FIFO(entrada)
		case "LRU":
			ReemplazarTLB_LRU(entrada)
		default:
			log.Printf("## ERROR -> Algoritmo de reemplazo incorrecto.")
		}
	}
}

func VaciarCache(pid int) {
	var nuevasEntradas []globals.EntradaCacheDePaginas

	for _, entrada := range globals.CachePaginas.Entradas {
		if entrada.PID == pid {
			if entrada.Modificada {
				Write(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory,
					entrada.Frame, string(entrada.PaginaCompleta),
				)
			}
		} else {
			nuevasEntradas = append(nuevasEntradas, entrada)
		}
	}

	globals.CachePaginas.Entradas = nuevasEntradas
	log.Printf("## Proceso PID: %d desalojado correctamente - VACIO CACHE.", pid) //lo que interpreto es que se eliminan las entradas del proceso desalojado no TODAS(de todos los procesos)
}
