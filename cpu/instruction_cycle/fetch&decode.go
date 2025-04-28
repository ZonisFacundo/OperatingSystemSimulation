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

type Instruccion struct { // instruccion obtenida de memoria
	ProcessValues   utilsCPU.Proceso      `json:"instruction"`  //Valores de PID y PC
	Interrup        utilsCPU.Interrupcion `json:"interruption"` //Valores de la interrupción.
	Direccion       int                   `json:"adress"`       //Para Read and Write -> Dirección lógica que pasa memoria.
	InstructionType string                `json:"message"`      //Contexto de la ejecución, es decir, la string que entra en el execute.
	Valor           *int                  `json:"value"`        //Parámetro para GOTO
	Tamaño          *int                  `json:"size"`         //Parámetro para el READ e INIT_PROC.
	Tiempo          *int                  `json:"time"`         //Parámetro para NOOP.
	Datos           *string               `json:"datos"`        //Parámetro para el WRITE.
}

type PaqueteRecibidoMemoria struct {
	Mensaje string `json:"message"`
}
type PaqueteRecibidoWRITE struct {
	Mensaje string `json:"message"`
}
type PaqueteRecibidoREAD struct {
	ValorInstruccion string `json:"message"`
}

/*
Fetch
La primera etapa del ciclo consiste en buscar la próxima instrucción a ejecutar. En este trabajo práctico cada instrucción deberá ser pedida al módulo Memoria
utilizando el Program Counter (también llamado Instruction Pointer) que representa el número de instrucción a buscar relativo al hilo en ejecución.*/

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

/*
Decode
Esta etapa consiste en interpretar qué instrucción es la que se va a ejecutar y si la misma requiere de una traducción de dirección lógica a dirección física.
*/

func Decode() {

}
