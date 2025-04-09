package utilsIO

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Handshakepaquete struct {
	Nombre string `json:"name"`
	Ip     string `json:"ip"`
	Puerto int    `json:"port"`
}

func ConfigurarLogger() {
	logFile, err := os.OpenFile("io.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

/*
conexion entre IO (Client) con Kernel (Server)
enviamos handshake con datos del modulo y esperamos respuesta
*/

func HandshakeAKernel(nombre string, ip string, puerto int) {

	var paquete Handshakepaquete
	paquete.Nombre = nombre
	paquete.Ip = ip
	paquete.Puerto = puerto

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		//aca tiene que haber un logger
		fmt.Printf("error al convertir a json")
		return
	}
	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/handshake", paquete.Ip, paquete.Puerto) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

	if err != nil {
		//aca tiene que haber un logger
		fmt.Printf("error al generar la peticion al server")
		return
	}

	req.Header.Set("Content-Type", "application/json") //le avisa al server que manda la data en json format

	respuesta, err := cliente.Do(req) //recibe la respuesta del server

	if err != nil {
		fmt.Printf("error al recibir respuesta")
		return

	}

	if respuesta.StatusCode != http.StatusOK {

		fmt.Printf("status de respuesta el server no fue la esperada")
		return
	}

}
