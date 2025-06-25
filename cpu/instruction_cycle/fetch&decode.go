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
	Mensaje string `json:"message"`
}
type PaqueteRecibidoREAD struct {
	Info byte `json:"info"`
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

		nroPagina := globals.ID.DireccionLog / memoryManagement.TamPagina

		if mmu.EstaTraducida(nroPagina) {
			log.Printf("entro aca (1)")
			Execute(globals.ID)
		} else {
			log.Printf("entro aca (2)")
			direccionAEnviar := mmu.TraducirDireccion(globals.ID.DireccionLog, memoryManagement, instruccion.ProcessValues.Pid, nroPagina)
			EnvioDirLogica(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, direccionAEnviar)
			globals.ID.DireccionFis = (globals.ID.Frame * globals.ClientConfig.Page_size) + globals.ID.Desplazamiento
			//Mandar direccion fisica a la TLB junto con el numero de página así queda guardada en "caché".
		}

	case "WRITE":
		instruccion.DireccionLog, _ = strconv.Atoi(partesDelString[1])
		instruccion.Datos = partesDelString[2]

		globals.ID.DireccionLog = instruccion.DireccionLog
		globals.ID.Datos = instruccion.Datos

		nroPagina := globals.ID.DireccionLog / memoryManagement.TamPagina
		// mmu despues deberiamos hacerlo global, porque son parametros que nos deberia pasar memoria (tabla de pags)

		if mmu.EstaTraducida(nroPagina) {
			log.Printf("entro aca (1)")
			Execute(globals.ID)
		} else {
			log.Printf("entro aca (2)")
			log.Printf("desplazamiento: %d", globals.ID.Desplazamiento)
			direccionAEnviar := mmu.TraducirDireccion(globals.ID.DireccionLog, memoryManagement, instruccion.ProcessValues.Pid, nroPagina)
			EnvioDirLogica(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, direccionAEnviar)
			globals.ID.DireccionFis = (globals.ID.Frame * globals.ClientConfig.Page_size) + globals.ID.Desplazamiento
			// aca habria que agregar la direccion traducida a la tlb y trabajar con un alg de reemplazo si la tlb esta llena
		}
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
