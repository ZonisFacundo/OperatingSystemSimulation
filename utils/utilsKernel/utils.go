package utilsKernel

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	//"encoding/json"
	"net/http"
)

type Handshakepaquete struct {
	Nombre string `json:"name"`
	Ip     string `json:"ip"`
	Puerto int    `json:"port"`
}

type respuestaalIO struct {
	Mensaje string `json:"message"`
}

func ConfigurarLogger() {
	logFile, err := os.OpenFile("kernel.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func ConexionRecibida(w http.ResponseWriter, r *http.Request) {

	var request Handshakepaquete

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	fmt.Printf("el cliente nos mando esto: \n nombre: %s  \n puerto: %d \n IP: %s \n", request.Nombre, request.Puerto, request.Ip)

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuesta respuestaalIO
	respuesta.Mensaje = "me pinto mandarle un string al cliente"
	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}
