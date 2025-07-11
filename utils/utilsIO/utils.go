package utilsIO

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Handshakepaquete struct {
	Nombre string `json:"name"`
	Ip     string `json:"ip"`
	Puerto int    `json:"port"`
}

type HandshakepaqueteFin struct {
	Nombre string `json:"name"`
}

type RespuestaHandshakeKernel struct {
	Mensaje string `json:"message"`
}

type PaqueteRecibidoIO struct {
	Pid    int `json:"pid"`
	Tiempo int `json:"tiempo"`
}

type PaqueteRespuestaKERNEL struct {
	Mensaje string `json:"message"`
}

func ConfigurarLogger(ioId string) {
	logFileName := fmt.Sprintf("IO-%s.log", ioId)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.SetPrefix(fmt.Sprintf("[IO-%s] ", ioId))
}

func NotificarFinalizacionAlKernel(nombre string, ipKernel string, puertoKernel int, ipIO string, puertoIO int) {
	var paquete HandshakepaqueteFin
	paquete.Nombre = nombre

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("error al convertir a json")
		return
	}
	cliente := http.Client{}

	url := fmt.Sprintf("http://%s:%d/finIO", ipKernel, puertoKernel)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson))

	if err != nil {
		log.Printf("error al generar la peticion al server")
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

	//log.Printf("Conexion establecida con exito \n")
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta RespuestaHandshakeKernel
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON")
		return
	}

}

func PeticionClienteIOServidorKERNEL(nombre string, ipKernel string, puertoKernel int, ipIO string, puertoIO int) {

	var paquete Handshakepaquete
	paquete.Nombre = nombre
	paquete.Ip = ipIO
	paquete.Puerto = puertoIO

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("error al convertir a json")
		return
	}
	cliente := http.Client{}

	url := fmt.Sprintf("http://%s:%d/IO", ipKernel, puertoKernel)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson))

	if err != nil {
		log.Printf("error al generar la peticion al server")
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

	//log.Printf("Conexion establecida con exito \n")
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta RespuestaHandshakeKernel
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON")
		return
	}
	//log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)

}

func RetornoClienteKERNELServidorIO(w http.ResponseWriter, r *http.Request) {

	var request PaqueteRecibidoIO

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("## PID: %d - Inicio de IO - Tiempo: %d", request.Pid, request.Tiempo)
	IniciarSleep(request.Tiempo)

	log.Printf("## PID: %d - Fin de IO", request.Pid)

	var respuestaIO PaqueteRespuestaKERNEL
	respuestaIO.Mensaje = "I/O Finalizado"
	respuestaJSON, err := json.Marshal(respuestaIO)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func IniciarSleep(tiempo int) {
	time.Sleep(time.Duration(tiempo) * time.Millisecond)

}
