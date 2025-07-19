package instruction_cycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

type PaqueteRecibidoMemoria struct {
	Mensaje string `json:"message"`
}
type PaqueteRecibidoWRITE struct {
	Mensaje string `json:"message"`
}
type PaqueteRecibidoREAD struct {
	Info           []byte `json:"info"`
    PaginaCompleta []byte `json:"pag"`
}

func Fetch(pid int, pc int, ip string, puerto int) {

	var paquete utilsCPU.HandshakeMemory

	paquete.Pc = pc
	paquete.Pid = pid

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("## ERROR -> Error al convertir a json.")
		return
	}
	cliente := http.Client{}

	// log.Printf("Process ID: %d - FETCH - Program Counter: %d.\n", pid, pc) // log mínimo y obligatorio

	url := fmt.Sprintf("http://%s:%d/INSTRUCCIONES", ip, puerto)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(PaqueteFormatoJson))

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

	//fmt.Printf("Conexion establecida con exito.\n")
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta PaqueteRecibidoMemoria

	err = json.Unmarshal(body, &(respuesta))
	if err != nil {
		log.Printf("## ERROR -> Error al decodificar el JSON.")
		return
	}

	log.Printf("## Instruccion recibida -> %s.", respuesta.Mensaje) // Nos manda memoria la instrucción.

	globals.ID.InstructionType = respuesta.Mensaje
}

func Decode(instruccion globals.Instruccion) {

	memoryManagement := mmu.MMU{
		Pc:                  instruccion.ProcessValues.Pc,
		Pid:                 instruccion.ProcessValues.Pid,
		TamPagina:           globals.ClientConfig.Page_size,
		Niveles:             globals.ClientConfig.Niveles,
		Cant_entradas_tabla: globals.ClientConfig.Entradas,
		TablasPaginas:       make(map[int]int),
	}

	partesDelString := strings.Fields(instruccion.InstructionType)
	instruccion.InstructionType = partesDelString[0]

	globals.ID.InstructionType = instruccion.InstructionType

	log.Printf("## PID: <%d> - FETCH - Program Counter: <%d>.", instruccion.ProcessValues.Pid, instruccion.ProcessValues.Pc)

	// Instruccion ¿tipo?
	switch instruccion.InstructionType {

	case "READ":
		instruccion.DireccionLog, _ = strconv.Atoi(partesDelString[1])
		instruccion.Tamaño, _ = strconv.Atoi(partesDelString[2])

		globals.ID.DireccionLog = instruccion.DireccionLog
		globals.ID.Tamaño = instruccion.Tamaño

		globals.ID.ProcessValues.Pid = instruccion.ProcessValues.Pid

		nroPagina := globals.ID.DireccionLog / memoryManagement.TamPagina
		globals.ID.NroPag = nroPagina

		traducida := mmu.EstaTraducida(nroPagina)

		if traducida {
			globals.Tlb.Entradas[globals.ID.PosicionPag].UltimoAcceso = time.Now().UnixNano()
		} else {
			direccionAEnviar := mmu.TraducirDireccion(globals.ID.DireccionLog, memoryManagement, instruccion.ProcessValues.Pid, nroPagina)
			EnvioDirLogica(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, direccionAEnviar)
			globals.ID.DireccionFis = (globals.ID.Frame * globals.ClientConfig.Page_size) + globals.ID.Desplazamiento
		}

	case "WRITE":
		instruccion.DireccionLog, _ = strconv.Atoi(partesDelString[1])
		instruccion.Datos = partesDelString[2]

		globals.ID.DireccionLog = instruccion.DireccionLog
		globals.ID.Datos = instruccion.Datos

		globals.ID.ProcessValues.Pid = instruccion.ProcessValues.Pid

		globals.ID.NroPag = globals.ID.DireccionLog / memoryManagement.TamPagina

		if mmu.EstaTraducida(globals.ID.NroPag) {
			globals.Tlb.Entradas[globals.ID.PosicionPag].UltimoAcceso = time.Now().UnixNano()

		} else {

			direccionAEnviar := mmu.TraducirDireccion(globals.ID.DireccionLog, memoryManagement, instruccion.ProcessValues.Pid, globals.ID.NroPag)

			EnvioDirLogica(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, direccionAEnviar)

			if globals.ID.Frame < 0 {
				log.Printf("## ERROR -> Frame inválido: %d", globals.ID.Frame)

			} else {
				globals.ID.DireccionFis = (globals.ID.Frame * globals.ClientConfig.Page_size) + globals.ID.Desplazamiento
				log.Printf("## Direccion Fisica: %d", globals.ID.DireccionFis)
			}
		}

	case "INIT_PROC":
		instruccion.Tamaño, _ = strconv.Atoi(partesDelString[2])
		instruccion.ArchiInstr = partesDelString[1]

		globals.ID.ArchiInstr = instruccion.ArchiInstr
		globals.ID.Tamaño = instruccion.Tamaño

	case "IO":
		instruccion.Dispositivo = partesDelString[1]
		instruccion.Tiempo, _ = strconv.Atoi(partesDelString[2])

		globals.ID.Dispositivo = instruccion.Dispositivo
		globals.ID.Tiempo = instruccion.Tiempo

	default:
		log.Printf("## Nada que modificar, continua la ejecución.")
	}
}
