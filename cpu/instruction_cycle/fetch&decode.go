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

// ¿Es correcto? porque, si no, no se guarda en la variable

/*
Fetch
La primera etapa del ciclo consiste en buscar la próxima instrucción a ejecutar. En este trabajo práctico cada instrucción deberá ser pedida al módulo Memoria
utilizando el Program Counter (también llamado Instruction Pointer) que representa el número de instrucción a buscar relativo al hilo en ejecución.*/

func SolicitudPIDyPCAKernel(ip string, puerto int, instancia string) {

	var paquete utilsCPU.HandshakeCPU

	paquete.Ip = ip
	paquete.Puerto = puerto

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.\n")
		return
	}
	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/handshake", ip, puerto) //url del server

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

	log.Printf("Conexion establecida con exito.\n")
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	// log.Println("Conexion establecida con exito en la instancia: ", instancia) -> Mostrar la instancia en ejecución (próximamente programada).

	err = json.Unmarshal(body, &globals.Instruction)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
	}

	log.Printf("PID: %d - PC: %d.\n", globals.Instruction.Pid, globals.Instruction.Pc)
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

	url := fmt.Sprintf("http://%s:%d/INSTRUCCIONES", globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)

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

	var respuesta globals.Instruccion

	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}

	log.Printf("Instruction given: %s.\n", respuesta.InstructionType) // Nos manda memoria la instrucción.
}

/*
Decode
Esta etapa consiste en interpretar qué instrucción es la que se va a ejecutar y si la misma requiere de una traducción de dirección lógica a dirección física.
*/

func Decode() {

}
