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
	Info []byte `json:"info"`
}

func Fetch(pid int, pc int, ip string, puerto int) {

	var paquete utilsCPU.HandshakeMemory

	paquete.Pc = pc
	paquete.Pid = pid

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.\n")
		return
	}
	cliente := http.Client{}

	// log.Printf("Process ID: %d - FETCH - Program Counter: %d.\n", pid, pc) // log mínimo y obligatorio

	url := fmt.Sprintf("http://%s:%d/INSTRUCCIONES", ip, puerto)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(PaqueteFormatoJson))

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

	var respuesta PaqueteRecibidoMemoria

	err = json.Unmarshal(body, &(respuesta))
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}

	log.Printf("Instruction given: %s\n", respuesta.Mensaje) // Nos manda memoria la instrucción.

	globals.ID.InstructionType = respuesta.Mensaje
}

func Decode(instruccion globals.Instruccion) {

	memoryManagement := mmu.MMU{
		Pc:                  instruccion.Pc,
		Pid:                 instruccion.Pid,
		TamPagina:           64,
		Niveles:             5,
		Cant_entradas_tabla: 4,
		TablasPaginas:       make(map[int]int)}

	partesDelString := strings.Fields(instruccion.InstructionType)

	instruccion.InstructionType = partesDelString[0]

	globals.ID.InstructionType = instruccion.InstructionType

	log.Printf("## PID: <%d> - FETCH - Program Counter: <%d>.", instruccion.Pid, instruccion.Pc)

	// Instruccion ¿tipo?
	switch instruccion.InstructionType {

	case "READ":
		instruccion.DireccionLog, _ = strconv.Atoi(partesDelString[1])
		instruccion.Tamaño, _ = strconv.Atoi(partesDelString[2])

		globals.ID.DireccionLog = instruccion.DireccionLog
		globals.ID.Tamaño = instruccion.Tamaño

		globals.ID.Pid = instruccion.Pid

		nroPagina := globals.ID.DireccionLog / memoryManagement.TamPagina
		globals.ID.NroPag = nroPagina

		traducida := mmu.EstaTraducida(nroPagina)

		if traducida {
			globals.Tlb.Entradas[globals.ID.PosicionPag].UltimoAcceso = time.Now().UnixNano()
		} else {
			direccionAEnviar := mmu.TraducirDireccion(globals.ID.DireccionLog, memoryManagement, instruccion.Pid, nroPagina)
			EnvioDirLogica(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, direccionAEnviar)
			globals.ID.DireccionFis = (globals.ID.Frame * globals.ClientConfig.Page_size) + globals.ID.Desplazamiento
		}

	case "WRITE":
		instruccion.DireccionLog, _ = strconv.Atoi(partesDelString[1])
		instruccion.Datos = partesDelString[2]

		globals.ID.DireccionLog = instruccion.DireccionLog
		globals.ID.Datos = instruccion.Datos

		globals.ID.Pid = instruccion.Pid

		globals.ID.NroPag = globals.ID.DireccionLog / memoryManagement.TamPagina
		// mmu despues deberiamos hacerlo global, porque son parametros que nos deberia pasar memoria (tabla de pags)

		if mmu.EstaTraducida(globals.ID.NroPag) {
			globals.Tlb.Entradas[globals.ID.PosicionPag].UltimoAcceso = time.Now().UnixNano()
		} else {

			direccionAEnviar := mmu.TraducirDireccion(globals.ID.DireccionLog, memoryManagement, instruccion.Pid, globals.ID.NroPag)

			EnvioDirLogica(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, direccionAEnviar)

			if globals.ID.Frame < 0 {
				log.Printf("ERROR: Frame inválido: %d", globals.ID.Frame)
			} else {
				globals.ID.DireccionFis = (globals.ID.Frame * globals.ClientConfig.Page_size) + globals.ID.Desplazamiento
				log.Printf("Dirección física calculada: %d", globals.ID.DireccionFis)
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
		log.Printf("Nada que modificar, continua la ejecución.")
	}
}
