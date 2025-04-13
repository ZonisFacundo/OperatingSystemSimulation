package utilsMemoria

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func ConfigurarLogger() {
	logFile, err := os.OpenFile("memory.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

type Handshakepaquete struct {
	Instruccion string `json:"instruccion"`
}

type HandshakepaqueteKernel struct {
	NombreCodigo   string `json:"nombreCodigo"`
	TamanioProceso int    `json:"taamanioProceso"`
}

type respuestaalCPU struct {
	Mensaje string `json:"message"`
}

func RetornoClienteCPUServidorMEMORIA(w http.ResponseWriter, r *http.Request) {

	var request Handshakepaquete

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("El cliente nos mando esto: \n instruccion: %s.\n", request.Instruccion)

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuestaCpu respuestaalCPU
	respuestaCpu.Mensaje = "Recibi de CPU"
	respuestaJSON, err := json.Marshal(respuestaCpu)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func RetornoClienteKernelServidorMEMORIA(w http.ResponseWriter, r *http.Request) {

	var request HandshakepaqueteKernel

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("El cliente nos mando esto: \n nombre pseudocodigo: %s, tamanio proceso: %d.\n", request.NombreCodigo, request.TamanioProceso)

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuestaCpu respuestaalCPU
	respuestaCpu.Mensaje = "Recibi de Kernel"
	respuestaJSON, err := json.Marshal(respuestaCpu)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}
