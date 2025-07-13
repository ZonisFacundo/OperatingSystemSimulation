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

		nroPagina := globals.ID.DireccionLog / globals.ClientConfig.Page_size
		globals.ID.NroPag = nroPagina

		if mmu.EstaTraducida(nroPagina) {
			globals.Tlb.Entradas[globals.ID.PosicionPag].UltimoAcceso = time.Now().UnixNano()
			Execute(globals.ID)
		} else {
			direccionAEnviar := mmu.TraducirDireccion(globals.ID.DireccionLog, instruccion.ProcessValues.Pid, nroPagina)
			log.Println("que es esto? (2): ", direccionAEnviar)
			EnvioDirLogica(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, direccionAEnviar)

			globals.ID.DireccionFis = (globals.ID.Frame * globals.ClientConfig.Page_size) + globals.ID.Desplazamiento

			log.Printf("desplazamiento: %d", globals.ID.Desplazamiento)
			log.Printf("frame a enviar %d", globals.ID.Frame)
			log.Printf("page size: %d", globals.ClientConfig.Page_size)
			log.Printf("direccion a enviar %d", globals.ID.DireccionFis)
			//Mandar direccion fisica a la TLB junto con el numero de página así queda guardada en "caché".
		}

	case "WRITE":
		instruccion.DireccionLog, _ = strconv.Atoi(partesDelString[1])
		instruccion.Datos = partesDelString[2]

		globals.ID.DireccionLog = instruccion.DireccionLog
		globals.ID.Datos = instruccion.Datos

		nroPagina := globals.ID.DireccionLog / globals.ClientConfig.Page_size
		globals.ID.NroPag = nroPagina

		if mmu.EstaTraducida(nroPagina) {
			globals.Tlb.Entradas[globals.ID.PosicionPag].UltimoAcceso = time.Now().UnixNano()
			Execute(globals.ID)
		} else {

			direccionAEnviar := mmu.TraducirDireccion(globals.ID.DireccionLog, instruccion.ProcessValues.Pid, nroPagina)

			log.Println("dir?: ", direccionAEnviar)

			EnvioDirLogica(globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, direccionAEnviar)
			if globals.ID.Frame < 0 {
				log.Printf("ERROR ERROR ERROR FACU, te imprimo el frame %d", globals.ID.Frame)
			} else {
				globals.ID.DireccionFis = (globals.ID.Frame * globals.ClientConfig.Page_size) + globals.ID.Desplazamiento

				log.Printf("desplazamiento: %d", globals.ID.Desplazamiento)
				log.Printf("frame a enviar %d", globals.ID.Frame)
				log.Printf("page size: %d", globals.ClientConfig.Page_size)
				log.Printf("direccion a enviar %d", globals.ID.DireccionFis)
			}
			// aca habria que agregar la direccion traducida a la tlb y trabajar con un alg de reemplazo si la tlb esta llena
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
