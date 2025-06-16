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

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/mmu"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

type PaqueteRecibidoMemoria struct {
	Mensaje string `json:"message"`
}
type PaqueteRecibidoWRITE struct {
	Mensaje int `json:"message"`
}
type PaqueteRecibidoREAD struct {
	Info string `json:"info"`
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
		ProcesoActual:       instruccion.ProcessValues,
		TamPagina:           64,
		Niveles:             5,
		Cant_entradas_tabla: 4,
		TablasPaginas:       make(map[int]int)}

	partesDelString := strings.Fields(instruccion.InstructionType)

	instruccion.InstructionType = partesDelString[0]

	globals.ID.InstructionType = instruccion.InstructionType

	// Instruccion ¿tipo?
	switch instruccion.InstructionType {

	case "READ":
		instruccion.DireccionLog, _ = strconv.Atoi(partesDelString[1])
		instruccion.Tamaño, _ = strconv.Atoi(partesDelString[2])

		globals.ID.DireccionLog = instruccion.DireccionLog
		globals.ID.Tamaño = instruccion.Tamaño

		direccionAEnviar := mmu.TraducirDireccion(globals.ID.DireccionLog, memoryManagement, instruccion.ProcessValues.Pid)

		EnvioDirLogica(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, direccionAEnviar)

		log.Printf("frame: %d", globals.ID.Frame)

		globals.ID.DireccionFis = (globals.ID.Frame * memoryManagement.TamPagina) + globals.ID.Desplazamiento

	case "WRITE":
		instruccion.DireccionLog, _ = strconv.Atoi(partesDelString[1])
		instruccion.Datos = partesDelString[2]

		globals.ID.DireccionLog = instruccion.DireccionLog
		globals.ID.Datos = instruccion.Datos

		log.Printf("dir: %d", instruccion.DireccionLog)

		direccionAEnviar := mmu.TraducirDireccion(globals.ID.DireccionLog, memoryManagement, instruccion.ProcessValues.Pid)

		log.Printf("direccion a enviar: %d", direccionAEnviar)

		EnvioDirLogica(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, direccionAEnviar)

		log.Printf("frame: %d", globals.ID.Frame)

		globals.ID.DireccionFis = (globals.ID.Frame * memoryManagement.TamPagina) + globals.ID.Desplazamiento

	case "GOTO":
		instruccion.Valor, _ = strconv.Atoi(partesDelString[1])
		globals.ID.Valor = instruccion.Valor

	case "IO":
		instruccion.Parametro1, _ = strconv.Atoi(partesDelString[2])
		instruccion.Parametro2 = partesDelString[1]

		globals.ID.Parametro1 = instruccion.Parametro1
		globals.ID.Parametro2 = instruccion.Parametro2

	case "INIT_PROC":
		instruccion.Parametro1, _ = strconv.Atoi(partesDelString[2])
		instruccion.Parametro2 = partesDelString[1]

		globals.ID.Parametro1 = instruccion.Parametro1
		globals.ID.Parametro2 = instruccion.Parametro2

	default:
		log.Printf("Nada que modificar, continua la ejecución.")
		//Execute(instruccion)
	}

	// READ & WRITE -> Traduzco -> Execute
	// Else -> Execute
}
