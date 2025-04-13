package utilsMemoria

import (
	"io"
	"log"
	"os"
	"net/http"
	"encoding/json"
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


type respuestaalCPU struct {
	Mensaje string `json:"message"`
}



func ConexionRecibidaCPU(w http.ResponseWriter, r *http.Request) {

	var request Handshakepaquete

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("el cliente nos mando esto: \n instruccion: %s", request.Instruccion)

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuestaCpu respuestaalCPU
	respuestaCpu.Mensaje = "me pinto mandarle un string al cliente"
	respuestaJSON, err := json.Marshal(respuestaCpu)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}
