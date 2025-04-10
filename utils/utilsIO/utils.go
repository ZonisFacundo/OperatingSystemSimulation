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
	Nombre string `json:"name"` //lo de name, ip, port es como va a construir el json la maquina cuando lo pasemos a json
	Ip     string `json:"ip"`   // es fundamental ponerlo
	Puerto int    `json:"port"`
}

type RespuestaHandshakeKernel struct { // aca va el formato que va a tener lo que esperas del server
	Mensaje string `json:"message"`
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

	respuestaJSON, err := cliente.Do(req) //recibe la respuesta del server
	/* que tipo de dato tiene respuestaJSON?
		type respuestaJSON struct {
	    Status     string
	    StatusCode int
	    Header     Header
	    Body       io.ReadCloser  // ← This is what you're accessing
	    // ... other fields ...


		ya definido por go de esa forma
	*/

	if err != nil {
		fmt.Printf("error al recibir respuesta")
		return

	}

	if respuestaJSON.StatusCode != http.StatusOK {

		fmt.Printf("status de respuesta el server no fue la esperada")
		return
	}
	defer respuestaJSON.Body.Close() //cerramos algo supuestamente importante de cerrar pero no se que hace

	fmt.Printf("conexion establecida con exito \n")
	//pasamos de JSON a formato bytes lo que nos paso el paquete
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	//pasamos la respuesta de JSON a formato paquete que nos mando el server

	var respuesta RespuestaHandshakeKernel //para eso declaramos una variable con el struct que esperamos que nos envie el server
	err = json.Unmarshal(body, &respuesta) //pasamos de bytes al formato de nuestro paquete lo que nos mando el server
	if err != nil {
		fmt.Printf("error al decodificar el JSON")
		return
	}
	fmt.Printf("la respuesta del server fue: %s/n", respuesta.Mensaje)
	//en mi caso era un mensaje, por eso el struct tiene mensaje string, vos por ahi estas esperando 14 ints, no necesariamente un struct

}
